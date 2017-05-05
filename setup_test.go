package apidCRUD

// this module does global setup for unit tests.
// it also contains some generic test support code.

import (
	"testing"
	"fmt"
	"os"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
)

// context for printing test failure messages.
type testContext struct {
	t *testing.T
	suiteName string
	funcName string
	testno int
}

// this Errorf() is like the testing package's Errorf(), but
// prefixes the message with info identifying the test.
func (cx *testContext) Errorf(form string, args ...interface{}) {
	var prefix string
	if cx.funcName == "" {
		prefix = fmt.Sprintf("%s #%d: ",
			cx.suiteName, cx.testno)
	} else {
		prefix = fmt.Sprintf("%s #%d: %s ",
			cx.suiteName, cx.testno, cx.funcName)
	}
	cx.t.Errorf(prefix + form, args...)
}

// advance the test counter
func (cx *testContext) bump() {
	cx.testno++
}

func (cx *testContext) setFuncName(name string) {
	cx.funcName = name
}

// newTestContext() creates a new test context.
// the optional final argument is the funcName for the context.
func newTestContext(t *testing.T,
		suiteName string,
		opt ...string) *testContext {
	funcName := suiteName
	if len(opt) > 0 {
		funcName = opt[0]
	}
	return &testContext{t, suiteName, funcName, 0}
}

var testServices = factory.DefaultServicesFactory()

// TestMain() is called by the test framework before running the tests.
// we use it to initialize the log variable.
func TestMain(m *testing.M) {
	// do this in case functions under test need to log something
	apid.Initialize(testServices)
	log = apid.Log()
	log.Debugf("in TestMain")

	// for testing purposes, set global maxRecs to some smallish value
	maxRecs = 7

	var err error
	db, err = fakeInitDB()
	if err != nil {
		panic(err.Error())
	}

	// required boilerplate
	os.Exit(m.Run())
}
