//go:build !custom || inputs || inputs.ipmi_sensor_formula

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/ipmi_sensor_formula" // register plugin
