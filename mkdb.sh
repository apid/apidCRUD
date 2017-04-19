#! /bin/bash
#	mkdb

# zap the file
DBFILE=apidCRUD.db
cp /dev/null "$DBFILE"

# create the bundles table, the users table, the nothing table.
sqlite3 "$DBFILE" << EOF
create table bundles(id integer not null primary key autoincrement,
name text not null,
uri text not null);
create table users (id integer not null primary key autoincrement,
name text not null);
insert into users (name) values ("djfong");
create table nothing(id integer not null primary key autoincrement, name text not null);
.quit
EOF

# create the tables table
echo ""
sqlite3 "$DBFILE" <<EOF
create table tables (name text unique not null);
insert into tables (name) values ("bundles");
insert into tables (name) values ("users");
insert into tables (name) values ("nothing");
.quit
EOF

# dump the bundles table.
echo ""
sqlite3 "$DBFILE" << EOF
select * from bundles;
.quit
EOF
