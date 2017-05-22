#! /bin/bash
#	mkdb

# ----- start of mainline code
PROGDIR=$(cd "$(dirname "$0")" && /bin/pwd)
. "$PROGDIR/tester-env.sh" || exit 1
. "$PROGDIR/test-common.sh" || exit 1

# zap the file
cp /dev/null "$DBFILE"

# create the bundles table, the users table, the nothing table, the file table.
sqlite3 "$DBFILE" << EOF
create table bundles(id integer not null primary key autoincrement,
name text not null,
uri text not null);
create table users (id integer not null primary key autoincrement,
name text not null);
insert into users (name) values ("djfong");
create table nothing(id integer not null primary key autoincrement, name text not null);
create table file(line text);
.quit
EOF

# create the _tables_ table
sqlite3 "$DBFILE" <<EOF
create table _tables_ (id integer not null primary key autoincrement, name text unique not null, schema text not null);
insert into _tables_ (name,schema) values ("bundles",
'{"fields":[{"name":"id",properties:["is_primary_key"]},{"name":"name"},{"name":"uri"}]}');
insert into _tables_ (name,schema) values ("users",
'{"fields":[{"name":"id",properties:["is_primary_key"]},{"name":"name"}]}');
insert into _tables_ (name,schema) values ("nothing",
'{"fields":[{"name":"id",properties:["is_primary_key"]},{"name":"name"}]}');
insert into _tables_ (name,schema) values ("file",
'{"files":[{"name":"line"}]}');
.quit
EOF

# dump the bundles table.
sqlite3 "$DBFILE" << EOF
select * from bundles;
.quit
EOF
