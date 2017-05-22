package apidCRUD

// this module contains definitions and functions of type testContext.

import (
	"testing"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"path"
)

// context for printing assertion failure messages including the tab name
// and item number for table-based testcases.
type testContext struct {
	t *testing.T
	tabName string
	testno int
}

// this Errorf() is like the testing package's Errorf(), but
// prefixes the message with info identifying the test.
func (cx *testContext) Errorf(form string, args ...interface{}) {
	var prefix string
	loc := getTestCaller(2)
	testName := cx.t.Name()
	if cx.tabName == "" {
		prefix = fmt.Sprintf("%s %s: ",
			loc, testName)
	} else {
		prefix = fmt.Sprintf("%s %s #%d %s: ",
			loc, testName, cx.testno, cx.tabName)
	}
	fmt.Printf(prefix + form + "\n", args...)
	cx.t.Fail()
}

// advance the test counter
func (cx *testContext) bump() {
	cx.testno++
}

// ----- definition of assertions

// assertEqual() fails the test unless exp and act are equal.
func (cx *testContext) assertEqual(exp interface{}, act interface{}, msg string) bool {
	if exp != act {
		cx.Errorf(`*** Assertion Failed: %s
   expected: (%T)<%v>
   got:      (%T)<%v>`,
			msg, exp, exp, act, act)
		return false
	}
	return true
}

// assertEqualObj() fails the test unless exp and act are deeply equal.
func (cx *testContext) assertEqualObj(exp interface{}, act interface{}, msg string) bool {
	if !reflect.DeepEqual(exp, act) {
		cx.Errorf(`assertion failed: %s
   expected: <%s>
   got:      <%s>`,
			msg, exp, act)
		return false
	}
	return true
}

// assertTrue() fails the test unless act is true.
func (cx *testContext) assertTrue(act bool, msg string) bool {
	if !act {
		cx.Errorf(`assertion failed: %s, is %t; s/b true`,
			msg, act)
		return false
	}
	return true
}

// assertErrorNil fails the test unless act is nil.
func (cx *testContext) assertErrorNil(err error, msg string) bool {
	if err != nil {
		cx.Errorf(`assertion failed: %s, gave error [%s]`,
			msg, err)
		return false
	}
	return true
}

// newTestContext() creates a new test context.
// tabName is intended to be a single optional argument.
// if it is specified and nonempty, it signifies that any
// assertion failures will be labelled with the tabName
// and an offset.
func newTestContext(t *testing.T,
		opt ...string) *testContext {
	tabName := ""
	if len(opt) > 0 {
		tabName = opt[0]
	}
	return &testContext{t, tabName, 0}
}

// return the name of the given function
func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// getTestCaller() returns the name of the n-th function up the call
// chain from the caller of "this" function.  ie, the function that
// called getTestCaller would be 0; the caller of that function
// would be 1.
// the returned name includes the file and line number of the
// function call.
func getTestCaller(nth int) string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(nth+2, pc)
	if n == 0 {
		return "?"
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fn := frame.Function
	// skip over the module name if any.
	i := strings.LastIndexByte(fn, '.')
	if i >= 0 {
		fn = fn[i+1:]
	}
	return fmt.Sprintf("%s@%s:%d",
		fn, path.Base(frame.File), frame.Line)
}
