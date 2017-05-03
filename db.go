package apidCRUD

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// initDB opens the named database and returns a handle wrapper.
func initDB(dbName string) (dbType, error) {
	h, err := sql.Open("sqlite3", dbName)
	return dbType{handle: h}, err
}
