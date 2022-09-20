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
  #    "host:port/service_name"
  #  or oracle net connect descriptor string like
  #    (DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=dbhost.example.com)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))
  connection_string = ""

  ## Database credentials
  username = ""
  password = ""

  ## Role, either SYSDBA, SYSASM, SYSOPER or empty
  role = ""

  ## Path to the Oracle Client library directory, optional.
  # Should be used if there is no LD_LIBRARY_PATH variable
  # or not possible to confugire it properly.
  client_lib_dir = ""

  ## Define the toml config where the sql queries are stored
  # Structure :
  # [[inputs.oracledb.query]]
  #   sqlquery string
  #   script string
  #   schema string
  #   tag_columns array of strings
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Memory Sort Ratio"
    # OracleDB sql query
    sqlquery = "select 'MEMORY_SORT' as METRIC_NAME, a.value as Sort_memory, b.value as Sort_disk, round(a.value/(a.value+b.value) * 100 ,2) as Memory_Sort_Ratio from v$sysstat a, v$sysstat b where a.name = 'sorts (memory)' and b.name = 'sorts (disk)'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Buffer Cache Hit Ratio"
    # OracleDB sql query
    sqlquery = "SELECT 'Buffer_Cache_Hit_Ratio' as METRIC_NAME, ROUND(((1-(SUM(DECODE(name, 'physical reads', value,0))/(SUM(DECODE(name, 'db block gets', value,0))+(SUM(DECODE(name, 'consistent gets', value, 0))))))*100),2) as Buffer_Cache_Hit_Ratio FROM V$SYSSTAT"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Cursor Cache Hit Ratio"
    # OracleDB sql query
    sqlquery = "select 'Cursor_Cache_Hit_Ratio' as METRIC_NAME, round((sum(decode(name, 'session cursor cache hits', value,0)) / (sum(decode(name, 'parse count (total)', value,0)) - sum(decode(name, 'parse count (hard)', value,0))))*100,2) as Cursor_Cache_Hit_Ratio from v$sysstat"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Library Cache Hit Ratio"
    # OracleDB sql query
    sqlquery = "SELECT 'Library_Cache_Hit_Ratio' as metric_name, round(sum(pinhits)/sum(pins)*100,2) as Library_Cache_Hit_Ratio FROM v$librarycache"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Shared_Pool_Free"
    # OracleDB sql query
    sqlquery = "SELECT 'Shared_Pool_Free' as metric_name, round((sum(decode(name,'free memory',bytes)) / sum(bytes))*100,2) as Shared_Pool_Free FROM  v$sgastat WHERE pool = 'shared pool'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Physical Reads"
    # OracleDB sql query
    sqlquery = "SELECT 'Physical_Reads' as metric_name, value as Physical_Reads FROM v$sysstat WHERE name='physical reads'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Physical Writes"
    # OracleDB sql query
    sqlquery = "SELECT 'Physical_Writes' as metric_name, value as Physical_Writes FROM v$sysstat WHERE name='physical writes'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Enqueue Timeouts"
    # OracleDB sql query
    sqlquery = "select 'Enqueue_Timeouts' as metric_name, sum(decode(name, 'enqueue timeouts', value,0)) as Enqueue_Timeouts from v$sysstat"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "GC CR block received"
    # OracleDB sql query
    sqlquery = "SELECT 'GC_CR_block_received' as metric_name, value as GC_CR_block_received FROM v$sysstat WHERE name = 'gc cr blocks received'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "GC blocks corrupt"
    # OracleDB sql query
    sqlquery = "SELECT 'GC_blocks_corrupt' as metric_name, value as GC_blocks_corrupt FROM v$sysstat WHERE name = 'gc blocks corrupt'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "GC blocks lost"
    # OracleDB sql query
    sqlquery = "SELECT 'GC_blocks_lost' as metric_name, value as GC_blocks_lost FROM v$sysstat WHERE name = 'gc blocks lost'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Active Sessions"
    # OracleDB sql query
    sqlquery = "SELECT 'Active_Sessions' as metric_name, count(*) as Active_Sessions FROM V$Session WHERE Status='ACTIVE' AND UserName IS NOT NULL"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Long Table Scans"
    # OracleDB sql query
    sqlquery = "SELECT 'Long_Table_Scans' as metric_name, value as Long_Table_Scans FROM v$sysstat WHERE name = 'table scans (long tables)'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "User Rollbacks"
    # OracleDB sql query
    sqlquery = "SELECT 'User_Rollbacks' as metric_name, value as User_Rollbacks FROM v$sysstat WHERE name = 'user rollbacks'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Sort Rows"
    # OracleDB sql query
    sqlquery = "SELECT 'Sort_Rows' as metric_name, value as Sort_Rows FROM v$sysstat WHERE name = 'sorts (rows)'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Memory Sorts"
    # OracleDB sql query
    sqlquery = "SELECT 'Memory_Sorts' as metric_name, value as Memory_Sorts FROM v$sysstat WHERE name = 'sorts (memory)'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Disk Sorts"
    # OracleDB sql query
    sqlquery = "SELECT 'Disk_Sorts' as metric_name, value as Disk_Sorts FROM v$sysstat WHERE name = 'sorts (disk)'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Session Limit"
    # OracleDB sql query
    sqlquery = "SELECT 'Session_Limit' as metric_name, to_number(LIMIT_VALUE) as Session_Limit FROM V$RESOURCE_LIMIT WHERE RESOURCE_NAME = 'sessions'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Session Limit Usage"
    # OracleDB sql query
    sqlquery = "SELECT 'Session_Limit_Usage' as metric_name, (CURRENT_UTILIZATION / LIMIT_VALUE) * 100 as Session_Limit_Usage FROM V$RESOURCE_LIMIT WHERE RESOURCE_NAME = 'sessions'"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Session Count"
    # OracleDB sql query
    sqlquery = "SELECT 'Session_Count' as metric_name, COUNT(*) as Session_Count FROM V$SESSION"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Pga used memory"
    # OracleDB sql query
    sqlquery = "SELECT 'PGA_Used_Memory_MB' as metric_name, ROUND(SUM(pga_used_mem)/(1024*1024),2) PGA_Used_Memory_MB FROM v$process"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Pga allocated memory"
    # OracleDB sql query
    sqlquery = "SELECT 'PGA_Allocated_Memory_MB' as metric_name, ROUND(SUM(pga_alloc_mem)/(1024*1024),2) PGA_Allocated_Memory_MB FROM v$process"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Pga freeable memory"
    # OracleDB sql query
    sqlquery = "SELECT 'PGA_Freeable_Memory_MB' as metric_name, ROUND(SUM(pga_freeable_mem)/(1024*1024),2) PGA_Freeable_Memory_MB FROM v$process"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Pga maximum memory"
    # OracleDB sql query
    sqlquery = "SELECT 'PGA_Maximum_Memory_MB' as metric_name, ROUND(SUM(pga_max_mem)/(1024*1024),2) PGA_Maximum_Memory_MB FROM v$process"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Tablespace"
    # OracleDB sql query
    sqlquery = "SELECT df.tablespace_name tablespace_name, round(df.maxbytes / (1024 * 1024), 2) ts_max_size, round(df.bytes / (1024 * 1024), 2) ts_curr_size, round((df.bytes - sum(fs.bytes)) / (1024 * 1024), 2) ts_used_size, round(sum(fs.bytes) / (1024 * 1024), 2) ts_free_size, round((df.bytes - sum(fs.bytes)) / (df.maxbytes) * 100, 2) ts_max_pct_used, round((df.bytes-sum(fs.bytes)) * 100 / df.bytes, 2) ts_pct_used, nvl(round(sum(fs.bytes) * 100 / df.bytes), 2) ts_pct_free FROM dba_free_space fs, (SELECT tablespace_name, sum(bytes) bytes, sum(decode(maxbytes, 0, bytes, maxbytes)) maxbytes FROM dba_data_files GROUP BY tablespace_name) df WHERE fs.tablespace_name (+) = df.tablespace_name GROUP BY df.tablespace_name, df.bytes, df.maxbytes UNION ALL SELECT df.tablespace_name tablespace_name, round(df.maxbytes / (1024 * 1024), 2) ts_max_size, round(df.bytes / (1024 * 1024), 2) ts_curr_size, round((df.bytes - sum(fs.bytes)) / (1024 * 1024), 2) ts_used_size, round(sum(fs.bytes) / (1024 * 1024), 2) ts_free_size, round((df.bytes - sum(fs.bytes)) / (df.maxbytes) * 100, 2) ts_max_pct_used, round((df.bytes-sum(fs.bytes)) * 100 / df.bytes, 2) ts_pct_used, nvl(round(sum(fs.bytes) * 100 / df.bytes), 2) ts_pct_free FROM (SELECT tablespace_name, bytes_used bytes FROM V$temp_space_header GROUP BY tablespace_name, bytes_free, bytes_used) fs, (SELECT tablespace_name, sum(bytes) bytes, sum(decode(maxbytes, 0, bytes, maxbytes)) maxbytes FROM dba_temp_files GROUP BY tablespace_name) df WHERE fs.tablespace_name (+) = df.tablespace_name GROUP BY df.tablespace_name, df.bytes, df.maxbytes ORDER BY 4 DESC"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["TABLESPACE_NAME"]
  [[inputs.oracledb.query]]
    # Query name, optional. Used in logging.
    name = "Total Table Scans"
    # OracleDB sql query
    sqlquery = "select 'Total_Table_Scans' as METRIC_NAME, sum(decode(name, 'table scans (long tables)', value,0)) + sum(decode(name, 'table scans (short tables)', value,0)) as Total_Table_Scans from v$sysstat"
    # The script option can be used to specify the .sql file path.
    # If script and sqlquery options specified at same time, sqlquery will be used.
    script = ""
    # Schema name. If provided, then ALTER SESSION SET CURRENT_SCHEMA query will be executed
    schema = ""
    # Query execution timeout, in seconds.
    timeout = 10
    # Array of column names, which would be stored as tags
    tag_columns = ["METRIC_NAME"]
```

### Example Output:

```bash
$ ./telegraf.exe --test --config=etc/telegraf_win_test.conf --input-filter oracledb
2021-05-14T06:28:50Z I! Starting Telegraf
2021-05-14T06:28:50Z E! Unable to open windows/telegraf.log (open windows/telegraf.log: The system cannot find the path specified.), using stderr
> oracledb,METRIC_NAME=MEMORY_SORT,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE MEMORY_SORT_RATIO=100,SORT_DISK=0,SORT_MEMORY=2652 1625123772000000000
> oracledb,METRIC_NAME=Buffer_Cache_Hit_Ratio,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE BUFFER_CACHE_HIT_RATIO=94.83 1625123772000000000
> oracledb,METRIC_NAME=Cursor_Cache_Hit_Ratio,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE CURSOR_CACHE_HIT_RATIO=342.31 1625123772000000000
> oracledb,METRIC_NAME=Library_Cache_Hit_Ratio,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE LIBRARY_CACHE_HIT_RATIO=78.11 1625123772000000000
> oracledb,METRIC_NAME=Shared_Pool_Free,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE SHARED_POOL_FREE=33.11 1625123772000000000
> oracledb,METRIC_NAME=Physical_Reads,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PHYSICAL_READS=4341 1625123772000000000
> oracledb,METRIC_NAME=Physical_Writes,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PHYSICAL_WRITES=52 1625123772000000000
> oracledb,METRIC_NAME=Enqueue_Timeouts,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE ENQUEUE_TIMEOUTS=0 1625123772000000000
> oracledb,METRIC_NAME=GC_CR_block_received,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE GC_CR_BLOCK_RECEIVED=0 1625123772000000000
> oracledb,METRIC_NAME=GC_blocks_corrupt,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE GC_BLOCKS_CORRUPT=0 1625123772000000000
> oracledb,METRIC_NAME=GC_blocks_lost,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE GC_BLOCKS_LOST=0 1625123772000000000
> oracledb,METRIC_NAME=Active_Sessions,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE ACTIVE_SESSIONS=1 1625123772000000000
> oracledb,METRIC_NAME=Long_Table_Scans,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE LONG_TABLE_SCANS=2 1625123772000000000
> oracledb,METRIC_NAME=User_Rollbacks,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE USER_ROLLBACKS=15 1625123773000000000
> oracledb,METRIC_NAME=Sort_Rows,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE SORT_ROWS=24210 1625123773000000000
> oracledb,METRIC_NAME=Memory_Sorts,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE MEMORY_SORTS=2689 1625123773000000000
> oracledb,METRIC_NAME=Disk_Sorts,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE DISK_SORTS=0 1625123773000000000
> oracledb,METRIC_NAME=Session_Limit,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE SESSION_LIMIT=172 1625123773000000000
> oracledb,METRIC_NAME=Session_Limit_Usage,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE SESSION_LIMIT_USAGE=18.023255813953487 1625123773000000000
> oracledb,METRIC_NAME=Session_Count,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE SESSION_COUNT=25 1625123773000000000
> oracledb,METRIC_NAME=PGA_Used_Memory_MB,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PGA_USED_MEMORY_MB=40.1 1625123773000000000
> oracledb,METRIC_NAME=PGA_Allocated_Memory_MB,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PGA_ALLOCATED_MEMORY_MB=58.64 1625123773000000000
> oracledb,METRIC_NAME=PGA_Freeable_Memory_MB,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PGA_FREEABLE_MEMORY_MB=4.63 1625123773000000000
> oracledb,METRIC_NAME=PGA_Maximum_Memory_MB,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE PGA_MAXIMUM_MEMORY_MB=69.45 1625123773000000000
> oracledb,TABLESPACE_NAME=SYSAUX,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TS_CURR_SIZE=690,TS_FREE_SIZE=45.94,TS_MAX_PCT_USED=1.97,TS_MAX_SIZE=32767.98,TS_PCT_FREE=7,TS_PCT_USED=93.34,TS_USED_SIZE=644.06 1625123773000000000
> oracledb,TABLESPACE_NAME=SYSTEM,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TS_CURR_SIZE=360,TS_FREE_SIZE=5.81,TS_MAX_PCT_USED=59.03,TS_MAX_SIZE=600,TS_PCT_FREE=2,TS_PCT_USED=98.39,TS_USED_SIZE=354.19 1625123773000000000
> oracledb,TABLESPACE_NAME=UNDOTBS1,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TS_CURR_SIZE=40,TS_FREE_SIZE=25,TS_MAX_PCT_USED=0.05,TS_MAX_SIZE=32767.98,TS_PCT_FREE=63,TS_PCT_USED=37.5,TS_USED_SIZE=15 1625123773000000000
> oracledb,TABLESPACE_NAME=USERS,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TS_CURR_SIZE=100,TS_FREE_SIZE=97.44,TS_MAX_PCT_USED=0.02,TS_MAX_SIZE=11264,TS_PCT_FREE=97,TS_PCT_USED=2.56,TS_USED_SIZE=2.56 1625123773000000000
> oracledb,TABLESPACE_NAME=TEMP,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TS_CURR_SIZE=32,TS_FREE_SIZE=32,TS_MAX_PCT_USED=0,TS_MAX_SIZE=32767.98,TS_PCT_FREE=100,TS_PCT_USED=0,TS_USED_SIZE=0 1625123773000000000
> oracledb,METRIC_NAME=Total_Table_Scans,db_unique_name=XE,host=S2100113,instance_name=XE,server_host=snet-oracle,service_name=XE TOTAL_TABLE_SCANS=347 1625123773000000000
```
