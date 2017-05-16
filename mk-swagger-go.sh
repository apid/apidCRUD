#! /bin/bash
#	mk-swagger-go.sh swagger.yaml > swagger.go
# create swagger.go from swagger.yaml

# ----- start of mainline code
INFILE=${1:-swagger.yaml}
echo '// THIS IS A GENERATED FILE - DO NOT EDIT'
echo 'package apidCRUD'
echo ''
echo -n 'var swaggerJSON = `'
go run cmd/yaml2json/yaml2json.go "$INFILE"
echo '`'
