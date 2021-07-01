# Oracle plugin

This oracle plugin provides metrics for your oracle database. It has been designed
to parse SQL queries in the plugin section of your `telegraf.conf`. Plugin requires
[Oracle Instant Client Library](https://www.oracle.com/database/technologies/instant-client/downloads.html) at run time.
Either Basic or Basic Light package is sufficient. Library path can be provided via
standard OS environment variables - `LD_LIBRARY_PATH`, `DYLD_LIBRARY_PATH` and `PATH` for
Linux, MacOS and Windows respectively, or via `client_lib_dir` plugin parameter for Windows and MacOS.

### Configuration example

```toml
[[inputs.oracledb]]
  ## Connection string, e.g. easy connect string like
  connection_string = "192.168.56.105:1521/xe"
  #  or oracle net connect descriptor string like
  #    (DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=dbhost.example.com)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))
  # connection_string = "(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=192.168.56.105)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=XE)))"
  ## Database credentials
  username = "system"
  password = "oracle"
  ## Role, either SYSDBA, SYSASM, SYSOPER or empty
  role = ""
  ## Path to the Oracle Client library directory, optional.
  # Should be used if there is no LD_LIBRARY_PATH variable
  # or not possible to confugire it properly.
  client_lib_dir = "C:/Oracle/instantclient_19_11"
  ## Define the toml config where the sql queries are stored
  # Structure :
  # [[inputs.oracledb.query]]
  #   sqlquery string
  #   script string
  #   schema string
  #   tag_columns array of strings
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = ""
    # OracleDB sql query
    sqlquery = "select n.wait_class as wait_class, round(m.time_waited/m.INTSIZE_CSEC,3) wait_value from   v$waitclassmetric  m, v$system_wait_class n where m.wait_class_id=n.wait_class_id and n.wait_class != 'Idle' union select  'CPU', round(value/100,3) wait_value from v$sysmetric where metric_name='CPU Usage Per Sec' and group_id=2 union select 'CPU_OS', round((prcnt.busy*parameter.cpu_count)/100,3) - aas.cpu from ( select value busy from v$sysmetric where metric_name='Host CPU Utilization (%)' and group_id=2 ) prcnt, ( select value cpu_count from v$parameter where name='cpu_count' )  parameter, ( select  'CPU', round(value/100,3) cpu from v$sysmetric where metric_name='CPU Usage Per Sec' and group_id=2) aas"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["WAIT_CLASS"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = ""
    # OracleDB sql query
    sqlquery = "select n.wait_class wait_class, n.name wait_name, m.wait_count cnt, round(10*m.time_waited/nullif(m.wait_count,0),3) avgms from v$eventmetric m, v$event_name n where m.event_id=n.event_id and n.wait_class <> 'Idle' and m.wait_count > 0 order by 1"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["WAIT_CLASS"]
```

### Example Output:

```bash
$ ./telegraf.exe --test --config=etc/telegraf_win_test.conf --input-filter oracledb
2021-05-14T06:28:50Z I! Starting Telegraf
2021-05-14T06:28:50Z E! Unable to open windows/telegraf.log (open windows/telegraf.log: The system cannot find the path specified.), using stderr
> oracledb,WAIT_CLASS=Application,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=CPU,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0.026 1620973731000000000
> oracledb,WAIT_CLASS=CPU_OS,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0.008 1620973731000000000
> oracledb,WAIT_CLASS=Commit,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Concurrency,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Configuration,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Network,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Other,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Scheduler,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=System\ I/O,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0.02 1620973731000000000
> oracledb,WAIT_CLASS=User\ I/O,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE WAIT_VALUE=0 1620973731000000000
> oracledb,WAIT_CLASS=Network,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE AVGMS=0.002,CNT=9,WAIT_NAME="SQL*Net message to client" 1620973732000000000
> oracledb,WAIT_CLASS=Other,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE AVGMS=0.004,CNT=1,WAIT_NAME="asynch descriptor resize" 1620973732000000000
> oracledb,WAIT_CLASS=System\ I/O,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE AVGMS=0.012,CNT=48,WAIT_NAME="control file sequential read" 1620973732000000000
> oracledb,WAIT_CLASS=System\ I/O,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE AVGMS=60.002,CNT=25,WAIT_NAME="control file parallel write" 1620973732000000000
> oracledb,WAIT_CLASS=User\ I/O,db_unique_name=XE,host=jack-office-host,instance_name=XE,server_host=a9d86139a205,service_name=XE AVGMS=0.061,CNT=1,WAIT_NAME="Disk file operations I/O"
```
