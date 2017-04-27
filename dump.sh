#! /bin/bash

DATABASE=apidCRUD.db
echo "select * from bundles;" | sqlite3 "$DATABASE"
