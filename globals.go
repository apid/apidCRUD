package apidCRUD

// this file contains defs that are global to the apidCRUD plugin.
// do not put functions here other than initConfig() and confGet().

import (
	"strconv"
	"net/http"
	"database/sql"
	"github.com/30x/apid-core"
)

// dbType is intended to encapsulate the database handle type.
type dbType struct {
	handle *sql.DB
}

// badStat is a convenience constant, the http status for a bad request.
const badStat = http.StatusBadRequest

// db is our global database handle
var db dbType

// log is our global log variable
var log apid.LogService

// the variables below are globals controlled by config file.
// in normal operation, the defaults will be overridden
// when initPlugin() calls initConfig().
//
// during unit tests, TestMain() calls initConfig() with
// a fake config reader to do the initialization.

// dbName is the name of the database that is implicitly used in these APIs.
var dbName = "apidCRUD.db"

// dbDriver is the name of the database driver to use
var dbDriver = "sqlite3"

// basePath is the prefix applied to paths in the API description table
var basePath = "/apid"

// maxRecs is the max number of results allowed in a bulk request.
var maxRecs = 1000

// tobleOfTables is the name of the internal table of table names/schemas
var  tableOfTables = "_tables_"

// getStringer is an interface that supports GetString().
// narrowed from apid.ConfigService.
type getStringer interface {
	GetString(vname string) string
}

// confGet() returns the config value of the named string,
// or if there is no configured value, the given default value.
func confGet(gsi getStringer, vname string, defval string) string {
	ret := gsi.GetString(vname)
	if ret == "" {
		return defval
	}
	return ret
}

// initConfig() sets up some global configuration parameters for this plugin.
func initConfig(gsi getStringer) {
	aMaxRecs := strconv.Itoa(maxRecs)

	// these are all global assignments!
	dbDriver = confGet(gsi, "apidCRUD_db_driver", dbDriver)
	dbName = confGet(gsi, "apidCRUD_db_name", dbName)
	basePath = confGet(gsi, "apidCRUD_base_path", basePath)
	maxRecs, _ = strconv.Atoi(			// nolint
		confGet(gsi, "apidCRUD_max_recs", aMaxRecs))
}
