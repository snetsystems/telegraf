//go:build !custom || inputs || inputs.nutanix

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/nutanix" // register plugin
