//go:build !custom || inputs || inputs.http_response_ext

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/http_response_ext" // register plugin
