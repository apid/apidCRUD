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
// template for unit tests of a function XYZ.

// inputs and outputs for one XYZ testcase.
type XYZ_TC struct {
	// CUSTOMIZE
	arg string
	result string
}

// table of XYZ testcases.
var XYZ_Tab = XYZ_TC {
	// CUSTOMIZE
	{ "arg1", "result1" },
}

// run one testcase for function XYZ.
func XYZ_Checker(t *testing.T, testno int, tc XYZ_TC) {
	// CUSTOMIZE
	result := XYZ(tc.args)
	if tc.result !- result {
		t.Errorf(`#%d: XYZ("%s")="%s"; expected "%s"`,
			testno, tc.args, result, tc.result)
	}
}

// the XYZ test suite.  run all XYZ testcases.
func Test_XYZ(t *testing.T) {
	for testno, tc := range XYZ_Tab {
		XYZ_Checker(t, testno, tc)
	}
}
```

## Code entry points

note that the code can be invoked in at least 3 ways:

   * as an apid plugin, called via apid.InitializePlugins() - see init.c, plugin.c
   * standalone test program - see cmd/apidCRUD/main.go, init.c, plugin.c
   * thru unit test framework - see setup_test.go

## Badges

[![Coverage Status](https://coveralls.io/repos/github/30x/apidCRUD/badge.svg?branch=master)](https://coveralls.io/github/30x/apidCRUD?branch=master)
