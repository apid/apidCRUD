package apidCRUD

// this file contains defs that are global to the apidCRUD plugin.
// do not put functions here.

import (
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
