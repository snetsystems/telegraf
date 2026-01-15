// Package nutanix implements an nutanix input plugin for Telegraf
//
//go:generate ../../../tools/readme_config_includer/generator
package nutanix

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/common/tls"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// DO NOT REMOVE THE NEXT TWO LINES! This is required to embed the sampleConfig data.
//
//go:embed sample.conf
var sampleConfig string

type Nutanix struct {
	URL             string          `toml:"url"`
	Username        string          `toml:"username"`
	Password        string          `toml:"password"`
	EnabledServices []string        `toml:"enabled_services"`
	Log             telegraf.Logger `toml:"-"`

	ResponseTimeout config.Duration `toml:"response_timeout"`
	MaxRetries      int             `toml:"max_retries"`
	tls.ClientConfig

	client        *http.Client
	sessionCookie []*http.Cookie

	mu              sync.Mutex
	failureCount    int  // Track failure count for the entire connection
	skipLogged      bool // Track if we've already logged the skip message
}

func (*Nutanix) SampleConfig() string {
	return sampleConfig
}

func (n *Nutanix) Init() error {
	if n.URL == "" {
		return errors.New("url cannot be empty")
	}

	if n.Username == "" || n.Password == "" {
		return errors.New("username or password can not be empty string")
	}

	if n.MaxRetries <= 0 {
		n.MaxRetries = 3 // Default to 3 retries
	}

	client, err := n.createHTTPClient()
	if err != nil {
		return err
	}

	n.client = client
	n.failureCount = 0
	n.skipLogged = false

	return nil
}

// Gather gathers resources from the Nutanix API and accumulates metrics.  This
// implements the Input interface.
func (n *Nutanix) Gather(acc telegraf.Accumulator) error {
	// Check if we've exceeded max retries for the entire connection
	n.mu.Lock()
	failCount := n.failureCount
	n.mu.Unlock()

	if failCount >= n.MaxRetries {
		// Log only once when we start skipping
		n.mu.Lock()
		defer n.mu.Unlock()
		if !n.skipLogged {
			n.Log.Infof("Skipping all services: exceeded max consecutive failures (%d). Will retry on next successful connection.", n.MaxRetries)
			n.skipLogged = true
		}
		return nil
	}

	callDuration := map[string]interface{}{}
	hasError := false
	var lastError error

	for _, service := range n.EnabledServices {
		start := time.Now()
		if err := n.gatherGeneric(service, acc); err != nil {
			hasError = true
			lastError = err
			acc.AddError(fmt.Errorf("failed to get resource %q: %w", service, err))
		}
		callDuration[service] = time.Since(start).Nanoseconds()
	}

	// Update failure count based on overall result
	var newFailCount int
	n.mu.Lock()
	if hasError {
		n.failureCount++
		newFailCount = n.failureCount
	} else {
		// Reset failure count on success
		n.failureCount = 0
		n.skipLogged = false // Reset skip log flag on success
	}
	n.mu.Unlock()

	// Log errors outside of lock to avoid holding lock during logging
	if hasError {
		if newFailCount >= n.MaxRetries {
			acc.AddError(fmt.Errorf("connection failed (exceeded max consecutive failures %d): %w", n.MaxRetries, lastError))
		} else {
			acc.AddError(fmt.Errorf("connection failed (consecutive failures %d/%d): %w", newFailCount, n.MaxRetries, lastError))
		}
	}

	return nil
}

func (n *Nutanix) gatherGeneric(service string, acc telegraf.Accumulator) error {
	payload, ok := metricPayloads[service]
	if !ok {
		return fmt.Errorf("%s payload not defined", service)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := n.doWithSessionRetry(data)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	return parseMetrics(service, body, acc)
}

func (n *Nutanix) createHTTPClient() (*http.Client, error) {
	tlsCfg, err := n.ClientConfig.TLSConfig()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
		Timeout: time.Duration(n.ResponseTimeout),
	}

	return client, nil
}

func (n *Nutanix) doWithSessionRetry(payload []byte) (*http.Response, error) {
	req1, err := http.NewRequest("POST", n.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req1.Header.Set("Content-Type", "application/json")

	n.mu.Lock()
	for _, c := range n.sessionCookie {
		req1.AddCookie(c)
	}
	n.mu.Unlock()

	resp, err := n.client.Do(req1)
	if err != nil {
		return nil, err
	}
	// If not unauthorized, check for other error status codes
	if resp.StatusCode != http.StatusUnauthorized {
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return resp, nil
	}
	resp.Body.Close()

	// Retry with basic auth
	req2, err := http.NewRequest("POST", n.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req2.SetBasicAuth(n.Username, n.Password)
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := n.client.Do(req2)
	if err != nil {
		return nil, err
	}

	// Check status code after retry
	if resp2.StatusCode == http.StatusOK {
		newCookies := []*http.Cookie{}
		for _, c := range resp2.Cookies() {
			if strings.HasPrefix(c.Name, "NTNX_") {
				newCookies = append(newCookies, c)
			}
		}
		n.mu.Lock()
		n.sessionCookie = newCookies
		n.mu.Unlock()
		return resp2, nil
	}

	// Authentication failed or other error
	body, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	if resp2.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("authentication failed (401 Unauthorized): invalid username or password")
	}
	return nil, fmt.Errorf("API request failed with status %d: %s", resp2.StatusCode, string(body))
}

func parseMetrics(key string, body []byte, acc telegraf.Accumulator) error {
	var resp GroupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse cluster JSON: %w", err)
	}

	var tagKey string

	switch key {
	case "cluster":
		tagKey = "cluster_id"
	case "hosts":
		tagKey = "node_id"
	case "disks":
		tagKey = "disk_id"
	case "vms":
		tagKey = "vm_id"
	default:
		tagKey = "entity_id"
	}

	for _, group := range resp.GroupResults {
		for _, entity := range group.EntityResults {
			tags := map[string]string{
				tagKey: entity.EntityID,
			}

			fields := make(map[string]interface{})

			for _, metric := range entity.Data {
				name := metric.Name

				if len(metric.Values) == 0 || len(metric.Values[0].Values) == 0 {
					continue
				}

				if name == "_created_timestamp_usecs_" {
					continue
				}

				raw := metric.Values[0].Values[0]

				switch key {
				case "cluster":
					if name == "cluster_name" {
						tags[name] = raw
						continue
					}
				case "hosts":
					if name == "cluster" {
						tags["cluster_id"] = raw
						continue
					}
					if name == "cluster_name" || name == "node_name" {
						tags[name] = raw
						continue
					}
				case "disks":
					if name == "cluster" {
						tags["cluster_id"] = raw
						continue
					}
					if name == "node" {
						tags["node_id"] = raw
						continue
					}
					if name == "cluster_name" || name == "node_name" ||
						name == "serial" || name == "storage_tier" {
						tags[name] = raw
						continue
					}
				case "vms":
					if name == "cluster" {
						tags["cluster_id"] = raw
						continue
					}
					if name == "node" {
						tags["node_id"] = raw
						continue
					}
					if name == "cluster_name" || name == "node_name" ||
						name == "project_name" || name == "vm_name" ||
						name == "project_reference" {
						tags[name] = raw
						continue
					}
				}

				fieldKey := strings.ReplaceAll(name, ".", "_")

				switch metric.DataType {
				case "int64", "uint64":
					if val, err := strconv.ParseInt(raw, 10, 64); err == nil {
						fields[fieldKey] = val
						if strings.HasSuffix(name, "_ppm") {
							percentName := strings.TrimSuffix(fieldKey, "_ppm") + "_percent"
							fields[percentName] = float64(val) / 10000.0
						}
					}
				case "double":
					if val, err := strconv.ParseFloat(raw, 64); err == nil {
						fields[fieldKey] = val
						if strings.HasSuffix(name, "_ppm") {
							percentName := strings.TrimSuffix(fieldKey, "_ppm") + "_percent"
							fields[percentName] = val / 10000.0
						}
					}
				case "string":
					fields[fieldKey] = raw
				default:
					fields[fieldKey] = raw
				}
			}

			if len(fields) > 0 {
				acc.AddFields("nutanix_"+key, fields, tags)
			}
		}
	}

	return nil
}

// init registers a callback which creates a new nutanix input instance.
func init() {
	inputs.Add("nutanix", func() telegraf.Input {
		return &Nutanix{}
	})
}
