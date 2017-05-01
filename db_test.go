package apidCRUD

import (
	"fmt"
	"database/sql"
	_ "github.com/proullon/ramsql/driver"
)

func fakeInitDb() (dbType, error) {
	dbHandle, err := sql.Open("ramsql", "fakedbfortesting")
	db = dbType{dbHandle}
	if err == nil {
		err = createDbData(db)
	}
	return db, err
}

var cmds = []string {
	`create table bundles(id integer not null primary key autoincrement, name text not null, uri text not null)`,
	`insert into bundles (name, uri) values ("b1", "http://localhost/~dfong/bundles/b1.zip")`,
	`insert into bundles (name, uri) values ("b2", "http://localhost/~dfong/bundles/b2.zip")`,
	`insert into bundles (name, uri) values ("b3", "http://localhost/~dfong/bundles/b3.zip")`,
}

func createDbData(db dbType) error {
	dbh := db.handle
	for _, cmd := range cmds {
		_, err := dbh.Exec(cmd)
		// fmt.Printf("cmd=%s\n", cmd)
		if err != nil {
			fmt.Printf(`Exec error [%s]\n`, err)
			return err
		}
	}
	return nil
}
