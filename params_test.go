package apidCRUD

import (
	"testing"
	"fmt"
	"strings"
	"strconv"
	"reflect"
	"runtime"
	"net/http"
)

// ---- generic support for testing validator functions

// the type of a validator function.
type validatorFunc func(string) (string, error)

// type validatorTC is the structure of one test case for a validator.
type validatorTC struct {
	arg string
	xres string
	xsucc bool
}

// run thru the table of test cases for the given validator function.
func run_validator(t *testing.T, vf validatorFunc, tab []validatorTC) {
	fname := getFunctionName(vf)
	for i, test := range tab {
		validator_Checker(t, fname, vf, i, test)
	}
}

// run one test case thru the given validator function.
func validator_Checker(t *testing.T,
		fname string,
		vf validatorFunc,
		i int,
		test validatorTC) {
	res, err := vf(test.arg)
	msg := errRep(err)
	if !((test.xsucc && err == nil && test.xres == res) ||
	   (!test.xsucc && err != nil)) {
		t.Errorf(`#%d: %s("%s")=("%s","%s"); expected ("%s",%t)`,
			i, fname, test.arg, res, msg,
			test.xres, test.xsucc)
	}
}

// return the name of the given function
func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// ----- unit tests for validate_id_field

var validate_id_field_Tab = []validatorTC {
	{ "", "id", true },
	{ "x", "x", true },
	{ "X", "X", true },
	{ "_", "_", true },
	{ "1", "1", false },
}

func Test_validate_id_field(t *testing.T) {
	run_validator(t, validate_id_field, validate_id_field_Tab)
}

// ----- unit tests for validate_fields

var validate_fields_Tab = []validatorTC {
	{ "", "*", true },
	{ "f1", "f1", true },
	{ "f1,f2", "f1,f2", true },
	{ "f1,", "f1,", false },
	{ ",f1,", ",f1", false },
	{ " f1,", " f1", false },
	{ "f1 ", "f1 ", false },
}

func Test_validate_fields(t *testing.T) {
	run_validator(t, validate_fields, validate_fields_Tab)
}

// ----- unit tests for validate_table_name

var validate_table_name_Tab = []validatorTC {
	{ "", "", false },
	{ "a", "a", true },
	{ "1", "1", false },
	{ "a-2", "a-2", false },
	{ ".", ".", false },
	{ "xyz", "xyz", true },
}

func Test_validate_table_name(t *testing.T) {
	run_validator(t, validate_table_name, validate_table_name_Tab)
}

// ----- unit tests for validate_id

var validate_id_Tab = []validatorTC {
	{ "", "", false },			// empty
	{ " ", " ", false },			// blank
	{ "0", "0", true },			// simple
	{ "-1", "-1", true },			// negative
	{ "0x21", "", false },			// go-ism
	{ "00021", "21", true },		// go-ism
	{ "1 ", "", false },			// trailing space
	{ " 1", "", false },			// leading space
	{ "2,1", "", false },			// multiple
	{ "1_000_000", "1_000_000", false },	// go-ism
	{ "1000", "1000", true },		// 1E3
	{ "1000000", "1000000", true },		// 1E6
	{ "1000000000", "1000000000", true },	// 1E9
	{ "1000000000000", "1000000000000", true },  // 1E12
	{ "1000000000000000", "1000000000000000", true },  // 1E15
	{ "1000000000000000000000", "1000000000000000000000", false },	// will overflow
}

func Test_validate_id(t *testing.T) {
	run_validator(t, validate_id, validate_id_Tab)
}

// ----- unit tests for validate_limit

var strMaxRecs = strconv.Itoa(maxRecs)

var validate_limit_Tab = []validatorTC {
	{ "", strMaxRecs, true },
	{ " ", "", false },
	{ " 1", "", false },
	{ "1 ", "", false },
	{ "1", "1", true },
	{ "-1", strMaxRecs, true },
	{ "100000", strMaxRecs, true },
	{ "1000000", strMaxRecs, true },
	{ "1000000000", strMaxRecs, true },
	{ "1000000000000", strMaxRecs, true },
}

func Test_validate_limit(t *testing.T) {
	run_validator(t, validate_limit, validate_limit_Tab)
}

// ----- unit tests for validate_ids()

var validate_ids_Tab = []validatorTC {
	{ "", "", true },			// empty list
	{ " ", " ", false },			// blanks
	{ "0x21", "", false },			// go-ism
	{ "00021", "21", true },		// go-ism
	{ "0", "0", true },
	{ "-1", "-1", true },
	{ "0x21", "", false },
	{ "0,0,1,1,1", "0,0,1,1,1", true },
	{ "1 ", "", false },
	{ " 1", "", false },
	{ "1, -1", "", false },
	{ "2,1,", "", false },
	{ "1_000_000", "1_000_000", false },
	{ "1000", "1000", true },
	{ "1000000", "1000000", true },
	{ "1000000000", "1000000000", true },
	{ "1000000000000", "1000000000000", true },
}

func Test_validate_ids(t *testing.T) {
	run_validator(t, validate_ids, validate_ids_Tab)
}

// ----- unit tests for validate_offset()

var validate_offset_Tab = []validatorTC {
	{ "", "0", true },
	{ "0", "0", true },
	{ "12345678", "12345678", true },
	{ "-12345678", "-12345678", true },
	{ "+12345678", "12345678", true },
	{ "12345678.", "", false },
	{ " 12345678", "", false },
	{ "12345678 ", "", false },
	{ "1000", "1000", true },
	{ "1000000", "1000000", true },
	{ "1000000000", "1000000000", true },
	{ "1000000000000", "1000000000000", true },
}

func Test_validate_offset(t *testing.T) {
	run_validator(t, validate_offset, validate_offset_Tab)
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

// ----- test table for a field with no validator

var validate_nofield_Tab = []validatorTC {
	{ "", "", false },
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

// ----- unit tests for newExtReq()

func mkRequest(path string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, path, nil)
}

// return an ExtReq object for testing, based on the given path.
func mkExtReq(path string) (*extReq, error) {
	req, err := mkRequest(path)
	if err != nil {
		return nil, err
	}
	return newExtReq(req, validators)
}

func Test_newExtReq(t *testing.T) {
	fn := "newExtReq"
	xr, err := mkExtReq("/apid/db")
	if err != nil {
		t.Errorf("%s failure: %s", fn, err)
		return
	}
	if xr == nil {
		t.Errorf("%s returned nil", fn)
	}
}

// ----- unit tests for getParam()

func getParam_Checker(t *testing.T,
		paramName string,
		val string) (string, error) {
	path := fmt.Sprintf("/apid/db?%s=%s", paramName, val)
	xr, err := mkExtReq(path)
	if err != nil {
		return "", nil
	}
	return xr.getParam(paramName)
}

func Test_getParam(t *testing.T) {

	// test getParam on id values
	run_validator(t,
		func(val string) (string, error) {
			return getParam_Checker(t, "id", val)
		},
		validate_id_Tab)

	// test getParam on ids values
	run_validator(t,
		func(val string) (string, error) {
			return getParam_Checker(t, "ids", val)
		},
		validate_ids_Tab)

	// test getParam on id_field values
	run_validator(t,
		func(val string) (string, error) {
			return getParam_Checker(t, "id_field", val)
		},
		validate_id_field_Tab)

	// test getParam on a field with no validator
	run_validator(t,
		func(val string) (string, error) {
			return getParam_Checker(t, "nofield", val)
		},
		validate_nofield_Tab)
}

// ----- unit tests for fetchParams()

type fetchParams_TC struct {
	arg string	// query params to use in call
	xsucc bool	// expected success
}

var fetchParams_Tab = []fetchParams_TC {
	{ "id=123", true },
	{ "id=123&ids=123,456", true },
	{ "id=1&fields=a,b,c", true },
	{ "junk=1&fields=a,b,c", false },
}

func fetchParamsHelper(qp string) (map[string]string, error) {
	qplist := strings.Split(qp, "&")
	names := make([]string, len(qplist))
	for i, parm := range qplist {
		nv := strings.SplitN(parm, "=", 2)
		names[i] = nv[0]
	}

	req, err := mkRequest("/api/db?" + qp)
	if err != nil {
		vmap := map[string]string{}
		return vmap, err
	}

	vmap, err := fetchParams(req, names...)
	if err != nil {
		return vmap, err
	}

	// check that the map has the expected number of keys
	nvmap := len(vmap)
	nnames := len(names)
	if nvmap != nnames {
		err := fmt.Errorf("map has %d entries; expected %d",
				nvmap, nnames)
		return vmap, err
	}

	// check that each expected name is there
	for _, name := range names {
		_, ok := vmap[name]
		if !ok {
			err := fmt.Errorf("map does not have %s", name)
			return vmap, err
		}
	}

	return vmap, nil
}

// handle one testcase
func fetchParams_Checker(t *testing.T, i int, qp string, xsucc bool) {
	_, err := fetchParamsHelper(qp)
	if xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: fetchParams("%s")=(%s); expected (%t)`,
			i, qp, msg, xsucc)
	}
}

func Test_fetchParams(t *testing.T) {
	for i, test := range fetchParams_Tab {
		fetchParams_Checker(t, i, test.arg, test.xsucc)
	}
}

// ----- unit tests for aToIdType()

type aToIdType_TC struct {
	arg string
	xval idType
}

var aToIdType_Tab = []aToIdType_TC {
	{ "", -1 },
	{ "x", -1 },
	{ "0x", -1 },
	{ "-1", -1 },
	{ "-1000000", -1000000 },
	{ "0", 0 },
	{ "1", 1 },
	{ "10", 10 },
	{ "10000", 10000 },
	{ "1000000000000", 1000000000000 },
}

// handle one testcase.
func aToIdType_Checker(t *testing.T, i int, test aToIdType_TC) {
	fn := "aToIdType"
	res := idType(aToIdType(test.arg))
	if test.xval != res {
		t.Errorf(`#%d: %s("%s")=%d; expected %d`,
			i, fn, test.arg, res, test.xval)
	}
}

func Test_aToIdType(t *testing.T) {
	for i, test := range aToIdType_Tab {
		aToIdType_Checker(t, i, test)
	}
}

// ----- unit tests for idTypeToA()

type idTypeToA_TC struct {
	arg idType
	xval string
}

var idTypeToA_Tab = []idTypeToA_TC {
	{ 0, "0" },
	{ 1, "1" },
	{ -1, "-1" },
	{ -10000000000, "-10000000000" },
	{ 10000000000000, "10000000000000" },
}

// handle one testcase
func idTypeToA_Checker(t *testing.T, i int, test idTypeToA_TC) {
	fn := "idTypeToA"
	res := idTypeToA(int64(test.arg))
	if test.xval != res {
		t.Errorf(`#%d: %s(%d)="%s"; expected "%s"`,
			i, fn, test.arg, res, test.xval)
	}
}

func Test_idTypeToA(t *testing.T) {
	for i, test := range idTypeToA_Tab {
		idTypeToA_Checker(t, i, test)
	}
}
