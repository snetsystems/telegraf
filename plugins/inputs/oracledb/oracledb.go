package oracledb

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/godror/godror"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const (
	defaultTimeout = 10
	pluginName     = "oracledb"
)

// OracleDB connection infos & operations
type OracleDB struct {
	ConnectionString string  `toml:"connection_string"`
	Username         string  `toml:"username"`
	Password         string  `toml:"password"`
	Role             string  `toml:"role"`
	ClientLibDir     string  `toml:"client_lib_dir"`
	Queries          []query `toml:"query"`

	DB  *sql.DB
	Log telegraf.Logger
}

type query struct {
	Name       string   `toml:"name"`
	Sqlquery   string   `toml:"sqlquery"`
	Script     string   `toml:"script"`
	Schema     string   `toml:"schema"`
	Timeout    int      `toml:"timeout"`
	TagColumns []string `toml:"tag_columns"`
}

var sampleConfig = `
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
`

// SampleConfig to implement PluginDescriber interface
func (o *OracleDB) SampleConfig() string {
	return sampleConfig
}

// Description to implement PluginDescriber interface
func (o *OracleDB) Description() string {
	return "Read metrics from one or many oracle database servers"
}

// Init to implement Initializer interface
func (o *OracleDB) Init() error {
	var err error

	if err = o.initDB(); err != nil {
		o.Log.Errorf("db: %s, couldn't init db due to error: %s", o.ConnectionString, err)
	}

	for i := range o.Queries {
		if o.Queries[i].Sqlquery == "" {
			o.Queries[i].Sqlquery, err = readScript(o.Queries[i].Script)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Gather to implement Input interface
func (o *OracleDB) Gather(acc telegraf.Accumulator) error {
	dbinfo, err := o.getDBInfo()
	if err != nil {
		return fmt.Errorf("db=%s couldn't obtain db info due to error: %s", o.ConnectionString, err)
	}

	for id, q := range o.Queries {
		queryID := strconv.Itoa(id)
		if q.Name != "" {
			queryID = q.Name
		}

		o.Log.Debugf("db=%s executing query %s: %s", o.ConnectionString, queryID, q.Sqlquery)

		if q.Schema != "" {
			if err := o.switchSchema(q.Schema); err != nil {
				o.Log.Errorf("db=%s skip query %s due to alter schema error: %s", o.ConnectionString, queryID, err)
				continue
			}
		}

		if q.Timeout <= 0 {
			q.Timeout = defaultTimeout
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(q.Timeout)*time.Second)
		defer cancel()

		rows, err := o.DB.QueryContext(ctx, q.Sqlquery)
		if err != nil {
			o.Log.Errorf("db=%s skip query %s due to error: %s", o.ConnectionString, queryID, err)
			continue
		}

		columns, err := rows.ColumnTypes()
		if err != nil {
			o.Log.Errorf("db=%s skip query %s due to obtaining column data error: %s", o.ConnectionString, queryID, err)
			_ = rows.Close()
			continue
		}

		for rows.Next() {
			if err := o.gatherRow(&q, rows, columns, dbinfo, acc); err != nil {
				o.Log.Errorf("db=%s skip row in query %s due to error: %s", o.ConnectionString, queryID, err)
				continue
			}
		}

		_ = rows.Close()
	}

	return nil
}

func (o *OracleDB) gatherRow(q *query, r *sql.Rows, ct []*sql.ColumnType, dbinfo *map[string]string, acc telegraf.Accumulator) error {
	var columnVars []interface{}
	columnMap := make(map[*sql.ColumnType]*interface{})
	for _, col := range ct {
		columnMap[col] = new(interface{})
		columnVars = append(columnVars, columnMap[col])
	}

	if err := r.Scan(columnVars...); err != nil {
		return err
	}

	tags := make(map[string]string)
	for tag, val := range *dbinfo {
		tags[tag] = val
	}
	fields := make(map[string]interface{})

COLUMN:
	for col, val := range columnMap {
		o.Log.Debugf("Column: %s = %s: %v", col.Name(), col.DatabaseTypeName(), *val)
		if *val == nil {
			o.Log.Debugf("skip column %s due to its value is nil", col.Name())
			continue
		}

		for _, tag := range q.TagColumns {
			if col.Name() != tag {
				continue
			}
			switch v := (*val).(type) {
			case string:
				tags[col.Name()] = v
			default:
				o.Log.Warnf("couldn't assign value %v to tag %s due to unknown type %T", *val, col.Name(), v)
			}
			continue COLUMN
		}

		switch col.DatabaseTypeName() {
		case "NUMBER":
			// godror exposes Oracle NUMBER type as string, so we have to do some parsing
			number := (*val).(string)
			result, err := convertOraNumberType(number)
			if err != nil {
				o.Log.Error(err)
				continue
			}
			fields[col.Name()] = result
		default:
			fields[col.Name()] = *val
		}
	}

	acc.AddFields(pluginName, fields, tags)
	return nil
}

// getDBInfo gathers db environment such as hostname and instance name to use as tags
func (o *OracleDB) getDBInfo() (*map[string]string, error) {
	var (
		serverHost   string
		serviceName  string
		instanceName string
		dbUniqueName string
	)

	query := "SELECT SYS_CONTEXT('USERENV', 'SERVER_HOST'), " +
		"SYS_CONTEXT('USERENV', 'SERVICE_NAME'), " +
		"SYS_CONTEXT('USERENV', 'INSTANCE_NAME'), " +
		"SYS_CONTEXT('USERENV', 'DB_UNIQUE_NAME') " +
		"FROM DUAL"

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout*time.Second)
	defer cancel()

	if err := o.DB.QueryRowContext(ctx, query).Scan(
		&serverHost,
		&serviceName,
		&instanceName,
		&dbUniqueName,
	); err != nil {
		return nil, err
	}

	return &map[string]string{
		"server_host":    serverHost,
		"service_name":   serviceName,
		"instance_name":  instanceName,
		"db_unique_name": dbUniqueName,
	}, nil
}

func (o *OracleDB) initDB() error {
	var params godror.ConnectionParams
	params.Username = o.Username
	params.Password = godror.NewPassword(o.Password)
	params.ConnectString = o.ConnectionString

	if o.ClientLibDir != "" {
		params.LibDir = o.ClientLibDir
	}

	switch strings.ToUpper(o.Role) {
	case "SYSDBA":
		params.IsSysDBA = true
	case "SYSOPER":
		params.IsSysOper = true
	case "SYSASM":
		params.IsSysASM = true
	}

	o.Log.Debugf("init db connection using parameters: %s", params.String())
	o.DB = sql.OpenDB(godror.NewConnector(params))
	if err := o.DB.Ping(); err != nil {
		return err
	}

	return nil
}

func (o *OracleDB) switchSchema(schema string) error {
	if _, err := o.DB.Exec("ALTER SESSION SET CURRENT_SCHEMA = " + schema); err != nil {
		return err
	}
	return nil
}

// convertOraNumberType converts string representation of Oracle's NUMBER to golang types.
func convertOraNumberType(number string) (interface{}, error) {
	// if number string contains dot(.), then it is float
	if strings.Contains(number, ".") {
		result, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert NUMBER %s to float: %s", number, err)
		}
		return result, nil
	}

	// if number string contains minus(-) sign, then it is integer
	if strings.HasPrefix(number, "-") {
		result, err := strconv.ParseInt(number, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert NUMBER %s to int64: %s", number, err)
		}
		return result, nil
	}

	// in other cases it is unsigned integer for maximum precision
	result, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert NUMBER %s to uint64: %s", number, err)
	}
	return result, nil
}

func readScript(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	query, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(query), err
}

func init() {
	inputs.Add(pluginName, func() telegraf.Input { return &OracleDB{} })
} 