package apidCRUD

import (
	"database/sql"
)

// initDB opens the named database and returns a handle wrapper.
func initDB(dbName string) (dbType, error) {
	h, err := sql.Open(dbDriver, dbName)
	return dbType{handle: h}, err
}
