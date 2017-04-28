package apidCRUD

import (
	"database/sql"
)

// initDB opens the named database and returns a handle wrapper.
func initDB(dbName string) (dbType, error) {
	h, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return dbType{}, err
	}

	// assign the global db variable
	return dbType{handle: h}, nil
}
