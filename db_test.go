package apidCRUD

import (
	"fmt"
	"os"
)

const ut_DBNAME = "unit-test.db"

// utInitDB initializes the fake DB used for testing.
// note that tests using this DB are not true unit tests,
// since they are implicitly filesystem dependent.
func utInitDB() {
	_ = os.Remove(ut_DBNAME)
	db, _ = initDB(dbName)
	_ = createDbData(db)
}

var cmds = []string {
	// create the special table _tables_
	`create table _tables_(name text unique not null, schema text)`,
	`insert into _tables_ (name,schema) values ("bundles", "xxx")`,
	`insert into _tables_ (name,schema) values ("users", "xxx")`,
	`insert into _tables_ (name,schema) values ("nothing", "xxx")`,

	// ordinary tables
	`create table bundles(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into bundles (name, uri) values ("b1", "http://localhost/~dfong/bundles/b1.zip")`,
	`insert into bundles (name, uri) values ("b2", "http://localhost/~dfong/bundles/b2.zip")`,
	`insert into bundles (name, uri) values ("b3", "http://localhost/~dfong/bundles/b3.zip")`,

	// xxx is an extra scratch table
	`create table xxx(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into xxx (name, uri) values ("x1", "url1")`,
	`insert into xxx (name, uri) values ("x2", "url2")`,
	`insert into xxx (name, uri) values ("x3", "url3")`,
	`insert into xxx (name, uri) values ("x4", "url4")`,
	`insert into xxx (name, uri) values ("x5", "url5")`,
	`insert into xxx (name, uri) values ("x6", "url6")`,

	// create table for testing behavior around maxRecs
	`create table toomany(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into toomany (name, uri) values ("x1", "url1")`,
	`insert into toomany (name, uri) values ("x2", "url2")`,
	`insert into toomany (name, uri) values ("x3", "url3")`,
	`insert into toomany (name, uri) values ("x4", "url4")`,
	`insert into toomany (name, uri) values ("x5", "url5")`,
	`insert into toomany (name, uri) values ("x6", "url6")`,
	`insert into toomany (name, uri) values ("x7", "url7")`,
	`insert into toomany (name, uri) values ("x8", "url8")`,
	`insert into toomany (name, uri) values ("x9", "url9")`,
	`insert into toomany (name, uri) values ("x10", "url10")`,
	`insert into toomany (name, uri) values ("x11", "url11")`,
	`insert into toomany (name, uri) values ("x12", "url12")`,
	`insert into toomany (name, uri) values ("x13", "url13")`,
	`insert into toomany (name, uri) values ("x14", "url14")`,
	`insert into toomany (name, uri) values ("x15", "url15")`,
	`insert into toomany (name, uri) values ("x16", "url16")`,
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
