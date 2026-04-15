//go:build !custom || inputs || inputs.diskio_ext

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/diskio_ext" // register plugin
