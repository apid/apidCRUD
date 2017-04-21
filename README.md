# apidCRUD

apidCRUD is a plugin for 
[apid](http://github.com/30x/apid).
it handles CRUD (Create/Read/Update/Delete) APIs,
with a simple local database on the back end.

this is still a WIP,
with some features unimplemented and still subject to change.

## Functional description

see the file [swagger.yaml](swagger.yaml).

## Apid Services Used

* Config Service
* Log Service
* API Service

## Building apidCRUD

to build apidCRUD, run:
```
make install
```

for now, this just builds the standalone plugin application.

## Running apidCRUD
 
```
make run
```

for now, this runs apidCRUD in background, listening on localhost:9000.

## Testing
 
```
make unit-test
make coverage

make run
make func-test
```

when changes are pushed to github.com,
`make install` and `make unit-test` are run automatically by
[travis-CI](https://travis-ci.org/getting_started).
see [.travis.yml](.travis.yml)

