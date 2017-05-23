package apidCRUD

import (
	"testing"
	"fmt"
	"strconv"
	"strings"
	"sort"
	"net/http"
	"database/sql"
)

// mySplit() is like strings.Split() except that
// it returns a 0-length slice when s is the empty string.
func mySplit(str string, sep string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, sep)
}

// ----- unit tests for mkSQLRow()

var mySQLRow_Tab = []int {
	0, 1, 2, 4,
}

func mkSQLRow_Checker(cx *testContext, N int) {
	res := mkSQLRow(N)
	cx.assertEqual(N, len(res), "number of rows")
	for _, v := range res {
		_, ok := v.(*sql.RawBytes)
		if !cx.assertTrue(ok, "sql conversion error") {
			return
		}
	}
}

func Test_mkSQLRow(t *testing.T) {
	cx := newTestContext(t, "mySQLRow_Tab")
	for _, tc := range mySQLRow_Tab {
		mkSQLRow_Checker(cx, tc)
		cx.bump()
	}
}

// ----- unit tests for validateSQLKeys()

type validateSQLKeys_TC struct {
	size int
	form string
	xres bool
}

var validateSQLKeys_Tab = []validateSQLKeys_TC {
	{0, "K%d", true},	// regular key - ok
	{1, "K%d", true},	// regular key
	{2, "K%d", true},	// regular key
	{3, "K%d", true},	// regular key
	{3, "%d", false},	// purely numeric key - bad
	{1, "", false},		// empty key - bad
}

func genList(form string, N int) []string {
	ret := make([]string, N)
	for i := 0; i < N; i++ {
		ret[i] = fmt.Sprintf(form, i)
	}
	return ret
}

func sqlKeys_Checker(cx *testContext, tc *validateSQLKeys_TC) {
	values := genList(tc.form, tc.size)
	err := validateSQLKeys(values)
	cx.assertEqual(tc.xres, err==nil, "success of call")
}

func Test_validateSQLKeys_positive(t *testing.T) {
	cx := newTestContext(t, "validateSQLKeys_Tab")
	for _, tc := range validateSQLKeys_Tab {
		sqlKeys_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for nstring()

type nstring_TC struct {
	n int
	s string
	xres string
}

var nstring_Tab = []nstring_TC {
	{0, "", ""},
	{1, "", ""},
	{2, "", ","},
	{3, "", ",,"},
	{0, "abc", ""},
	{1, "abc", "abc"},
	{2, "abc", "abc,abc"},
	{3, "abc", "abc,abc,abc"},
}

func nstring_Checker(cx *testContext, tc *nstring_TC) {
	result := nstring(tc.s, tc.n)
	cx.assertEqual(tc.xres, result, "result")
}

func Test_nstring(t *testing.T) {
	cx := newTestContext(t, "nstring_Tab")
	for _, tc := range nstring_Tab {
		nstring_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for errorRet()

type errorRet_TC struct {
	code int
	msg string
	dmsg string
}

var errorRet_Tab = []errorRet_TC {
	{ 1, "abc", "" },
	{ 2, "", "" },
	{ 3, "xyz", "" },
	{ 4, "xyz", "with msg" },
}

func errorRet_Checker(cx *testContext, tc *errorRet_TC) {
	err := fmt.Errorf("%s", tc.msg)
	res := errorRet(tc.code, err, tc.dmsg)
	cx.assertEqual(tc.code, res.code, "returned code")
	eresp, ok := res.data.(ErrorResponse)
	if !cx.assertTrue(ok, "ErrorResponse conversion error") {
		return
	}
	cx.assertEqual(tc.code, eresp.Code, "ErrorResponse.Code")
	cx.assertEqual(tc.msg, eresp.Message, "ErrorResponse.Message")
}

func Test_errorRet(t *testing.T) {
	cx := newTestContext(t, "errorRet_Tab")
	for _, tc := range errorRet_Tab {
		errorRet_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for mkIdClause()

func fakeParams(paramstr string) map[string]string {
	ret := map[string]string{}
	if paramstr == "" {
		return ret
	}
	var name, value string
	strlist := mySplit(paramstr, "&")
	for _, s := range strlist {
		if s == "" {
			continue
		}
		words := strings.SplitN(s, "=", 2)
		switch len(words) {
		case 1:
			name = words[0]
			value = ""
		case 2:
			name = words[0]
			value = words[1]
		default:
			name = ""
			value = ""
		}
		// fmt.Printf("in fakeParams, name=%s, value=%s\n", name, value)
		ret[name] = value
	}
	return ret
}

type idclause_TC struct {
	paramstr string
	xres string
	xids string
	xsucc bool
}

var idclause_Tab = []idclause_TC {
	{ "id_field=id&id=123", "WHERE id = ?", "123", true },
	{ "id_field=id&ids=123", "WHERE id in (?)", "123", true },
	{ "id_field=id&ids=123,456", "WHERE id in (?,?)", "123,456", true },
	{ "id_field=id", "", "", true },
}

func aToIdList(s string) []interface{} {	// nolint
	if s == "" {
		return []interface{}{}
	}
	slist := mySplit(s, ",")
	N := len(slist)
	ret := make([]interface{}, len(slist))
	for i := 0; i < N; i++ {
		ret[i] = aToIdType(slist[i])
	}
	return ret
}

// convert idlist (list of int64) to ascii csv.
func idListToA(idlist []interface{}) (string, error) {
	alist := make([]string, len(idlist))
	for i, ival := range idlist {
		val, ok := ival.(int64)
		if !ok {
			return "",
				fmt.Errorf(`idListToA conversion error on "%s"`, ival)
		}
		alist[i] = idTypeToA(val)
	}
	return strings.Join(alist, ","), nil
}

func mkIdClause_Checker(cx *testContext, tc *idclause_TC) {
	params := fakeParams(tc.paramstr)
	res, idlist := mkIdClause(params)
	cx.assertEqual(tc.xres, res, "mkIdClause query string")

	resids, err := idListToA(idlist)
	cx.assertErrorNil(err, "idListToA")
	cx.assertEqual(tc.xids, resids, "mkIdClause idlist")
}

func Test_mkIdClause(t *testing.T) {
	cx := newTestContext(t, "mkIdClause_Tab")
	for _, tc := range idclause_Tab {
		mkIdClause_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for mkIdClauseUpdate()

var mkIdClauseUpdate_Tab = []idclause_TC {
	{ "id_field=id&id=123", "WHERE id = 123", "", true },
	{ "id_field=id&ids=123", "WHERE id in (123)", "", true },
	{ "id_field=id&ids=123,456", "WHERE id in (123,456)", "", true },
	{ "id_field=id", "", "", true },
}

func mkIdClauseUpdate_Checker(cx *testContext, tc *idclause_TC) {
	params := fakeParams(tc.paramstr)
	res := mkIdClauseUpdate(params)
	cx.assertEqual(tc.xres, res, "result")
}

func Test_mkIdClauseUpdate(t *testing.T) {
	cx := newTestContext(t, "mkIdClauseUpdate_Tab")
	for _, tc := range mkIdClauseUpdate_Tab {
		mkIdClauseUpdate_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for idTypesToInterface()

type idTypesToInterface_TC struct {
	ids string
}

var idTypesToInterface_Tab = []string {
	"",
	"987",
	"654,321",
	"987,654,32,1,0",
}

func idTypesToInterface_Checker(cx *testContext, tc string) {
	alist := mySplit(tc, ",")
	res := idTypesToInterface(alist)
	str, err := idListToA(res)
	if !cx.assertErrorNil(err, "idListToA") {
		return
	}
	cx.assertEqual(tc, str, "result")
}

func Test_idTypesToInterface(t *testing.T) {
	cx := newTestContext(t, "idTypesToInterface_Tab")
	for _, tc := range idTypesToInterface_Tab {
		idTypesToInterface_Checker(cx, tc)
		cx.bump()
	}
}


// ----- unit tests for mkSelectString()

type mkSelectString_TC struct {
	paramstr string
	xres string
	xids string
	xsucc bool
}

var mkSelectString_Tab = []mkSelectString_TC {
	{"table_name=T&id_field=id&id=456&fields=a&limit=1&offset=0",
		"SELECT id,a FROM T WHERE id = ? LIMIT 1 OFFSET 0",
		"456", true},
	{"table_name=T&id_field=id&ids=123,456&fields=a,b,c&limit=1&offset=0",
		"SELECT id,a,b,c FROM T WHERE id in (?,?) LIMIT 1 OFFSET 0",
		"123,456", true},
}

// run one tc case
func mkSelectString_Checker(cx *testContext, tc *mkSelectString_TC) {
	params := fakeParams(tc.paramstr)
	res, idlist := mkSelectString(params)
	if !cx.assertEqual(tc.xres, res, "result") {
		return
	}
	ids, err := idListToA(idlist)
	if !cx.assertErrorNil(err, "idListToA") {
		return
	}
	cx.assertEqual(tc.xids, ids, "idlist")
}

func Test_mkSelectString(t *testing.T) {
	cx := newTestContext(t, "mkSelectString_Tab")
	for _, tc := range mkSelectString_Tab {
		mkSelectString_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for getBodyRecord()

type getBodyRecord_TC struct {
	data string
	keys string
	values string
}

// turn a list of strings masked as interfaces, back to list of strings.
func unmaskStrings(ilist []interface{}) []string {
	N := len(ilist)
	ret := make([]string, N)
	for i := 0; i < N; i++ {
		s, _ := ilist[i].(string)
		ret[i] = s
	}
	return ret
}

// turn a list of strings, into strings masked by interface.
func maskStrings(slist []string) []interface{} {
	N := len(slist)
	ret := make([]interface{}, N)
	for i := 0; i < N; i++ {
		ret[i] = slist[i]
	}
	return ret
}

var getBodyRecord_Tab = []getBodyRecord_TC {
	{`{"Records":[{"Keys":[], "Values":[]}]}`,
		"",
		""},
	{`{"Records":[{"Keys":["k1","k2","k3"], "Values":["v1","v2","v3"]}]}`,
		"k1,k2,k3",
		"v1,v2,v3"},
	{`{"Records":[{"Keys":["k1","k2","k3"], "Values":["v1","v2","v3"]},{"Keys":["k4","k5","k6"], "Values":["v4","v5","v6"]}]}`,
		"k1,k2,k3&k4,k5,k6",
		"v1,v2,v3&v4,v5,v6"},
}

func getBodyRecord_Checker(cx *testContext, tc *getBodyRecord_TC) {
	rdr := strings.NewReader(tc.data)
	req, _ := http.NewRequest(http.MethodPost, "/xyz", rdr)

	tckeys := strings.Split(tc.keys, "&")	// not mySplit
	tcvalues := strings.Split(tc.values, "&")  // not mySplit
	nkeys := len(tckeys)

	body, err := getBodyRecord(mkApiHandlerArg(req, nil))
	if !cx.assertErrorNil(err, "returned err") {
		return
	}

	records := body.Records
	nrecs := len(records)

	cx.assertEqual(nkeys, nrecs, "nkeys")

	for j := 0; j < nrecs; j++ {
		rec := records[j]
		keystr := strings.Join(rec.Keys, ",")
		cx.assertEqual(tckeys[j], keystr, "item keys")
		valstr := strings.Join(unmaskStrings(rec.Values), ",")
		cx.assertEqual(tcvalues[j], valstr, "item vals")
	}
}

func Test_getBodyRecord(t *testing.T) {
	cx := newTestContext(t, "getBodyRecord_Tab")
	for _, tc := range getBodyRecord_Tab {
		getBodyRecord_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for convTableNames() and grabNameField()

type convTableNames_TC struct {
	names string
	xsucc bool
}

var convTableNames_Tab = []convTableNames_TC {
	{"", true},
	{"a", true},
	{"a,b", true},
	{"abc,def,ghi", true},
	{"abc,bogus,ghi", false},
}

// mimicTableNamesQuery() returns an object that mimics the return from
// the query to the "_tables_" table.
func mimicTableNamesQuery(names []string) []*KVResponse {
	N := len(names)
	ret := make([]*KVResponse, N)
	for i := 0; i < N; i++ {
		Keys := []string{"name"}
		val := names[i]
		var ival interface{}
		if val == "bogus" {
			// an inconvertible value
			ival = interface{}(func() {})
		} else {
			ival = interface{}(val)
		}
		Values := []interface{}{ival}
		ret[i] = &KVResponse{Keys, Values, "KVResponse", ""}
	}
	return ret
}

func convTableNames_Checker(cx *testContext, tc *convTableNames_TC) {
	names := mySplit(tc.names, ",")
	obj := mimicTableNamesQuery(names)
	res, err := convTableNames(obj)
	ok := (err == nil && tc.names == strings.Join(res, ","))
	cx.assertEqual(tc.xsucc, ok, "conversion success")
}

func Test_convTableNames(t *testing.T) {
	cx := newTestContext(t, "convTableNames_Tab")
	for _, tc := range convTableNames_Tab {
		convTableNames_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for validateRecords()

type validateRecords_TC struct {
	desc string
	xsucc bool
}

// mkRecords() turns a description string into an array of records.
// the description is parsed as follows.
// record descriptions are separated by ';' chars.
// within each record, '|' separates the list of keys from the list of values.
// the keys (names) are comma-separated.
func mkRecords(desc string) []KVRecord {
	desclist := mySplit(desc, ";")
	nrecs := len(desclist)
	ret := make([]KVRecord, nrecs)
	for i, rdesc := range desclist {
		parts := mySplit(rdesc, "|")
		ret[i].Keys = mySplit(parts[0], ",")
		ret[i].Values = maskStrings(mySplit(parts[1], ","))
	}
	return ret
}

var validateRecords_Tab = []validateRecords_TC {
	{"", true},		// 0 records
	{"k1,k2,k3|v1,v2,v3", true},  // 1 record, valid
	{"k1,,k3|v1,v2,v3", false},   // 1 record, invalid
	{"k1,k2,k3|v1,v2,v3;k4|v4", true}, // 2 records, valid
	{"k1,k2,k3|v1,v2,v3;|v4", false}, // 2 records, invalid
}

func validateRecords_Checker(cx *testContext, tc *validateRecords_TC) {
	records := mkRecords(tc.desc)
	res := validateRecords(records)
	cx.assertEqual(tc.xsucc, res == nil, "result")
}

func Test_validateRecords(t *testing.T) {
	cx := newTestContext(t, "validateRecords_Tab")
	for _, tc := range validateRecords_Tab {
		validateRecords_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for convValues()

// inputs and outputs for one convValues testcase.
type convValues_TC struct {
	arg string
	xres bool
}

// table of convValues testcases.
var convValues_Tab = []convValues_TC {
	{ "", true },
	{ "abc", true },
	{ "abc,def", true },
	{ "abc,def,ghi", true },
	{ "1,def,ghi", false },
}

func strToSQLValues(arg string) []interface{} {
	args := mySplit(arg, ",")
	N := len(args)
	ret := make([]interface{}, N)
	for i, s := range args {
		j, err := strconv.Atoi(s)
		if err == nil {
			// an unconvertible value
			ret[i] = &j
		} else {
			rb := sql.RawBytes(s)
			ret[i] = &rb;
		}
	}
	return ret
}

// run one testcase for function convValues.
func convValues_Checker(cx *testContext, tc *convValues_TC) {
	argInter := strToSQLValues(tc.arg)
	err := convValues(argInter)
	ok := (err == nil &&
		tc.arg == strings.Join(unmaskStrings(argInter), ","))
	cx.assertEqual(tc.xres, ok, "conversion result")
}

// main test suite for convValues().
func Test_convValues(t *testing.T) {
	cx := newTestContext(t, "convValues_Tab")
	for _, tc := range convValues_Tab {
		convValues_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for support for testing of api calls

type apiCall_TC struct {
	title string
	hf apiHandler
	verb string
	argDesc string
	xcode int
	xdata string
}

// this special value for xdata signifies do not check the returned data.
const noCheck = `NOCHECK`

func apiCall_Checker(cx *testContext, tc *apiCall_TC) apiHandlerRet {
	log.Debugf("----- %s #%d: [%s]", cx.tabName, cx.testno, tc.title)
	result := callApiHandler(tc.hf, tc.verb, tc.argDesc)
	cx.assertEqual(tc.xcode, result.code, tc.title)
	// check the returned data only if expected data is non-nil.
	if tc.xdata != noCheck {
		rdata, err := convData(result.data)
		if !cx.assertErrorNil(err, "data conversion error") {
			return result
		}
		cx.assertEqual(tc.xdata, string(rdata), tc.title)
	}
	return result
}

func apiCalls_Runner(t *testing.T, tabName string, tab []apiCall_TC) {
	cx := newTestContext(t, tabName)
	for _, tc := range tab {
		apiCall_Checker(cx, &tc)
		cx.bump()
	}
}

func callApiHandler(hf apiHandler, verb string, desc string) apiHandlerRet {
	return hf(parseHandlerArg(verb, desc))
}

// ----- unit tests for various implemented handlers.

// note that the success or failure of a given call can be order dependent.

var createDbRecords_Tab = []apiCall_TC {
	{"create record 1",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc1","xyz1"]}]}`,
		http.StatusCreated, noCheck},
	{"create record 2",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc2","xyz2"]}]}`,
		http.StatusCreated, noCheck},
	{"create record 3",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc3","xyz3"]}]}`,
		http.StatusCreated, noCheck},
	{"create record 4",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc4","xyz4"]}]}`,
		http.StatusCreated, noCheck},

	{"get record 123",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=123|fields=name,uri`,
		http.StatusBadRequest, noCheck},
	{"get record 1",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=1|fields=name,uri`,
		http.StatusOK, noCheck},
	{"get record 2",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=2|fields=name,uri`,
		http.StatusOK, noCheck},

	{"get records 1,2",
		getDbRecordsHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles|ids=1,2&fields=name,uri`,
		http.StatusOK, noCheck},

	{"get record 1 bad field",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=123|fields=name,uri,bogus`,
		http.StatusBadRequest, noCheck},

	{"delete records 2,4",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bundles|ids=2,4`,
		http.StatusOK, noCheck},

	{"delete records no id or ids",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest, noCheck},

	{"delete record no id or ids",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest, noCheck},

	{"delete record bad table_name",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bogus|id=1`,
		http.StatusBadRequest, noCheck},

	{"get record 2 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=2`,
		http.StatusBadRequest, noCheck},

	{"get record 4 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=4`,
		http.StatusBadRequest, noCheck},

	{"delete record 1",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusOK, noCheck},

	{"delete record 1 expecting failure",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusBadRequest, noCheck},

	{"get record 1 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusBadRequest, noCheck},

	{"update records missing id",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest, noCheck},

	{"update records missing table_name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|id=1`,
		http.StatusBadRequest, noCheck},

	{"update records bogus table_name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|table_name=bogus&id=1||{"records":[{"keys":["name", "uri"], "values":["name9", "uri9"]}]}`,
		http.StatusBadRequest, noCheck},

	{"update records bogus field name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|table_name=xxx&id=1||{"records":[{"keys":["bogus", "uri"], "values":["name9", "uri9"]}]}`,
		http.StatusBadRequest, noCheck},

	{"update records no body records",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|table_name=xxx&id=1||{"records":[]}`,
		http.StatusBadRequest, noCheck},

	{"update record missing id",
		updateDbRecordHandler,
		http.MethodPatch,
		`/test/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest, noCheck},

	{"create record missing id",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest, noCheck},

	{"create records missing body",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles|id=1`,
		http.StatusBadRequest, noCheck},

	{"create records missing table_name",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|id=1`,
		http.StatusBadRequest, noCheck},

	{"create records bogus field",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","bogus"],"Values":["abc3","xyz3"]}]}`,
		http.StatusBadRequest, noCheck},

	{"get records missing table_name",
		getDbRecordsHandler,
		http.MethodGet,
		`/test/db/_table/tabname|id=1`,
		http.StatusBadRequest, noCheck},

	{"get record missing table_name",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|id=1`,
		http.StatusBadRequest, noCheck},

	{"delete records missing table_name",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/test/db/_table/tabname||ids=1`,
		http.StatusBadRequest, noCheck},

	{"delete record missing table name",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname/1234|id=1`,
		http.StatusBadRequest, noCheck},

	{"delete record nonexistent record",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/tabname/1234|id=1001`,
		http.StatusBadRequest, noCheck},

	{"create records with excess values",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=xxx||{"Records":[{"Keys":["name","uri"],"Values":["abc4","xyz4","superfluous"]}]}`,
		http.StatusBadRequest, noCheck},
}

// the handlers must be called in a certain order, in order for the
// calls to succeed or fail as expected.
func Test_createDbRecordsHandler(t *testing.T) {
	apiCalls_Runner(t, "createDbRecords_Tab", createDbRecords_Tab)
}

// ----- unit tests of getDbTablesHandler()

var getDbTables_Tab = []apiCall_TC {
	{"get tablenames",
		getDbTablesHandler,
		http.MethodGet,
		`http://localhost/test/db/_tables`,
		http.StatusOK, noCheck},
}

func getDbTables_Checker(cx *testContext, tc *apiCall_TC) {
	result := apiCall_Checker(cx, tc)
	if result.code != tc.xcode {
		// would have already failed.
		return
	}
	// if the code was success, data should be of this type.
	data, ok := result.data.(TablesResponse)
	cx.assertTrue(ok, "TablesResponse data type")
	xtabNames := "bundles,nothing,users"
	dataNames := data.Names
	sort.Strings(dataNames)
	resNames := strings.Join(dataNames, ",")
	cx.assertEqualObj(xtabNames, resNames, "retrieved names")
	cx.assertEqual("TablesResponse", data.Kind, "tables response kind")
	w := strings.SplitN(data.Self, "?", 2)
	cx.assertEqual(tc.argDesc, w[0], "tables response self")
}

func Test_getDbTablesHandler(t *testing.T) {
	cx := newTestContext(t, "getDbTables_Tab")
	for _, tc := range getDbTables_Tab {
		getDbTables_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit test for tablesQuery()

type tablesQuery_TC struct {
	self string
	tableName string
	fieldName string
	xcode int
}

var tablesQuery_Tab = []tablesQuery_TC {
	{"xyz", "bogus_table", "name", http.StatusBadRequest},
	{"xyz", "_tables_", "bogus_field", http.StatusBadRequest},
	{"xyz", "_tables_", "name", http.StatusOK},
}

func tablesQuery_Checker(cx *testContext, tc *tablesQuery_TC) {
	result := tablesQuery(tc.self, tc.tableName, tc.fieldName)
	cx.assertEqual(tc.xcode, result.code, "returned code")
}

func Test_tablesQuery_bogusField(t *testing.T) {
	cx := newTestContext(t, "tablesQuery_Tab")
	for _, tc := range tablesQuery_Tab {
		tablesQuery_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests of updateDbRecordHandler()

func Test_updateDbRecordHandler(t *testing.T) {

	cx := newTestContext(t)
	tabName := "xxx"
	recno := "2"
	newurl := "host9:xyz"

	argDesc := fmt.Sprintf(`/test/db/_table/tabname|table_name=%s&id=%s|fields=name|{"records":[{"keys":["name", "uri"], "values":["name9", "%s"]}]}`,
		tabName, recno, newurl)

	// do an update record
	tc := apiCall_TC{"update record in xxx",
		updateDbRecordHandler,
		http.MethodPatch,
		argDesc,
		http.StatusOK, noCheck}

	result := apiCall_Checker(cx, &tc)
	if result.code != tc.xcode {
		// would have already failed.
		return
	}

	// if the code was success, data should be of this type.
	data, ok := result.data.(NumChangedResponse)
	if !cx.assertTrue(ok, "data of wrong type") {
		return
	}
	if !cx.assertEqual(1, int(data.NumChanged),
			"number of changed records") {
		return
	}

	vals, ok := retrieveValues(cx, tabName, recno)
	if !cx.assertTrue(ok, "return from retrieveValues") {
		return
	}

	cx.assertEqual(newurl, vals[2], "retrieved url")
}

// return the values of the given row in the given table.
func retrieveValues(cx *testContext,
		tabName string,
		recno string) ([]string, bool) {
	argDesc := fmt.Sprintf(`/test/db/_table/tabname|table_name=%s&id=%s`,
		tabName, recno)
	tc := apiCall_TC{"get record in retrieveValues",
		getDbRecordHandler,
		http.MethodGet,
		argDesc,
		http.StatusOK, noCheck}
	result := apiCall_Checker(cx, &tc)
	if !cx.assertEqual(tc.xcode, result.code,
		"return code from getDbRecordHandler") {
		return nil, false
	}

	// fetch the changed record
	rdata, ok := result.data.(RecordsResponse)
	if !cx.assertTrue(ok, "data should be of type RecordsResponse") {
		return nil, false
	}

	recs := rdata.Records
	nr := len(recs)
	if !cx.assertEqual(1, nr, "number of records") {
		return nil, false
	}
	ivals := recs[0].Values
	_ = convValues(ivals)
	return unmaskStrings(ivals), true
}

// ----- unit tests for updateDbRecordsHandler()

func Test_updateDbRecordsHandler(t *testing.T) {

	cx := newTestContext(t)
	tabName := "xxx"
	recno := "1"
	newname := "name7"
	newurl := "host7:abc"

	argDesc := fmt.Sprintf(`/test/db/_table/tabname|table_name=%s|ids=%s&fields=name|{"records":[{"keys":["name", "uri"], "values":["%s", "%s"]}]}`,
		tabName, recno, newname, newurl)

	// do an update record
	tc := apiCall_TC{"update record",
		updateDbRecordsHandler,
		http.MethodPatch,
		argDesc,
		http.StatusOK, noCheck}

	result := apiCall_Checker(cx, &tc)
	if result.code != tc.xcode {
		// would have already failed.
		return
	}

	// if the code was success, data should be of this type.
	data, ok := result.data.(NumChangedResponse)
	if !cx.assertTrue(ok, "data should be of type NumChanedResponse") {
		return
	}
	if !cx.assertEqual(1, int(data.NumChanged),
			"number changed") {
		return
	}

	argDesc = fmt.Sprintf(`/test/db/_table/tabname|table_name=%s&id=%s`,
		tabName, recno)

	// read the record back, and check the data
	tc = apiCall_TC{"get record",
		getDbRecordHandler,
		http.MethodGet,
		argDesc,
		http.StatusOK, noCheck}
	result = apiCall_Checker(cx, &tc)
	if tc.xcode != result.code {
		// would have already failed.
		return
	}

	vals, ok := retrieveValues(cx, tabName, recno)
	if !ok {
		return
	}

	cx.assertEqual(newurl, vals[2], "retrieved url field")
}

// ----- unit tests for getDbRecordsHandler()

func nameRange(start int, n int) []string {
	nameList := make([]string, n)
	for i := 0; i < n; i++ {
		nameList[i] = fmt.Sprintf("x%d", i+start)
	}
	return nameList
}

func readNamesWithOffset(cx *testContext,
		tab string,
		offset int) []string {
	ret := []string{}
	argDesc := fmt.Sprintf(`/test/db/_table|table_name=%s|fields=name&offset=%d`,
		tab, offset)
	result := callApiHandler(getDbRecordsHandler, http.MethodGet, argDesc)
	if !cx.assertEqual(http.StatusOK, result.code,
		"return code from getDbRecordsHandler") {
		return ret
	}

	resp, ok := result.data.(RecordsResponse)
	if !cx.assertTrue(ok, "data of type RecordsResponse") {
		return ret
	}

	// grab the name field from each record
	ret = make([]string, len(resp.Records))
	for i, rec := range resp.Records {
		_ = convValues(rec.Values)
		svals := unmaskStrings(rec.Values)
		ret[i] = svals[0]
	}

	return ret
}

func Test_getDbRecordHandler_offset(t *testing.T) {
	// the table "toomany" has been seeded with more than maxRecs
	// records that can be read.
	cx := newTestContext(t)
	tableName := "toomany"

	// read first batch, checking values
	names := readNamesWithOffset(cx, tableName, 0)
	xnames := nameRange(1, len(names))
	cx.assertEqualObj(xnames, names, "names from #1 batch")

	// expect exactly maxRecs results
	nrecs := len(names)
	cx.assertEqual(maxRecs, nrecs, "nrecs")

	// read second batch, checking values
	names = readNamesWithOffset(cx, tableName, maxRecs)
	xnames = nameRange(maxRecs+1, len(names))
	cx.assertEqualObj(xnames, names, "names from #2 batch")
}

// ----- unit tests for createDbTableHandler()

var users_schema = `{"fields":[{"name":"id","properties":["is_primary_key","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}`

// table of createDbTable testcases.
var createDbTable_Tab = []apiCall_TC {
	{"create table w/ missing table_name",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema|||`+users_schema,
		http.StatusBadRequest, noCheck},
	{"create table w/ invalid table_name",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema|table_name=XYZ.DEF||`+users_schema,
		http.StatusBadRequest, noCheck},
	{"create table w/ malformed body",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABC|table_name=ABC||bogus`+users_schema,
		http.StatusBadRequest, noCheck},
	{"create table ABC expecting success",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABC|table_name=ABC||`+users_schema,
		http.StatusCreated, noCheck},
	{"create table GHI expecting success, absent properties",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/GHI|table_name=GHI||{"fields":[{"name":"id","properties":["is_primary_key","int32"]},{"name":"uri"},{"name":"name"}]}`,
		http.StatusCreated, noCheck},
	{"create table ABC expecting failure, pre-existing",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABC|table_name=ABC||`+users_schema,
		http.StatusBadRequest, noCheck},
	{"create record in ABC",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=ABC||{"Records":[{"Keys":["name","uri"],"Values":["xyz-abc4","abc-xyz4"]}]}`,
		http.StatusCreated, noCheck},
	{"get record 1 in ABC",
		getDbRecordHandler,
		http.MethodGet,
		`/test/db/_table/tabname|table_name=ABC&id=1|fields=name,uri`,
		http.StatusOK, noCheck},
}

// the createDbTable test suite.  run all createDbTable testcases.
func Test_createDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, "createDbTable_Tab", createDbTable_Tab)
}

// ----- unit tests for deleteDbTableHandler()

// table of deleteDbTableHandler testcases.
var deleteDbTable_Tab = []apiCall_TC {
	{"delete table ABC missing table_name",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema`,
		http.StatusBadRequest, noCheck},
	{"delete table ABC empty table_name",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema|table_name=`,
		http.StatusBadRequest, noCheck},
	{"delete table ABCD expecting failure",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/ABCD|table_name=ABCD`,
		http.StatusBadRequest, noCheck},
	{"create table ABCD expecting success",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABCD|table_name=ABCD||{"fields":[{"name":"id","properties":["is_primary_key","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}`,
		http.StatusCreated, noCheck},
	{"delete table ABCD expecting success",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/ABCD|table_name=ABCD`,
		http.StatusOK, noCheck},
	{"delete table ABCD expecting failure",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/ABCD|table_name=ABCD`,
		http.StatusBadRequest, noCheck},
}

// the deleteDbTable test suite.  run all deleteDbTable testcases.
func Test_deleteDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, "deleteDbTable_Tab", deleteDbTable_Tab)
}

// ----- unit tests for schemaQuery().

// inputs and outputs for one schemaQuery testcase.
type schemaQuery_TC struct {
	self string
	tableName string
	fieldName string
	selector string
	item string
	xcode int
	xdata string
}

// table of schemaQuery testcases.
var schemaQuery_Tab = []schemaQuery_TC {
	// a good request
	{ "http://abc", "_tables_", "schema", "name", "users", http.StatusOK, "{users_schema SchemaResponse http://abc}" },

	// a good request
	{ "http://abc", "_tables_", "schema", "name", "bundles", http.StatusOK, "{bundles_schema SchemaResponse http://abc}" },

	// bogus table
	{ "http://abc", "bogus", "schema", "name", "users", http.StatusBadRequest, "xxx" },

	// bogus field
	{ "http://abc", "_tables_", "schema", "bogus", "users", http.StatusBadRequest, "xxx" },

	// bogus item
	{ "http://abc", "_tables_", "schema", "name", "bogus", http.StatusBadRequest, "xxx" },
}

// run one testcase for function schemaQuery.
func schemaQuery_Checker(cx *testContext, tc *schemaQuery_TC) {
	res := schemaQuery(tc.self, tc.tableName, tc.fieldName,
			tc.selector, tc.item)
	cx.assertEqual(tc.xcode, res.code, "returned code")
	if tc.xcode == http.StatusOK {
		dataStr := fmt.Sprintf("%v", res.data)
		cx.assertEqual(tc.xdata, dataStr, "returned data")
	}
}

// the schemaQuery test suite.  run all schemaQuery testcases.
func Test_schemaQuery(t *testing.T) {
	cx := newTestContext(t, "schemaQuery_Tab")
	for _, tc := range schemaQuery_Tab {
		schemaQuery_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit tests for describeDbTableHandler().

// table of describeDbTable testcases.
var describeDbTable_Tab = []apiCall_TC {
	{"get schema for users table",
		describeDbTableHandler,
		http.MethodGet,
		`/test/db/_schema/users|table_name=users`,
		http.StatusOK, noCheck},
	{"get schema for bundles table",
		describeDbTableHandler,
		http.MethodGet,
		`/test/db/_schema/bundles|table_name=bundles`,
		http.StatusOK, noCheck},
	{"get schema for bogus table",
		describeDbTableHandler,
		http.MethodGet,
		`/test/db/_schema/bogus|table_name=bogus`,
		http.StatusBadRequest, noCheck},
	{"get schema for no table_name",
		describeDbTableHandler,
		http.MethodGet,
		`/test/db/_schema/|`,
		http.StatusBadRequest, noCheck},
}

// the describeDbTable test suite.  run all describeDbTable testcases.
func Test_describeDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, "describeDbTable_Tab", describeDbTable_Tab)
}

// ----- unit tests for execN()

// execN() is pretty well tested thru API test cases.
// this test is only to exercise an error condition
// which can't be easily reproduced thru the API.
func Test_execN(t *testing.T) {
	cx := newTestContext(t)
	err := execN(mkBadDb())
	cx.assertTrue(err != nil, "expected error")
}

// ----- unit tests for runExec()

// runExec() is pretty well tested thru API test cases.
// this test is only to exercise an error condition
// which can't be easily reproduced thru the API.
func Test_runExec(t *testing.T) {
	cx := newTestContext(t)
	query := "select * from bundles where name = ?"
	values := []interface{}{}  // insufficient values
	_, err := runExec(db, query, values)
	cx.assertTrue(err != nil, "expected error")
}

// ----- unit tests for getDbResourcesHandler().

// table of getDbResources testcases.
var getDbResources_Tab = []apiCall_TC {
	{"get db resources",
		getDbResourcesHandler,
		http.MethodGet,
		`/test/db`,
		http.StatusOK, noCheck},
}

// the getDbResources test suite.  run all getDbResources testcases.
func Test_getDbResourcesHandler(t *testing.T) {
	apiCalls_Runner(t, "getDbResources_Tab", getDbResources_Tab)
}

// ----- unit tests for deleteDbRecordHandler().

// this test suite is somewhat minimal since deleteDbRecordHandler
// gets more thorough testing by the unit tests for createDbRecordsHandler.

// table of deleteDbRecord testcases.
var deleteDbRecord_Tab = []apiCall_TC {
	{"create table xxxdel",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABC|table_name=xxxdel||`+users_schema,
		http.StatusCreated, noCheck},
	{"create db resources",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=xxxdel||{"Records":[{"Keys":["name","uri"],"Values":["name-1","uri-1"]}]}`,
		http.StatusCreated, `{"Ids":[1],"Kind":"Collection"}`},
	{"delete db resources expecting success",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/test/db/_table/xxx|table_name=xxxdel&id=1`,
		http.StatusOK, `{"NumChanged":1,"Kind":"NumChangedResponse"}`},
	{"delete table xxxdel",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/ABCD|table_name=xxxdel`,
		http.StatusOK, noCheck},
}

// the deleteDbRecord test suite.  run all deleteDbRecord testcases.
func Test_deleteDbRecordHandler(t *testing.T) {
	apiCalls_Runner(t, "deleteDbRecord_Tab", deleteDbRecord_Tab)
}

// ----- unit tests for deleteDbRecordsHandler().

// this test suite is somewhat minimal since deleteDbRecordsHandler
// gets more thorough testing by the unit tests for createDbRecordsHandler.

// table of deleteDbRecords testcases.
var deleteDbRecords_Tab = []apiCall_TC {
	{"setup: create table xxxdels",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/ABC|table_name=xxxdels||`+users_schema,
		http.StatusCreated, noCheck},
	{"setup: create db record 1",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=xxxdels||{"Records":[{"Keys":["name","uri"],"Values":["xyz-abc5","abc-xyz5"]}]}`,
		http.StatusCreated, noCheck},
	{"create db record 2",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/tabname|table_name=xxxdels||{"Records":[{"Keys":["name","uri"],"Values":["xxxdel1","abc-xyz5"]}]}`,
		http.StatusCreated, noCheck},
	{"delete db records 1,2",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/test/db/_table/xxx|table_name=xxxdels|ids=1,2`,
		http.StatusOK, `{"NumChanged":2,"Kind":"NumChangedResponse"}`},
	{"teardown: delete table xxxdels",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/xxxdels|table_name=xxxdels`,
		http.StatusOK, noCheck},
}

// the deleteDbRecords test suite.  run all deleteDbRecords testcases.
func Test_deleteDbRecordsHandler(t *testing.T) {
	apiCalls_Runner(t, "deleteDbRecords_Tab", deleteDbRecords_Tab)
}

// ----- unit tests for getDbRecordHandler().

// this suite is somewhat miminal since getDbRecordHandler is well
// exercised by the test suite for createDbRecordsHandler.

// table of getDbRecord testcases.
var getDbRecordHandler_Tab = []apiCall_TC {
	{"setup: create table xxxget",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/xxxget|table_name=xxxget||`+users_schema,
		http.StatusCreated, noCheck},
	{"setup: create db record 1",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/xxxget|table_name=xxxget||{"Records":[{"Keys":["uri","name"],"Values":["uri-a","name-a"]}]}`,
		http.StatusCreated, noCheck},
	{"get db record 1",
		getDbRecordHandler,
		http.MethodGet,
		`http://localhost/test/db/_table/xxxget|table_name=xxxget&id=1`,
		http.StatusOK, `{"Records":[{"Keys":["id","uri","name"],"Values":["1","uri-a","name-a"],"Kind":"KVResponse","Self":"http://localhost/test/db/_table/xxxget/1"}],"Kind":"Collection"}`},
	{"teardown: delete table xxxget",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/xxxget|table_name=xxxget`,
		http.StatusOK, noCheck},
}

// the getDbRecord test suite.  run all getDbRecord testcases.
func Test_getDbRecordHandler(t *testing.T) {
	apiCalls_Runner(t, "getDbRecordHandler_Tab", getDbRecordHandler_Tab)
}

// ----- unit tests for getDbRecordHandler().

// this suite is somewhat miminal since getDbRecordsHandler is well
// exercised by the test suite for createDbRecordsHandler.

// table of getDbRecords testcases.
var getDbRecordsHandler_Tab = []apiCall_TC {
	{"setup: create table xxxget",
		createDbTableHandler,
		http.MethodPost,
		`/test/db/_schema/xxxget|table_name=xxxget||`+users_schema,
		http.StatusCreated, noCheck},
	{"setup: create db record 1",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/xxxget|table_name=xxxget||{"Records":[{"Keys":["name","uri"],"Values":["name-a","uri-a"]}]}`,
		http.StatusCreated, noCheck},
	{"setup: create db record 2",
		createDbRecordsHandler,
		http.MethodPost,
		`/test/db/_table/xxxget|table_name=xxxget||{"Records":[{"Keys":["uri","name"],"Values":["uri-b","name-b"]}]}`,
		http.StatusCreated, noCheck},
	{"get db records 1,2",
		getDbRecordsHandler,
		http.MethodGet,
		`http://localhost/db/_table/xxxget|table_name=xxxget|ids=1,2`,
		http.StatusOK,
		`{"Records":[{"Keys":["id","uri","name"],"Values":["1","uri-a","name-a"],"Kind":"KVResponse","Self":"http://localhost/test/db/_table/xxxget/1"},{"Keys":["id","uri","name"],"Values":["2","uri-b","name-b"],"Kind":"KVResponse","Self":"http://localhost/test/db/_table/xxxget/2"}],"Kind":"Collection"}`},
	{"teardown: delete table xxxget",
		deleteDbTableHandler,
		http.MethodDelete,
		`/test/db/_schema/xxxget|table_name=xxxget`,
		http.StatusOK, noCheck},
}

// the getDbRecords test suite.  run all getDbRecords testcases.
func Test_getDbRecordsHandler(t *testing.T) {
	apiCalls_Runner(t, "getDbRecordsHandler_Tab", getDbRecordsHandler_Tab)
}

// ----- unit tests for listToMap().

// inputs and outputs for one listToMap testcase.
type listToMap_TC struct {
	arg string
}

// table of listToMap testcases.
var listToMap_Tab = []listToMap_TC {
	{""},
	{"a"},
	{"a,b"},
	{"x,y,z"},
}

// run one testcase for function listToMap.
func listToMap_Checker(cx *testContext, tc *listToMap_TC) {
	m := mySplit(tc.arg, ",")
	result := listToMap(m)
	cx.assertEqual(len(m), len(result), "number of items")
	for _, mitem := range m {
		cx.assertTrue(result[mitem] != 0, "item in map")
	}
}

// the listToMap test suite.  run all listToMap testcases.
func Test_listToMap(t *testing.T) {
	cx := newTestContext(t, "listToMap_Tab")
	for _, tc := range listToMap_Tab {
		listToMap_Checker(cx, &tc)
		cx.bump()	// increment testno.
	}
}
