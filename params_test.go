package apidCRUD

import (
	"testing"
	"fmt"
	"strings"
	"strconv"
	"net/http"
)

// parsedUrl represents the information from a "url description".
type parsedUrl struct {
	path string
	pathParams map[string]string
	queryStr string
	body string
}

// parseUrlDesc() breaks down a "url description" string into its
// component parts.  the parts are separated by "|" chars.
//	PATH | PATH_PARAMS | QUERY_PARMS | BODY
// all parts are optional.
// PATH_PARAMS and QUERY_PARAMS are of the form VAR=VALUE&VAR=VALUE&...
func parseUrlDesc(urlStr string) *parsedUrl {
	pathParams := make(map[string]string)
	w := strings.SplitN(urlStr, "|", 4)

	switch len(w) {
	case 1:
		return &parsedUrl{w[0], pathParams, "", ""}
	case 2:
		return &parsedUrl{w[0], strToMap(w[1]), "", ""}
	case 3:
		return &parsedUrl{w[0], strToMap(w[1]), w[2], ""}
	case 4:
		return &parsedUrl{w[0], strToMap(w[1]), w[2], w[3]}
	default:
		return &parsedUrl{}
	}
}

// ---- generic support for testing validator functions

// the type of a validator function.
type validatorFunc func(string) (string, error)

// type validator_TC is the structure of one test case for a validator.
type validator_TC struct {
	arg string
	xres string
	xsucc bool
}

// run thru the table of test cases for the given validator function.
func run_validator(cx *testContext, vf validatorFunc, tab []validator_TC) {
	for _, tc := range tab {
		validator_Checker(cx, vf, tc)
		cx.bump()
	}
}

// run one test case thru the given validator function.
func validator_Checker(cx *testContext,
		vf validatorFunc,
		tc validator_TC) {
	result, err := vf(tc.arg)
	ok := (err == nil && tc.xres == result)
	cx.assertEqual(tc.xsucc, ok, "result")
}

// ----- unit tests for validate_id_field

var validate_id_field_Tab = []validator_TC {
	{ "", "id", true },
	{ "x", "x", true },
	{ "X", "X", true },
	{ "_", "_", true },
	{ "1", "1", false },
}

func Test_validate_id_field(t *testing.T) {
	cx := newTestContext(t, "validate_id_field_Tab")
	run_validator(cx, validate_id_field, validate_id_field_Tab)
}

// ----- unit tests for validate_fields

var validate_fields_Tab = []validator_TC {
	{ "", "*", true },
	{ "f1", "f1", true },
	{ "f1,f2", "f1,f2", true },
	{ "f1,", "f1,", false },
	{ ",f1,", ",f1", false },
	{ " f1,", " f1", false },
	{ "f1 ", "f1 ", false },
}

func Test_validate_fields(t *testing.T) {
	cx := newTestContext(t, "validate_fields_Tab")
	run_validator(cx, validate_fields, validate_fields_Tab)
}

// ----- unit tests for validate_table_name

var validate_table_name_Tab = []validator_TC {
	{ "", "", false },
	{ "a", "a", true },
	{ "1", "1", false },
	{ "a-2", "a-2", false },
	{ ".", ".", false },
	{ "xyz", "xyz", true },
}

func Test_validate_table_name(t *testing.T) {
	cx := newTestContext(t, "validate_table_name_Tab")
	run_validator(cx, validate_table_name, validate_table_name_Tab)
}

// ----- unit tests for validate_id

var validate_id_Tab = []validator_TC {
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
	cx := newTestContext(t, "validate_id_Tab")
	run_validator(cx, validate_id, validate_id_Tab)
}

// ----- unit tests for validate_limit

func Test_validate_limit(t *testing.T) {
	var strMaxRecs = strconv.Itoa(maxRecs)	// maxRecs converted to a string

	var validate_limit_Tab = []validator_TC {
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

	cx := newTestContext(t, "validate_limit_Tab")
	run_validator(cx, validate_limit, validate_limit_Tab)
}

// ----- unit tests for validate_ids()

var validate_ids_Tab = []validator_TC {
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
	cx := newTestContext(t, "validate_ids_Tab")
	run_validator(cx, validate_ids, validate_ids_Tab)
}

// ----- unit tests for validate_offset()

var validate_offset_Tab = []validator_TC {
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
	cx := newTestContext(t, "validate_offset_Tab")
	run_validator(cx, validate_offset, validate_offset_Tab)
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

func notIdentChar_Checker(cx *testContext, tc notIdentChar_TC) {
	result := notIdentChar(tc.c)
	cx.assertEqual(tc.res, result,
		fmt.Sprintf("<%q>", tc.c))
}

func Test_notIdentChar(t *testing.T) {
	cx := newTestContext(t, "notIdentChar_Tab")
	for _, tc := range notIdentChar_Tab {
		notIdentChar_Checker(cx, tc)
		cx.bump()
	}
}

// ----- test table for a field with no validator

var validate_nofield_Tab = []validator_TC {
	{ "", "", false },
}

// ----- unit tests for isValidIdent()

type isValidIdent_TC struct {
	arg string
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

func isValidIdent_Checker(cx *testContext, tc isValidIdent_TC) {
	res := isValidIdent(tc.arg)
	cx.assertEqual(tc.res, res,
		fmt.Sprintf("<%s>", tc.arg))
}

func Test_isValidIdent(t *testing.T) {
	cx := newTestContext(t, "isValidIdent_Tab")
	for _, tc := range isValidIdent_Tab {
		isValidIdent_Checker(cx, tc)
		cx.bump()
	}
}

// ----- unit tests for getParam()

func parseHandlerArg(verb string, descStr string) *apiHandlerArg {
	desc := parseUrlDesc(descStr)
	path := desc.path + "?" + desc.queryStr
	req, _ := http.NewRequest(verb, path, strings.NewReader(desc.body))
	return mkApiHandlerArg(req, desc.pathParams)
}

// checker for params that come from the path
func getPathParam_Checker(cx *testContext,
		paramName string, val string) (string, error) {
	descStr := fmt.Sprintf("/apid/db|%s=%s", paramName, val)
	harg := parseHandlerArg(http.MethodGet, descStr)
	return harg.getParam(paramName)
}

// checker for params that come from the query portion of url
func getQueryParam_Checker(cx *testContext,
		paramName string, val string) (string, error) {
	descStr := fmt.Sprintf("/apid/db||%s=%s", paramName, val)
	harg := parseHandlerArg(http.MethodGet, descStr)
	return harg.getParam(paramName)
}

func Test_getParam_idInPath(t *testing.T) {

	// test getParam on id values (as path param)
	cx := newTestContext(t, "validate_id_Tab")
	run_validator(cx,
		func(val string) (string, error) {
			return getPathParam_Checker(cx, "id", val)
		},
		validate_id_Tab)
}

func Test_getParam_idInQuery(t *testing.T) {

	// test getParam on id values (as query param)
	cx := newTestContext(t, "validate_id_Tab")
	run_validator(cx,
		func(val string) (string, error) {
			return getQueryParam_Checker(cx, "id", val)
		},
		validate_id_Tab)
}

func Test_getParam_ids(t *testing.T) {
	// test getParam on ids values (as query param)
	cx := newTestContext(t, "validate_ids_Tab")
	run_validator(cx,
		func(val string) (string, error) {
			return getQueryParam_Checker(cx, "ids", val)
		},
		validate_ids_Tab)
}

func Test_getParam_id_field(t *testing.T) {
	// test getParam on id_field values (as query param)
	cx := newTestContext(t, "validate_id_field_Tab")
	run_validator(cx,
		func(val string) (string, error) {
			return getQueryParam_Checker(cx, "id_field", val)
		},
		validate_id_field_Tab)
}

func Test_getParam_nofield(t *testing.T) {

	// test getParam on a field with no validator (as query param)
	cx := newTestContext(t, "validate_nofield_Tab")
	run_validator(cx,
		func(val string) (string, error) {
			return getQueryParam_Checker(cx, "nofield", val)
		},
		validate_nofield_Tab)
}

// ----- unit tests for fetchParams()

type fetchParams_TC struct {
	verb string
	desc string
	nameStr string		// list of names to fetch
	xsucc bool		// expected success
}

var fetchParams_Tab = []fetchParams_TC {
	{ http.MethodGet, "/db/abc|table_name=faketab&id=123", "id,table_name", true },
	{ http.MethodGet, "/db/abc||id=123&ids=123,456", "id,ids", true },
	{ http.MethodGet, "/db/abc||id=1&fields=a,b,c", "id,fields", true },
	{ http.MethodGet, "/db/abc||junk=1&fields=a,b,c", "junk,fields", false },
}

// strToMap() constructs a map object from a string
// in which mappings K=V are separated by & chars.
func strToMap(vars string) map[string]string {
	vlist := mySplit(vars, "&")
	ret := map[string]string{}
	for _, parm := range vlist {
		words := strings.SplitN(parm, "=", 2)
		switch len(words) {
		case 1:
			ret[words[0]] = ""
		case 2:
			ret[words[0]] = words[1]
		}
	}
	return ret
}

func fetchParamsHelper(verb string,
		descStr string,
		nameStr string) error {

	harg := parseHandlerArg(verb, descStr)

	namesList := mySplit(nameStr, ",")
	vmap, err := fetchParams(harg, namesList...)
	if err != nil {
		return err
	}

	// check that the map has the expected number of keys
	nvmap := len(vmap)
	nnames := len(namesList)
	if nvmap != nnames {
		return fmt.Errorf("map has %d entries; expected %d",
				nvmap, nnames)
	}

	// check that each expected name is there
	for _, name := range namesList {
		_, ok := vmap[name]
		if !ok {
			return fmt.Errorf("map does not have %s", name)
		}
	}

	return nil
}

// handle one testcase
func fetchParams_Checker(cx *testContext, tc fetchParams_TC) {
	err := fetchParamsHelper(tc.verb,
			tc.desc,
			tc.nameStr)
	cx.assertEqual(tc.xsucc, err == nil, tc.desc)
}

func Test_fetchParams(t *testing.T) {
	cx := newTestContext(t, "fetchParams_Tab")
	for _, tc := range fetchParams_Tab {
		fetchParams_Checker(cx, tc)
		cx.bump()
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
func aToIdType_Checker(cx *testContext, tc aToIdType_TC) {
	res := idType(aToIdType(tc.arg))
	cx.assertEqualObj(tc.xval, res, "result")
}

func Test_aToIdType(t *testing.T) {
	cx := newTestContext(t, "aToIdType_Tab")
	for _, tc := range aToIdType_Tab {
		aToIdType_Checker(cx, tc)
		cx.bump()
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
func idTypeToA_Checker(cx *testContext, tc idTypeToA_TC) {
	result := idTypeToA(int64(tc.arg))
	cx.assertEqual(tc.xval, result, "result")
}

func Test_idTypeToA(t *testing.T) {
	cx := newTestContext(t, "idTypeToA_Tab")
	for _, tc := range idTypeToA_Tab {
		idTypeToA_Checker(cx, tc)
		cx.bump()
	}
}
