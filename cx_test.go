package apidCRUD

// this module contains definitions and functions of type testContext.

import (
	"testing"
	"fmt"
	"reflect"
	"runtime"
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

// ----- definition of assertions

func (cx *testContext) assertEqualInt(exp int, act int, msg string) bool {
	if exp != act {
		cx.Errorf(`assertion failed: %s, got %d; expected %d`,
			msg, act, exp)
		return false
	}
	return true
}

func (cx *testContext) assertEqualStr(exp string, act string, msg string) bool {
	if exp != act {
		cx.Errorf(`assertion failed: %s, got %s; expected %s`,
			msg, act, exp)
		return false
	}
	return true
}

func (cx *testContext) assertEqualBool(exp bool, act bool, msg string) bool {
	if exp != act {
		cx.Errorf(`assertion failed: %s, got %t; expected %t`,
			msg, act, exp)
		return false
	}
	return true
}

func (cx *testContext) assertEqualObj(exp interface{}, act interface{}, msg string) bool {
	if !reflect.DeepEqual(exp, act) {
		cx.Errorf(`assertion failed: %s, got %s; expected %s`,
			msg, act, exp)
		return false
	}
	return true
}

func (cx *testContext) assertTrue(act bool, msg string) bool {
	if !act {
		cx.Errorf(`assertion failed: %s, is %t s/b true`,
			msg, act)
		return false
	}
	return true
}

func (cx *testContext) assertErrorNil(err error, msg string) bool {
	if err != nil {
		cx.Errorf(`assertion failed: %s, gave error [%s]`,
			msg, err)
		return false
	}
	return true
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

// return the name of the given function
func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
