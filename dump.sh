#! /bin/bash

DATABASE=apidCRUD.db
echo ".tables" | sqlite3 "$DATABASE"
echo "select * from bundles;" | sqlite3 "$DATABASE"
