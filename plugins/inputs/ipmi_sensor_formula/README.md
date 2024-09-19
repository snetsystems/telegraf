# IPMI Sensor Extension Input Plugin

Get bare metal metrics using the command line utility
[`ipmitool`](https://github.com/ipmitool/ipmitool).

If no servers are specified, the plugin will query the local machine sensor
stats via the following command:

```sh
ipmitool sdr
```

or with the version 2 schema:

```sh
ipmitool sdr elist
```

When one or more servers are specified, the plugin will use the following
command to collect remote host sensor stats:

```sh
ipmitool -I lan -H SERVER -U USERID -P PASSW0RD sdr
```

Any of the following parameters will be added to the aforementioned query if
they're configured:

```sh
-y hex_key -L privilege
```

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: ../../../docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Read metrics from the bare metal servers via IPMI
[[inputs.ipmi_sensor_formula]]
  ## Specify the path to the ipmitool executable
  # path = "/usr/bin/ipmitool"

  ## Use sudo
  ## Setting 'use_sudo' to true will make use of sudo to run ipmitool.
  ## Sudo must be configured to allow the telegraf user to run ipmitool
  ## without a password.
  # use_sudo = false

  ## Specify one or more servers via a connection string
  ## Configure tags, extra tags, and formulas for custom calculations
  ##
  ## 'tags' are user-defined tags applied to the collected metrics.
  ## Format: tags = {"tag_key": "tag_value", "tag_key2": "tag_value2", ...}
  ## e.g.
  ##    tags = {"model": "dell", "vender": "dell"}
  ##
  ## 'extra_tags' are additional tags applied to the metrics collected from each server.
  ## Format: extra_tags = ["extratag1", "extratag2", ...]
  ## e.g.
  ##    extra_tags = ["hostname", "line", "rack", "pdu"]
  ##
  ## 'formula' defines custom calculations based on IPMI sensor data.
  ## Format: formula = {"calc_tag": "calculation_expression", ...}
  ## e.g.
  ##    formula = {"calc_psu": "[current] * [voltage]"}
  
  ## Session privilege level
  ## Choose from: CALLBACK, USER, OPERATOR, ADMINISTRATOR
  # privilege = "ADMINISTRATOR"

  ## Timeout
  ## Timeout for the ipmitool command to complete.
  # timeout = "20s"

  ## Metric schema version
  ## See the plugin readme for more information on schema versioning.
  # metric_version = 1

  ## Sensors to collect
  ## Choose from:
  ##   * sdr: default, collects sensor data records
  ##   * chassis_power_status: collects the power status of the chassis
  ##   * dcmi_power_reading: collects the power readings from the Data Center Management Interface
  # sensors = ["sdr"]

  ## Hex key
  ## Optionally provide the hex key for the IMPI connection.
  # hex_key = ""

  ## Cache
  ## If ipmitool should use a cache
  ## Using a cache can speed up collection times depending on your device.
  # use_cache = false

  ## Path to the ipmitools cache file (defaults to OS temp dir)
  ## The provided path must exist and must be writable
  # cache_path = ""

## Server-specific configuration (you can define multiple servers)
[[inputs.ipmi_sensor_formula.servers]] 
  ## 'connection' is the connection string for the IPMI server.
  ## Format: [username[:password]@][protocol[(address)]]
  ## e.g.
  ##    "user:password@lanplus(192.168.1.1)"
  # connection = "user:password@lanplus(192.168.1.1)"  # Example: connection string format

[inputs.ipmi_sensor_formula.servers.extra_tags]
  ## Define additional tags for the server (optional)
  ## 'extra_tags' allows you to define tags specific to each server's collected metrics.
  ## Format: {"extratag_key": "extratag_value", ...}
  ## e.g.
  ##    {"hostname": "server1", "rack": "02", "pdu": {"pdu_1": ["calc_psu1"]}}

## To add more servers, you can define additional `servers` blocks.
# [[inputs.ipmi_sensor_formula.servers]]
# connection = "another_user:another_password@lanplus(192.168.1.2)"
# [inputs.ipmi_sensor_formula.servers.extra_tags]
# hostname = "server2"
# line = "B"
# rack = "03"
# pdu = {"pdu_1" = ["calc_psu1"], "pdu_2" = ["calc_psu2"]}
```

## Sensors

By default the plugin collects data via the `sdr` command and returns those
values. However, there are additonal sensor options that be call on:

- `chassis_power_status` - returns 0 or 1 depending on the output of
  `chassis power status`
- `dcmi_power_reading` - Returns the watt values from `dcmi power reading`

These sensor options are not affected by the metric version.

## Metrics

Version 1 schema:

- ipmi_sensor:
  - tags:
    - name
    - unit
    - host
    - server (only when retrieving stats from remote servers)
    - extra_tags (only when extra_tags entered in server IPMI info extra_tags)
  - fields:
    - status (int, 1=ok status_code/0=anything else)
    - value (float)

Version 2 schema:

- ipmi_sensor:
  - tags:
    - name
    - entity_id (can help uniquify duplicate names)
    - status_code (two letter code from IPMI documentation)
    - status_desc (extended status description field)
    - unit (only on analog values)
    - host
    - server (only when retrieving stats from remote)
    - extra_tags (only when extra_tags entered in server IPMI info extra_tags)
  - fields:
    - value (float)

### Permissions

When gathering from the local system, Telegraf will need permission to the
ipmi device node. When using udev you can create the device node giving
`rw` permissions to the `telegraf` user by adding the following rule to
`/etc/udev/rules.d/52-telegraf-ipmi.rules`:

```sh
KERNEL=="ipmi*", MODE="660", GROUP="telegraf"
```

Alternatively, it is possible to use sudo. You will need the following in your
telegraf config:

```toml
[[inputs.ipmi_sensor_formula]]
  use_sudo = true
```

You will also need to update your sudoers file:

```bash
$ visudo
# Add the following line:
Cmnd_Alias IPMITOOL = /usr/bin/ipmitool *
telegraf  ALL=(root) NOPASSWD: IPMITOOL
Defaults!IPMITOOL !logfile, !syslog, !pam_session
```

## Example Output

### Version 1 Schema

When retrieving stats from a remote server:

```shell
ipmi_sensor,server=10.20.2.203,name=uid_light value=0,status=1i 1517125513000000000
ipmi_sensor,server=10.20.2.203,name=sys._health_led status=1i,value=0 1517125513000000000
ipmi_sensor,server=10.20.2.203,name=power_supply_1,unit=watts status=1i,value=110 1517125513000000000
ipmi_sensor,server=10.20.2.203,name=power_supply_2,unit=watts status=1i,value=120 1517125513000000000
ipmi_sensor,server=10.20.2.203,name=power_supplies value=0,status=1i 1517125513000000000
ipmi_sensor,server=10.20.2.203,name=fan_1,unit=percent status=1i,value=43.12 1517125513000000000
```

When retrieving stats from a remote server(extra_tags={"hostname": "example_host", "pdu": {"pdu_1": ["calc_psu1","calc_psu2"]}} specified):

```shell
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=uid_light value=0,status=1i 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=sys._health_led status=1i,value=0 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=power_supply_1,unit=watts status=1i,value=110 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=power_supply_2,unit=watts status=1i,value=120 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=power_supplies value=0,status=1i 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=fan_1,unit=percent status=1i,value=43.12 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=calc_psu1,pdu=pdu_1 status=1i,value=128.4 1517125513000000000
ipmi_sensor,server=10.20.2.203,hostname=example_host,name=calc_psu2,pdu=pdu_1 status=1i,value=42.4 1517125513000000000
```

When retrieving stats from the local machine (no server specified):

```shell
ipmi_sensor,name=uid_light value=0,status=1i 1517125513000000000
ipmi_sensor,name=sys._health_led status=1i,value=0 1517125513000000000
ipmi_sensor,name=power_supply_1,unit=watts status=1i,value=110 1517125513000000000
ipmi_sensor,name=power_supply_2,unit=watts status=1i,value=120 1517125513000000000
ipmi_sensor,name=power_supplies value=0,status=1i 1517125513000000000
ipmi_sensor,name=fan_1,unit=percent status=1i,value=43.12 1517125513000000000
```

#### Version 2 Schema

When retrieving stats from the local machine (no server specified):

```shell
ipmi_sensor,name=uid_light,entity_id=23.1,status_code=ok,status_desc=ok value=0 1517125474000000000
ipmi_sensor,name=sys._health_led,entity_id=23.2,status_code=ok,status_desc=ok value=0 1517125474000000000
ipmi_sensor,entity_id=10.1,name=power_supply_1,status_code=ok,status_desc=presence_detected,unit=watts value=110 1517125474000000000
ipmi_sensor,name=power_supply_2,entity_id=10.2,status_code=ok,unit=watts,status_desc=presence_detected value=125 1517125474000000000
ipmi_sensor,name=power_supplies,entity_id=10.3,status_code=ok,status_desc=fully_redundant value=0 1517125474000000000
ipmi_sensor,entity_id=7.1,name=fan_1,status_code=ok,status_desc=transition_to_running,unit=percent value=43.12 1517125474000000000
```
