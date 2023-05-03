//go:generate ../../../tools/readme_config_includer/generator
// Package vrops implements an vrops input plugin for Telegraf
package vrops

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/common/tls"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers/graphite"
)

// DO NOT REMOVE THE NEXT TWO LINES! This is required to embed the sampleConfig data.
//go:embed sample.conf
var sampleConfig string

// vROps is the main structure associated with a collection instance.
type vROps struct {
	URL                            string                 `toml:"url"`
	Username                       string                 `toml:"username"`
	Password                       string                 `toml:"password"`
	EnabledServices                []string               `toml:"enabled_services"`
	ProjectStatKey                 []string               `toml:"project_stat_key"`
	VMStatKey                      []string               `toml:"vm_stat_key"`
	TanzuProjectStatKey            []string               `toml:"tanzu_project_stat_key"`
	TanzuVMStatKey                 []string               `toml:"tanzu_vm_stat_key"`
	TanzuProjectResourceKind       string                 `toml:"tanzu_project_resource_kind"`
	TanzuProjectPropertyConditions map[string]interface{} `toml:"tanzu_project_property_conditions"`
	TanzuVMResourceKind            string                 `toml:"tanzu_vm_resource_kind"`
	TanzuVMPropertyConditions      map[string]interface{} `toml:"tanzu_vm_property_conditions"`
	// bucket -> influx templates
	Templates []string
	// MetricSeparator is the separator between parts of the metric name.
	MetricSeparator string

	tls.ClientConfig

	vRealizeOpsToken string

	ResponseTimeout config.Duration

	client *http.Client
	Log    telegraf.Logger

	graphiteParser *graphite.GraphiteParser
}

func (*vROps) SampleConfig() string {
	return sampleConfig
}

func (o *vROps) Init() error {
	if o.URL == "" {
		return fmt.Errorf("url cannot be empty")
	}

	if o.Username == "" || o.Password == "" {
		return fmt.Errorf("username or password can not be empty string")
	}

	client, err := o.createHTTPClient()
	if err != nil {
		return err
	}

	o.client = client

	err = o.getvRealizeOpsToken()
	if err != nil {
		return err
	}

	return nil
}

// Gather gathers resources from the vROps API and accumulates metrics.  This
// implements the Input interface.
func (o *vROps) Gather(acc telegraf.Accumulator) error {
	// Gather resources.  Note service harvesting must come first as the other
	// gatherers are dependant on this information.
	gatherers := map[string]func(telegraf.Accumulator) error{
		"projects":      o.gatherProject,
		"vms":           o.gatherVMs,
		"tanzuProjects": o.gatherTanzuProject,
		"tanzuVMs":      o.gatherTanzuVMs,
	}

	callDuration := map[string]interface{}{}
	for _, service := range o.EnabledServices {
		// As Services are already gathered in Init(), using this to accumulate them.
		start := time.Now()
		gatherer := gatherers[service]
		if err := gatherer(acc); err != nil {
			acc.AddError(fmt.Errorf("failed to get resource %q %v", service, err))
		}
		callDuration[service] = time.Since(start).Nanoseconds()
	}

	return nil
}

func (o *vROps) parseName(bucket string) (name string, field string, tags map[string]string) {
	p := o.graphiteParser
	var err error

	if p == nil || o.graphiteParser.Separator != o.MetricSeparator {
		var metricSeparator string
		if o.MetricSeparator == "" {
			metricSeparator = "_"
		} else {
			metricSeparator = o.MetricSeparator
		}
		p, err = graphite.NewGraphiteParser(metricSeparator, o.Templates[:], nil)
		o.graphiteParser = p
	}

	if err == nil {
		pattern := `\.{2,}`
		re := regexp.MustCompile(pattern)

		name, tags, field, _ = p.ApplyTemplate(re.ReplaceAllString(bucket, "."))
	}

	return name, field, tags
}

func (o *vROps) gatherProject(acc telegraf.Accumulator) error {
	for {
		projects, err := o.getProjects()
		if err != nil {
			return err
		}

		projectIds := make([]string, 0, len(projects))
		for project := range projects {
			projectIds = append(projectIds, project)
		}

		statReqbody := &statBody{
			ResourceID:  projectIds,
			StatKey:     nil,
			CurrentOnly: true,
			MaxSamples:  1,
		}

		statsURL, err := url.Parse(o.URL + apiURL + resourceStats)
		if err != nil {
			return err
		}

		statsReq, err := json.Marshal(statReqbody)
		if err != nil {
			return nil
		}

		reqStats, err := http.NewRequest("POST", statsURL.String(), bytes.NewBuffer(statsReq))
		if err != nil {
			return err
		}

		o.setRequestHeaders(reqStats)

		resStats, err := o.client.Do(reqStats)
		if err != nil {
			return err
		}
		defer resStats.Body.Close()

		if resStats.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return err
			}
			continue
		}

		resourceStats, err := io.ReadAll(resStats.Body)
		if err != nil {
			return err
		}

		var stats statResource

		if err := json.Unmarshal(resourceStats, &stats); err != nil {
			return err
		}

		for _, resourceIds := range stats.Values {
			for _, stat := range resourceIds.StatList.Stat {

				re := regexp.MustCompile(`[:/|\s]`)
				bucket := re.ReplaceAllString(strings.ReplaceAll(stat.StatKey.Key, " ", ""), ".")

				name, field, tags := o.parseName(bucket)

				measurement := strings.ToLower(name)
				tags["tenant_id"] = resourceIds.ResourceID
				tags["tenant_name"] = strings.ReplaceAll(projects[resourceIds.ResourceID], " ", "_")

				fields := make(map[string]interface{})
				if stat.Data != nil {
					fields[field] = stat.Data[0]
				}

				if stat.Values != nil {
					fields[field] = stat.Values[0]
				}

				acc.AddFields(strings.Join([]string{"vrops_project", measurement}, "_"), fields, tags)

			}
		}

		return nil
	}
}

func (o *vROps) gatherVMs(acc telegraf.Accumulator) error {
	for {
		projects, err := o.getProjects()
		if err != nil {
			return err
		}

		deployments, err := o.getDeployments(projects)
		if err != nil {
			return err
		}

		vms, err := o.getVMs(deployments)
		if err != nil {
			return err
		}

		vmIds := make([]string, 0, len(vms))
		for vm := range vms {
			vmIds = append(vmIds, vm)
		}

		var statKey []string
		if o.VMStatKey != nil {
			statKey = o.VMStatKey
		} else {
			statKey = vmStatKey[:]
		}

		statReqbody := &statBody{
			ResourceID:  vmIds,
			StatKey:     statKey,
			CurrentOnly: true,
			MaxSamples:  1,
		}

		statsURL, err := url.Parse(o.URL + apiURL + resourceStats)
		if err != nil {
			return err
		}

		statsReq, err := json.Marshal(statReqbody)
		if err != nil {
			return nil
		}

		reqStats, err := http.NewRequest("POST", statsURL.String(), bytes.NewBuffer(statsReq))
		if err != nil {
			return err
		}

		o.setRequestHeaders(reqStats)

		resStats, err := o.client.Do(reqStats)
		if err != nil {
			return err
		}
		defer resStats.Body.Close()

		if resStats.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return err
			}
			continue
		}

		resourceStats, err := io.ReadAll(resStats.Body)
		if err != nil {
			return err
		}

		var stats statResource

		if err := json.Unmarshal(resourceStats, &stats); err != nil {
			return err
		}

		for _, resourceIds := range stats.Values {
			for _, stat := range resourceIds.StatList.Stat {

				re := regexp.MustCompile(`[:/|\s]`)
				bucket := re.ReplaceAllString(strings.ReplaceAll(stat.StatKey.Key, " ", ""), ".")

				name, field, tags := o.parseName(bucket)

				measurement := strings.ToLower(name)
				tags["tenant_id"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["projectId"], " ", "_")
				tags["tenant_name"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["projectName"], " ", "_")
				tags["deployment_id"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["deploymentId"], " ", "_")
				tags["deployment_name"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["deploymentName"], " ", "_")
				tags["vm_id"] = resourceIds.ResourceID
				tags["vm_name"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["name"], " ", "_")

				fields := make(map[string]interface{})
				if stat.Data != nil {
					fields[field] = stat.Data[0]
				}

				if stat.Values != nil {
					fields[field] = stat.Values[0]
				}

				acc.AddFields(strings.Join([]string{"vrops_vm", measurement}, "_"), fields, tags)

			}
		}

		return nil
	}
}

func (o *vROps) gatherTanzuProject(acc telegraf.Accumulator) error {
	for {
		projects, err := o.getTanzuProjects()
		if err != nil {
			return err
		}

		projectIds := make([]string, 0, len(projects))
		for project := range projects {
			projectIds = append(projectIds, project)
		}

		statReqbody := &statBody{
			ResourceID:  projectIds,
			StatKey:     nil,
			CurrentOnly: true,
			MaxSamples:  1,
		}

		statsURL, err := url.Parse(o.URL + apiURL + resourceStats)
		if err != nil {
			return err
		}

		statsReq, err := json.Marshal(statReqbody)
		if err != nil {
			return nil
		}

		reqStats, err := http.NewRequest("POST", statsURL.String(), bytes.NewBuffer(statsReq))
		if err != nil {
			return err
		}

		o.setRequestHeaders(reqStats)

		resStats, err := o.client.Do(reqStats)
		if err != nil {
			return err
		}
		defer resStats.Body.Close()

		if resStats.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return err
			}
			continue
		}

		resourceStats, err := io.ReadAll(resStats.Body)
		if err != nil {
			return err
		}

		var stats statResource

		if err := json.Unmarshal(resourceStats, &stats); err != nil {
			return err
		}

		for _, resourceIds := range stats.Values {
			for _, stat := range resourceIds.StatList.Stat {

				re := regexp.MustCompile(`[:/|\s]`)
				bucket := re.ReplaceAllString(strings.ReplaceAll(stat.StatKey.Key, " ", ""), ".")

				name, field, tags := o.parseName(bucket)

				measurement := strings.ToLower(name)
				tags["tenant_id"] = resourceIds.ResourceID
				tags["tenant_name"] = strings.ReplaceAll(projects[resourceIds.ResourceID], " ", "_")

				fields := make(map[string]interface{})
				if stat.Data != nil {
					fields[field] = stat.Data[0]
				}

				if stat.Values != nil {
					fields[field] = stat.Values[0]
				}

				acc.AddFields(strings.Join([]string{"vrops_tanzu_project", measurement}, "_"), fields, tags)

			}
		}

		return nil
	}
}

func (o *vROps) gatherTanzuVMs(acc telegraf.Accumulator) error {
	for {
		projects, err := o.getTanzuProjects()
		if err != nil {
			return err
		}

		vms, err := o.getTanzuVMs(projects)
		if err != nil {
			return err
		}

		vmIds := make([]string, 0, len(vms))
		for vm := range vms {
			vmIds = append(vmIds, vm)
		}

		var statKey []string
		if o.TanzuVMStatKey != nil {
			statKey = o.TanzuVMStatKey
		} else {
			statKey = tanzuVMStatKey[:]
		}

		statReqbody := &statBody{
			ResourceID:  vmIds,
			StatKey:     statKey,
			CurrentOnly: true,
			MaxSamples:  1,
		}

		statsURL, err := url.Parse(o.URL + apiURL + resourceStats)
		if err != nil {
			return err
		}

		statsReq, err := json.Marshal(statReqbody)
		if err != nil {
			return nil
		}

		reqStats, err := http.NewRequest("POST", statsURL.String(), bytes.NewBuffer(statsReq))
		if err != nil {
			return err
		}

		o.setRequestHeaders(reqStats)

		resStats, err := o.client.Do(reqStats)
		if err != nil {
			return err
		}
		defer resStats.Body.Close()

		if resStats.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return err
			}
			continue
		}

		resourceStats, err := io.ReadAll(resStats.Body)
		if err != nil {
			return err
		}

		var stats statResource

		if err := json.Unmarshal(resourceStats, &stats); err != nil {
			return err
		}

		for _, resourceIds := range stats.Values {
			for _, stat := range resourceIds.StatList.Stat {

				re := regexp.MustCompile(`[:/|\s]`)
				bucket := re.ReplaceAllString(strings.ReplaceAll(stat.StatKey.Key, " ", ""), ".")

				name, field, tags := o.parseName(bucket)

				measurement := strings.ToLower(name)
				tags["tenant_id"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["projectId"], " ", "_")
				tags["tenant_name"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["projectName"], " ", "_")
				tags["vm_id"] = resourceIds.ResourceID
				tags["vm_name"] = strings.ReplaceAll(vms[resourceIds.ResourceID]["name"], " ", "_")

				fields := make(map[string]interface{})
				if stat.Data != nil {
					fields[field] = stat.Data[0]
				}

				if stat.Values != nil {
					fields[field] = stat.Values[0]
				}

				acc.AddFields(strings.Join([]string{"vrops_tanzu_vm", measurement}, "_"), fields, tags)

			}
		}

		return nil
	}
}

func (o *vROps) getProjects() (map[string]string, error) {
	for {
		resourceURL, err := url.Parse(o.URL + apiURL + resourceInfo)
		if err != nil {
			return nil, err
		}

		requestbody := &resourceBody{
			ResourceKind:   []string{"Project"},
			ResourceState:  []string{"STARTED"},
			ResourceStatus: []string{"DATA_RECEIVING"},
		}

		body, err := json.Marshal(requestbody)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", resourceURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		o.setRequestHeaders(request)

		response, err := o.client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return nil, err
			}
			continue
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resources resources

		if err := json.Unmarshal(resBody, &resources); err != nil {
			return nil, err
		}

		projects := make(map[string]string)

		for _, resource := range resources.ResourceList {
			projects[resource.Identifier] = resource.ResourceKey.Name
		}

		return projects, nil
	}
}

func (o *vROps) getDeployments(projects map[string]string) (map[string]map[string]string, error) {
	for {
		resourceQuery := &resourceBody{
			ResourceKind:   []string{"Deployment"},
			ResourceState:  []string{"STARTED"},
			ResourceStatus: []string{"DATA_RECEIVING"},
		}

		projectIds := make([]string, 0, len(projects))
		for project := range projects {
			projectIds = append(projectIds, project)
		}

		requestbody := &relationshipBody{
			RelationshipType: "CHILD",
			ResourceIds:      projectIds,
			ResourceQuery:    *resourceQuery,
		}

		bulkURL, err := url.Parse(o.URL + apiURL + resourceBulk)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(requestbody)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", bulkURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		o.setRequestHeaders(request)

		response, err := o.client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return nil, err
			}
			continue
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resources bulkResource

		if err := json.Unmarshal(resBody, &resources); err != nil {
			return nil, err
		}

		deployments := make(map[string]map[string]string)

		for _, resource := range resources.ResourcesRelations {
			deployments[resource.Resource.Identifier] = map[string]string{
				"name":        resource.Resource.ResourceKey.Name,
				"projectId":   resource.Resource.RelatedResources[0],
				"projectName": projects[resource.Resource.RelatedResources[0]]}
		}

		return deployments, nil
	}
}

func (o *vROps) getVMs(deployments map[string]map[string]string) (map[string]map[string]string, error) {
	for {
		resourceQuery := &resourceBody{
			ResourceKind:   []string{"VirtualMachine"},
			ResourceState:  []string{"STARTED"},
			ResourceStatus: []string{"DATA_RECEIVING"},
		}

		deploymentIds := make([]string, 0, len(deployments))
		for deployment := range deployments {
			deploymentIds = append(deploymentIds, deployment)
		}

		requestBody := &relationshipBody{
			RelationshipType: "CHILD",
			ResourceIds:      deploymentIds,
			ResourceQuery:    *resourceQuery,
		}

		bulkURL, err := url.Parse(o.URL + apiURL + resourceBulk)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", bulkURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		o.setRequestHeaders(request)

		response, err := o.client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return nil, err
			}
			continue
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resources bulkResource

		if err := json.Unmarshal(resBody, &resources); err != nil {
			return nil, err
		}

		vms := make(map[string]map[string]string)

		for _, resource := range resources.ResourcesRelations {
			vms[resource.Resource.Identifier] = map[string]string{
				"name":           resource.Resource.ResourceKey.Name,
				"deploymentId":   resource.Resource.RelatedResources[0],
				"deploymentName": deployments[resource.Resource.RelatedResources[0]]["name"],
				"projectId":      deployments[resource.Resource.RelatedResources[0]]["projectId"],
				"projectName":    deployments[resource.Resource.RelatedResources[0]]["projectName"]}
		}

		return vms, nil
	}
}

func (o *vROps) getTanzuProjects() (map[string]string, error) {
	for {
		resourceURL, err := url.Parse(o.URL + apiURL + resourceInfo)
		if err != nil {
			return nil, err
		}

		var resourceKind string
		if o.TanzuProjectResourceKind != "" {
			resourceKind = o.TanzuProjectResourceKind
		} else {
			resourceKind = tanzuProjectResourceKind
		}

		propertyConditions := make(map[string]interface{})
		if o.TanzuProjectPropertyConditions != nil {
			propertyConditions = o.TanzuProjectPropertyConditions
		} else {
			propertyConditions = tanzuProjectPropertyConditions
		}

		requestbody := &resourceBody{
			ResourceKind:       []string{resourceKind},
			ResourceState:      []string{"STARTED"},
			ResourceStatus:     []string{"DATA_RECEIVING"},
			PropertyConditions: propertyConditions,
		}

		body, err := json.Marshal(requestbody)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", resourceURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		o.setRequestHeaders(request)

		response, err := o.client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return nil, err
			}
			continue
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resources resources

		if err := json.Unmarshal(resBody, &resources); err != nil {
			return nil, err
		}

		projects := make(map[string]string)

		for _, resource := range resources.ResourceList {
			projects[resource.Identifier] = resource.ResourceKey.Name
		}

		return projects, nil
	}
}

func (o *vROps) getTanzuVMs(projects map[string]string) (map[string]map[string]string, error) {
	for {
		var resourceKind string
		if o.TanzuVMResourceKind != "" {
			resourceKind = o.TanzuVMResourceKind
		} else {
			resourceKind = tanzuVMResourceKind
		}

		propertyConditions := make(map[string]interface{})
		if o.TanzuVMPropertyConditions != nil {
			propertyConditions = o.TanzuVMPropertyConditions
		} else {
			propertyConditions = tanzuVMPropertyConditions
		}

		resourceQuery := &resourceBody{
			ResourceKind:       []string{resourceKind},
			ResourceState:      []string{"STARTED"},
			ResourceStatus:     []string{"DATA_RECEIVING"},
			PropertyConditions: propertyConditions,
		}

		projectIds := make([]string, 0, len(projects))
		for project := range projects {
			projectIds = append(projectIds, project)
		}

		requestBody := &relationshipBody{
			RelationshipType: "CHILD",
			ResourceIds:      projectIds,
			ResourceQuery:    *resourceQuery,
		}

		bulkURL, err := url.Parse(o.URL + apiURL + resourceBulk)
		if err != nil {
			return nil, err
		}

		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", bulkURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		o.setRequestHeaders(request)

		response, err := o.client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 401 {
			err := o.getvRealizeOpsToken()
			if err != nil {
				return nil, err
			}
			continue
		}

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resources bulkResource

		if err := json.Unmarshal(resBody, &resources); err != nil {
			return nil, err
		}

		vms := make(map[string]map[string]string)

		for _, resource := range resources.ResourcesRelations {
			vms[resource.Resource.Identifier] = map[string]string{
				"name":        resource.Resource.ResourceKey.Name,
				"projectId":   resource.Resource.RelatedResources[0],
				"projectName": projects[resource.Resource.RelatedResources[0]]}
		}

		return vms, nil
	}
}

func (o *vROps) getvRealizeOpsToken() error {
	resourceURL, err := url.Parse(o.URL + apiURL + authToken)
	if err != nil {
		return err
	}

	requestbody := &tokenBody{
		Username: o.Username,
		Password: o.Password,
	}

	body, err := json.Marshal(requestbody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", resourceURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := o.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 401 {
		return fmt.Errorf("unable to authenticate vROps user %v", response.StatusCode)
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var vROpsToken vROpsToken

	if err := json.Unmarshal(resBody, &vROpsToken); err != nil {
		return err
	}

	o.vRealizeOpsToken = vROpsToken.Token

	return nil
}

func (o *vROps) createHTTPClient() (*http.Client, error) {
	tlsCfg, err := o.ClientConfig.TLSConfig()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
		Timeout: time.Duration(o.ResponseTimeout),
	}

	return client, nil
}

// Set common request headers
func (o *vROps) setRequestHeaders(req *http.Request) {
	bearer := "vRealizeOpsToken " + strings.Trim(o.vRealizeOpsToken, "\n")

	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

// init registers a callback which creates a new vROps input instance.
func init() {
	inputs.Add("vrops", func() telegraf.Input {
		return &vROps{}
	})
}
