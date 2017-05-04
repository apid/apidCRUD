package apidCRUD

import (
	"fmt"
	"database/sql"
	// _ "github.com/proullon/ramsql/driver"
	_ "github.com/mattn/go-sqlite3"
)

const ut_DBTYPE = "sqlite3"
const ut_DBNAME = "unit-test.db"

// fakeInitDB initializes the fake DB used for testing.
// note that tests using this DB are not true unit tests,
// since they are implicitly filesystem dependent.
func fakeInitDB() (dbType, error) {
	
	dbHandle, err := sql.Open(ut_DBTYPE, ut_DBNAME)
	db = dbType{dbHandle}	// assigns to global
	if err == nil {
		err = createDbData(db)
	}
	return db, err
}

var cmds = []string {
	`drop table if exists tables`,
	`drop table if exists bundles`,
	`drop table if exists users`,
	`create table tables(name text unique not null)`,
	`create table bundles(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into bundles (name, uri) values ("b1", "http://localhost/~dfong/bundles/b1.zip")`,
	`insert into bundles (name, uri) values ("b2", "http://localhost/~dfong/bundles/b2.zip")`,
	`insert into bundles (name, uri) values ("b3", "http://localhost/~dfong/bundles/b3.zip")`,
	`insert into tables (name) values ("bundles")`,
	`insert into tables (name) values ("users")`,
	`insert into tables (name) values ("nothing")`,
	// xxx is an extra scratch table
	`drop table if exists xxx`,
	`create table xxx(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into xxx (name, uri) values ("x1", "url1")`,
	`insert into xxx (name, uri) values ("x2", "url2")`,
}

func createDbData(db dbType) error {
	dbh := db.handle
	for _, cmd := range cmds {
		_, err := dbh.Exec(cmd)
		// fmt.Printf("cmd=%s\n", cmd)
		if err != nil {
			fmt.Printf(`Exec error on "%s": [%s]\n`, cmd, err)
			return err
		}
	}
	return nil
}
