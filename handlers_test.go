package apidCRUD

import "testing"
import "os"
import "github.com/30x/apid-core"
import "github.com/30x/apid-core/factory"

// TestMain() is called by the test framework before running the tests.
// we use it to initialize the log variable.
func TestMain(m *testing.M) {
	// do this in case functions under test need to log something
	apid.Initialize(factory.DefaultServicesFactory())
	log = apid.Log()

	// required boilerplate
	os.Exit(m.Run())
}

// ---- unit tests for validate_xxx()

type validate_offset_TC struct {
	arg string		// input string
	xres string		// expected result
	xsucc bool		// expected success
}

var validate_offset_Tab = []validate_offset_TC {
	{ "", "0", true },
	{ "0", "0", true },
	{ "12345678", "12345678", true },
	{ "-12345678", "-12345678", true },
	{ "+12345678", "12345678", true },
	{ "12345678.", "", false },
	{ " 12345678", "", false },
	{ "12345678 ", "", false },
}

func Test_validate_offset(t *testing.T) {
	fn := "validate_offset"
	for i, test := range validate_offset_Tab {
		ret, err := validate_offset(test.arg)
		msg := "true"
		if err != nil {
			msg = err.Error()
		}
		if !((test.xsucc && err == nil && test.xres == ret) ||
		   (!test.xsucc && err != nil)) {
			t.Errorf(`#%d: %s("%s")=("%s","%s"); expected ("%s",%t)`,
				i, fn, test.arg, ret, msg,
				test.xres, test.xsucc)
		}
	}
}


// ---- unit tests for notIdentChar()

type notIdentChar_TC struct {
	c rune
	res bool
}

var notIdentChar_Tab = []notIdentChar_TC {
	{'&', true},
	{'a', false},
	{'z', false},
	{'A', false},
	{'Z', false},
	{'0', false},
	{'9', false},
	{'_', false},
	{'|', true},
	{'\000', true},
	{'.', true},
	{',', true},
	{'/', true},
}

func Test_notIdentChar(t *testing.T) {
	fn := "isValidIdent"
	for i, test := range notIdentChar_Tab {
		res := notIdentChar(test.c)
		if res != test.res {
			t.Errorf(`#%d: %s('%c')=%t; expected %t`, i, fn, test.c, res, test.res)
		}
	}
}

// ----- unit tests for isValidIdent()

type isValidIdent_TC struct {
	s string
	res bool
}

var isValidIdent_Tab = []isValidIdent_TC {
	{"_ABCXYZabcxyz0123456789", true},
	{"_ABCabc0123.", false},
	{"abc.def", false},
	{"abc:def", false},
	{"abc/def", false},
	{"abc!def", false},
	{"abc?def", false},
	{"abc$def", false},
	{"", false},
}

func Test_isValidIdent(t *testing.T) {
	fn := "isValidIdent"
	for i, test := range isValidIdent_Tab {
		res := isValidIdent(test.s)
		if res != test.res {
			t.Errorf(`#%d: %s("%s")=%t; expected %t`, i, fn, test.s, res, test.res)
		}
	}
}
