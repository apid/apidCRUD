package apidCRUD

import (
	"testing"
	"fmt"
	"strings"
	"strconv"
	"os"
	"reflect"
	"runtime"
	"net/http"
	"database/sql"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
)

// TestMain() is called by the test framework before running the tests.
// we use it to initialize the log variable.
func TestMain(m *testing.M) {
	// do this in case functions under test need to log something
	apid.Initialize(factory.DefaultServicesFactory())
	log = apid.Log()

	// required boilerplate
	os.Exit(m.Run())
}

// ---- generic support for testing validator functions

type validatorFunc func(string) (string, error)

type validatorTC struct {
	arg string
	xres string
	xsucc bool
}

func run_validator(t *testing.T, vf validatorFunc, tab []validatorTC) {
	fname := getFunctionName(vf)
	for i, test := range tab {
		call_validator(t, fname, vf, i, test)
	}
}

func call_validator(t *testing.T,
		fname string,
		vf validatorFunc,
		i int,
		test validatorTC) {
	res, err := vf(test.arg)
	msg := "true"
	if err != nil {
		msg = err.Error()
	}
	if !((test.xsucc && err == nil && test.xres == res) ||
	   (!test.xsucc && err != nil)) {
		t.Errorf(`#%d: %s("%s")=("%s","%s"); expected ("%s",%t)`,
			i, fname, test.arg, res, msg,
			test.xres, test.xsucc)
	}
}

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// ----- unit tests for validate_table_name

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

func call_fetchParams(t *testing.T, i int, qp string, xsucc bool) {
	_, err := fetchParamsHelper(qp)
	if xsucc != (err == nil) {
		msg := "true"
		if err != nil {
			msg = err.Error()
		}
		t.Errorf(`#%d: fetchParams("%s")=(%s); expected (%t)`,
			i, qp, msg, xsucc)
	}
}

func Test_fetchParams(t *testing.T) {
	for i, test := range fetchParams_Tab {
		call_fetchParams(t, i, test.arg, test.xsucc)
	}
}

// ----- unit tests for mkVmap()

func strToRawBytes(data string) interface{} {
	rb := sql.RawBytes([]byte(data))
	return &rb
}

func interfaceToStr(data interface{}) (string, error) {
	sp, ok := data.(*string)
	if !ok {
		return "", fmt.Errorf("string conversion error")
	}
	return *sp, nil
}

func rawBytesHelper(strlist []string) []interface{} {
	ret := make([]interface{}, len(strlist))
	for i, s := range strlist {
		ret[i] = interface{}(strToRawBytes(s))
	}
	return ret
}

func mkVmap_Checker(t *testing.T,
		i int,
		keys []string,
		values []string) {
	fn := "mkVmap"
	N := len(keys)
	res, err := mkVmap(keys, rawBytesHelper(values))
	if err != nil {
		t.Errorf("#%d: %s(...) failed", i, fn)
		return
	}
	if N != len(*res) {
		t.Errorf("#%d: %s(...) map length mismatch nkeys", i, fn)
		return
	}
	for j, k := range keys {
		v, err := interfaceToStr((*res)[k])
		if err != nil {
			t.Errorf("#%d: %s(...) rawBytesToStr: %s", j, fn, err)
			return
		}
		if values[j] != v {
			t.Errorf("#%d: %s(...) map value mismatch", j, fn)
			return
		}
	}
}

func Test_mkVmap(t *testing.T) {
	N := 4

	// create the keys and values arrays, with canned values.
	keys := make([]string, N)
	values := make([]string, N)
	for i := 0; i < N; i++ {
		keys[i] = fmt.Sprintf("K%d", i)
		values[i] = fmt.Sprintf("V%d", i)
	}

	// test against initial slices of keys and values arrays.
	for i := 0; i < N+1; i++ {
		mkVmap_Checker(t, i, keys[0:i], values[0:i])
	}
}

// ----- unit tests for mkSQLRow()

func mkSQLRow_Checker(t *testing.T, i int, N int) {
	fn := "mkSQLRow"
	res := mkSQLRow(N)
	if len(res) != N {
		t.Errorf("#%d: %s(%d) failed", i, fn, N)
		return
	}
	for _, v := range res {
		_, ok := v.(*sql.RawBytes)
		if !ok {
			t.Errorf("#%d: %s(%d) sql conversion error", i, fn, N)
			return
		}
	}
}

func Test_mkSQLRow(t *testing.T) {
	M := 5
	for i := 0; i < M; i++ {
		mkSQLRow_Checker(t, i, i)
	}
}

// ----- unit tests for notImplemented()

func Test_notImplemented(t *testing.T) {
	fn := "notImplemented"
	xcode := http.StatusNotImplemented
	code, err := notImplemented()
	if code != xcode {
		t.Errorf("%s returned code %d; expected %d", fn, code, xcode)
	}
	if err == nil {
		t.Errorf("%s returned nil error; expected non-nil", fn)
	}
}

// ----- unit tests for validateSQLValues()

func genList(form string, N int) []string {
	ret := make([]string, N)
	for i := 0; i < N; i++ {
		ret[i] = fmt.Sprintf(form, i)
	}
	return ret
}

func sqlValues_Checker(t *testing.T, form string, N int) {
	fn := "validateSQLValues"
	values := genList(form, N)
	err := validateSQLValues(values)
	if err != nil {
		t.Errorf("%s(...) failed on length=%d", fn, N)
	}
}

func Test_validateSQLValues(t *testing.T) {
	M := 5
	for j := 0; j < M; j++ {
		sqlValues_Checker(t, "V%d", j)
	}

	// empty values OK
	sqlValues_Checker(t, "", 3)
}

// ----- unit tests for validateSQLKeys()

func sqlKeys_Checker(t *testing.T, form string, N int, xsucc bool) {
	fn := "validateSQLKeys"
	values := genList(form, N)
	err := validateSQLKeys(values)
	if xsucc != (err == nil) {
		msg := "true"
		if err != nil {
			msg = err.Error()
		}
		t.Errorf(`%s("%s"...)=%s; expected %t`,
			fn, form, msg, xsucc)
	}
}

func Test_validateSQLKeys(t *testing.T) {
	M := 3
	for j := 0; j < M; j++ {
		sqlKeys_Checker(t, "K%d", j, true)
	}

	// numeric key not OK
	sqlKeys_Checker(t, "%d", 1, false)

	// empty key not OK
	sqlKeys_Checker(t, "", 1, false)
}

// ----- unit tests for nstring()

func nstring_Checker(t *testing.T, s string, n int) {
	fn := "nstring"
	res := nstring(s, n)
	rlist := strings.Split(res, ",")
	if n == 0 {
		// this must be handled as a special case
		// because strings.Split() returns a list of length 1
		// on empty string.
		if res != "" {
			t.Errorf(`%s("%s",%d)="%s"; expected ""`,
				fn, s, n, res)
		}
		return
	} else if n != len(rlist) {
		t.Errorf(`%s("%s",%d)="%s" failed split test`,
			fn, s, n, res)
		return
	}
	for _, v := range rlist {
		if v != s {
			t.Errorf(`%s("%s",%d) bad item "%s"`,
				fn, s, n, v)
		}
	}
}

func Test_nstring(t *testing.T) {
	M := 3
	for j := 0; j < M; j++ {
		nstring_Checker(t, "", j)
		nstring_Checker(t, "abc", j)
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

// ----- unit tests for strListToInterfaces()

func strListToInterfaces_Checker(t *testing.T, form string, M int) {
	fn := "strListToInterfaces"
	strlist := genList(form, M)
	res := strListToInterfaces(strlist)
	n := len(res)
	if M != n {
		t.Errorf(`%s returned length is %d; expected %d`, fn, n, M)
	}
	for i, si := range res {
		str, ok := si.(string)
		if !ok {
			t.Errorf("%s length %d: result item is not a string",
				fn, M)
		}
		if str != strlist[i] {
			t.Errorf(`%s length %d: item="%s"; expected "%s"`,
				fn, M, str, strlist[i])
		}
	}
}

func Test_strListToInterfaces(t *testing.T) {
	M := 3
	for j := 0; j < M; j++ {
		strListToInterfaces_Checker(t, "S%d", j)
	}
}

// ----- unit tests for errorRet()

type errorRet_TC struct {
	code int
	msg string
}

var errorRet_Tab = []errorRet_TC {
	{ 1, "abc" },
	{ 2, "" },
	{ 3, "xyz" },
}

func errorRet_Checker(t *testing.T, i int, code int, msg string) {
	fn := "errorRet"
	err := fmt.Errorf("%s", msg)
	rescode, resdata := errorRet(code, err)
	if code != rescode {
		t.Errorf(`#%d: %s returned (%d,); expected %d`,
			i, fn, rescode, code)
		return
	}
	eresp, ok := resdata.(ErrorResponse)
	if !ok {
		t.Errorf(`#%d: %s ErrorResponse conversion error`, i, fn)
		return
	}
	if code != eresp.Code {
		t.Errorf(`#%d: %s ErrorResponse.Code=%d; expected %d`,
			i, fn, eresp.Code, code)
		return
	}
	if msg != eresp.Message {
		t.Errorf(`#%d: %s ErrorResponse.Message="%s"; expected "%s"`,
			i, fn, eresp.Message, msg)
	}
}

func Test_errorRet(t *testing.T) {
	for i, test := range errorRet_Tab {
		errorRet_Checker(t, i, test.code, test.msg)
	}
}

// ----- unit tests for okRet()

func Test_okRet(t *testing.T) {
	fn := "okRet"
	xcode := http.StatusOK
	data := DeleteResponse{0}
	code, idata := okRet(DeleteResponse{0})
	if xcode != code {
		t.Errorf("%s returned code=%d; expected %d", fn, code, xcode)
		return
	}
	cdata, ok := idata.(DeleteResponse)
	if !ok {
		t.Errorf("%s returned data could not convert", fn)
		return
	}
	if data != cdata {
		t.Errorf("%s returned data does not match", fn)
	}
}
