# apidCRUD

apidCRUD is a plugin for 
[apid](http://github.com/30x/apid).
it handles CRUD (Create/Read/Update/Delete) APIs,
with a simple local database on the back end.

Status:

this is still a WIP,
with some features unimplemented and still subject to change.
at the moment, the /db/_table APIs are implemented,
the /db/_schema APIs are not.

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

if go get errors occur during the glide install phase, try doing `make update`.

## Running apidCRUD
 
for now, this runs apidCRUD in background, listening on localhost:9000.
```
make run
```

several configuration parameters are supported.
see [apid_config.yaml](./apid_config.yaml) .

## Running unit tests
 
this command runs all the unit tests, and also prints overall percent coverage.
```
make unit-test
```

when changes are pushed to github.com,
`make install` and `make unit-test` are run automatically by
[travis-CI](https://travis-ci.org/getting_started).
see [.travis.yml](.travis.yml)

## Running functional tests

this command starts the server and runs tester.sh, which does a series of canned API calls exercising most of the currently-implemented APIs.
```
make func-test
```

there are also a handful of \*test.sh scripts,
each of which calls one of the APIs.

## Adding new APIs

to add a new API:

* create definition in swagger.yaml.
the definition should include an operationId.
be sure to up the version number.
* upload the new swagger.yaml to apistudio.com .
* update the apistudio link in the [Resources](#resources) section of README.md .
* remove the generated file gen_swag.go
to ensure that apiTable will be regenerated.
* put the new handler function in handlers.go,
named after the operationId with "Handler" appended.
* add at least one unit test in handlers_test.go .
* add at least one functional test script in functests/ .
modify tester.sh to call the new functional test.

## Adding unit tests

this project follows the somewhat standard go testing pattern.

* unit test functions are placed in files whose names end with _test.go .
* generally (but not always), functions defined in *FILE*.go have their test functions located in *FILE*_test.go .
* *FILE*_test.go imports the module "testing".
* the test suite for a function named *XYZ*() is a function named Test_*XYZ*() .
* Test_*XYZ*() takes one argument *t* of type \*testing.T .
* Test_*XYZ*() in turn may call other functions from *FILE*_test.go, to run the individual test cases.
* to report a test failure, call `t.Errorf(fmt, ...)` .
* generally, but not always, there may be a variable *XYZ*\_Tab which is an array of test cases, each of type *XYZ*\_TC .
* the type *XYZ*\_TC has fields defining the test case: the function arguments, and the expected results of calling the function with those arguments.
* Test_*XYZ*() just uses a for/range statement to loop over the test cases.
* on each test case, the function *XYZ*_Checker() checks one test case.

Template:
```
// ----- unit tests for XYZ().

// inputs and outputs for one XYZ testcase.
type XYZ_TC struct {
	// CUSTOMIZE
	arg string
	result string
}

// table of XYZ testcases.
var XYZ_Tab = []XYZ_TC {
	// CUSTOMIZE
	{ "arg1", "result1" },
}

// run one testcase for function XYZ.
func XYZ_Checker(cx *testContext, tc *XYZ_TC) {
	// CUSTOMIZE
	result := XYZ(tc.arg)
	cx.assertEqual(tc.result, result)
}

// the XYZ test suite.  run all XYZ testcases.
func Test_XYZ(t *testing.T) {
	cx := newTestContext(t, "XYZ_Tab")
	for _, tc := range XYZ_Tab {
		XYZ_Checker(cx, &tc)
		cx.bump()	// increment testno.
	}
}
```

## Code entry points

note that the code in this module can be invoked in at least 3 ways:

* as an apid plugin, called via apid.InitializePlugins() -
see init.c, plugin.c .
when called this way, initPlugin() uses configuration info from
apid's config file.
* standalone test program - see cmd/apidCRUD/main.go, init.c, plugin.c .
when called this way, initPlugin() uses configuration info from
the apid_config.yaml file in this directory.
* thru unit test framework - see setup_test.go .
when called this way, initPlugin() uses configuration info from
utConfData in globals_test.go .

## Resources

   * [travis-ci for apidCRUD](https://travis-ci.org/30x/apidCRUD)
   * [coveralls for apidCRUD](https://coveralls.io/github/30x/apidCRUD)
   * [godoc for apidCRUD](https://godoc.org/github.com/30x/apidCRUD)
   * [swagger.yaml for apidCRUD](./swagger.yaml)
   * [apistudio for apidCRUD](http://playground.apistudio.io/8548bd01-cb5e-47c7-b2f4-5452c9ca4e66/#/)
