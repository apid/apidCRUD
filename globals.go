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

// dbName is the name of the database that is implicitly used in these APIs.
var dbName = "apidCRUD.db"

// DEBUG enables debugging printouts
var DEBUG = true

// db is our global database handle
var db dbType

// log is our global log variable
var log apid.LogService
