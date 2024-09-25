//go:build !custom || inputs || inputs.vrops

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/vrops" // register plugin
