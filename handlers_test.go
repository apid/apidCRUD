package apidCRUD

import (
	"testing"
	"fmt"
	"strconv"
	"strings"
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

// ----- unit tests for notImplemented()

func Test_notImplemented(t *testing.T) {
	cx := newTestContext(t)
	res := notImplemented()
	cx.assertEqual(http.StatusNotImplemented, res.code, "returned code")
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
		"SELECT a FROM T WHERE id = ? LIMIT 1 OFFSET 0",
		"456", true},
	{"table_name=T&id_field=id&ids=123,456&fields=a,b,c&limit=1&offset=0",
		"SELECT a,b,c FROM T WHERE id in (?,?) LIMIT 1 OFFSET 0",
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
// the query to the "tables" table.
func mimicTableNamesQuery(names []string) []*KVRecord {
	N := len(names)
	ret := make([]*KVRecord, N)
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
		ret[i] = &KVRecord{Keys, Values}
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
}

func apiCall_Checker(cx *testContext, tc *apiCall_TC) apiHandlerRet {
	log.Debugf("----- %s #%d: [%s]", cx.tabName, cx.testno, tc.title)
	result := callApiHandler(tc.hf, tc.verb, tc.argDesc)
	cx.assertEqual(tc.xcode, result.code, tc.title)
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

// ----- unit tests for not-implemented handlers.

var notimpl_Tab = []apiCall_TC {
	{"API not implemented",
		getDbResourcesHandler,
		http.MethodGet,
		"/db",
		http.StatusNotImplemented},
	{"API not implemented",
		getDbSchemasHandler,
		http.MethodGet,
		"/db/_schema",
		http.StatusNotImplemented},
	{"API not implemented",
		updateDbTablesHandler,
		http.MethodPatch,
		"/db/_schema",
		http.StatusNotImplemented},
	{"API not implemented",
		describeDbTableHandler,
		http.MethodGet,
		"/db/_schema/tabname",
		http.StatusNotImplemented},
	{"API not implemented",
		describeDbFieldHandler,
		http.MethodDelete,
		"/db/_schema/tabname",
		http.StatusNotImplemented},
}

func Test_notimpl(t *testing.T) {
	apiCalls_Runner(t, "notimpl_Tab", notimpl_Tab)
}

// ----- unit tests for various implemented handlers.

// note that the success or failure of a given call can be order dependent.

var createDbRecords_Tab = []apiCall_TC {
	{"create record 1",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc1","xyz1"]}]}`,
		http.StatusCreated},
	{"create record 2",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc2","xyz2"]}]}`,
		http.StatusCreated},
	{"create record 3",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc3","xyz3"]}]}`,
		http.StatusCreated},
	{"create record 4",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","uri"],"Values":["abc4","xyz4"]}]}`,
		http.StatusCreated},

	{"get record 123",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=123|fields=name,uri`,
		http.StatusBadRequest},
	{"get record 1",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=1|fields=name,uri`,
		http.StatusOK},
	{"get record 2",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=2|fields=name,uri`,
		http.StatusOK},

	{"get records 1,2",
		getDbRecordsHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles|ids=1,2&fields=name,uri`,
		http.StatusOK},

	{"get record 1 bad field",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=123|fields=name,uri,bogus`,
		http.StatusBadRequest},

	{"delete records 2,4",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bundles|ids=2,4`,
		http.StatusOK},

	{"delete records no id or ids",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest},

	{"delete record no id or ids",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest},

	{"delete record bad table_name",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bogus|id=1`,
		http.StatusBadRequest},

	{"get record 2 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=2`,
		http.StatusBadRequest},

	{"get record 4 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=4`,
		http.StatusBadRequest},

	{"delete record 1",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusOK},

	{"delete record 1 expecting failure",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusBadRequest},

	{"get record 1 expecting failure",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=bundles&id=1`,
		http.StatusBadRequest},

	{"update records missing id",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest},

	{"update records missing table_name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/db/_table/tabname|id=1`,
		http.StatusBadRequest},

	{"update records bogus table_name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/db/_table/tabname|table_name=bogus&id=1||{"records":[{"keys":["name", "uri"], "values":["name9", "uri9"]}]}`,
		http.StatusBadRequest},

	{"update records bogus field name",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/db/_table/tabname|table_name=xxx&id=1||{"records":[{"keys":["bogus", "uri"], "values":["name9", "uri9"]}]}`,
		http.StatusBadRequest},

	{"update records no body records",
		updateDbRecordsHandler,
		http.MethodPatch,
		`/db/_table/tabname|table_name=xxx&id=1||{"records":[]}`,
		http.StatusBadRequest},

	{"update record missing id",
		updateDbRecordHandler,
		http.MethodPatch,
		`/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest},

	{"create record missing id",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles`,
		http.StatusBadRequest},

	{"create records missing body",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles|id=1`,
		http.StatusBadRequest},

	{"create records missing table_name",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|id=1`,
		http.StatusBadRequest},

	{"create records bogus field",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=bundles||{"Records":[{"Keys":["name","bogus"],"Values":["abc3","xyz3"]}]}`,
		http.StatusBadRequest},

	{"get records missing table_name",
		getDbRecordsHandler,
		http.MethodGet,
		`/db/_table/tabname|id=1`,
		http.StatusBadRequest},

	{"get record missing table_name",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|id=1`,
		http.StatusBadRequest},

	{"delete records missing table_name",
		deleteDbRecordsHandler,
		http.MethodDelete,
		`/db/_table/tabname||ids=1`,
		http.StatusBadRequest},

	{"delete record missing table name",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname/1234|id=1`,
		http.StatusBadRequest},

	{"delete record nonexistent record",
		deleteDbRecordHandler,
		http.MethodDelete,
		`/db/_table/tabname/1234|id=1001`,
		http.StatusBadRequest},

	{"create records with excess values",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=xxx||{"Records":[{"Keys":["name","uri"],"Values":["abc4","xyz4","superfluous"]}]}`,
		http.StatusBadRequest},
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
		"/db/_tables",
		http.StatusOK},
}

func getDbTablesHandler_Checker(cx *testContext, tc *apiCall_TC) {
	result := apiCall_Checker(cx, tc)
	if result.code != tc.xcode {
		// would have already failed.
		return
	}
	// if the code was success, data should be of this type.
	data, ok := result.data.(TablesResponse)
	cx.assertTrue(ok, "TablesResponse data type")
	xtabnames := "bundles,users,nothing"
	resnames := strings.Join(data.Names, ",")
	cx.assertEqualObj(xtabnames, resnames, "retrieved names")
}

func Test_getDbTablesHandler(t *testing.T) {
	cx := newTestContext(t, "getDbTables_Tab")
	for _, tc := range getDbTables_Tab {
		getDbTablesHandler_Checker(cx, &tc)
		cx.bump()
	}
}

// ----- unit test for tablesQuery()

type tablesQuery_TC struct {
	tableName string
	fieldName string
	xcode int
}

var tablesQuery_Tab = []tablesQuery_TC {
	{"bogus_table", "name", http.StatusBadRequest},
	{"tables", "bogus_field", http.StatusBadRequest},
	{"tables", "name", http.StatusOK},
}

func tablesQuery_Checker(cx *testContext, tc *tablesQuery_TC) {
	harg := parseHandlerArg(http.MethodGet, `/db/_tables`)
	result := tablesQuery(harg, tc.tableName, tc.fieldName)
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

	argDesc := fmt.Sprintf(`/db/_table/tabname|table_name=%s&id=%s|fields=name|{"records":[{"keys":["name", "uri"], "values":["name9", "%s"]}]}`,
		tabName, recno, newurl)

	// do an update record
	tc := apiCall_TC{"update record in xxx",
		updateDbRecordHandler,
		http.MethodPatch,
		argDesc,
		http.StatusOK}

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
	argDesc := fmt.Sprintf(`/db/_table/tabname|table_name=%s&id=%s`,
		tabName, recno)
	tc := apiCall_TC{"get record in retrieveValues",
		getDbRecordHandler,
		http.MethodGet,
		argDesc,
		http.StatusOK}
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

	argDesc := fmt.Sprintf(`/db/_table/tabname|table_name=%s|ids=%s&fields=name|{"records":[{"keys":["name", "uri"], "values":["%s", "%s"]}]}`,
		tabName, recno, newname, newurl)

	// do an update record
	tc := apiCall_TC{"update record",
		updateDbRecordsHandler,
		http.MethodPatch,
		argDesc,
		http.StatusOK}

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

	argDesc = fmt.Sprintf(`/db/_table/tabname|table_name=%s&id=%s`,
		tabName, recno)

	// read the record back, and check the data
	tc = apiCall_TC{"get record",
		getDbRecordHandler,
		http.MethodGet,
		argDesc,
		http.StatusOK}
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
	argDesc := fmt.Sprintf(`/db/_table|table_name=%s|fields=name&offset=%d`,
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

// table of createDbTable testcases.
var createDbTable_Tab = []apiCall_TC {
	{"create table w/ missing table_name",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema|||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create table w/ invalid table_name",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema|table_name=XYZ.DEF||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create table w/ malformed body",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema|||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name"}]}`,
		http.StatusBadRequest},
	{"create table ABC excess tables in body",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema/ABC|table_name=ABC||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]},{"name":"DEF","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create table ABC expecting success",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema/ABC|table_name=ABC||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusCreated},
	{"create table ABC expecting failure",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema/ABC|table_name=ABC||{"resource":[{"name":"ABC","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create record in ABC",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table/tabname|table_name=ABC||{"Records":[{"Keys":["name","uri"],"Values":["xyz-abc4","abc-xyz4"]}]}`,
		http.StatusCreated},
	{"get record 1 in ABC",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=ABC&id=1|fields=name,uri`,
		http.StatusOK},
}

// the createDbTable test suite.  run all createDbTable testcases.
func Test_createDbTable(t *testing.T) {
	apiCalls_Runner(t, "createDbTable_Tab", createDbTable_Tab)
}

// ----- unit tests for deleteDbTableHandler()

// table of deleteDbTableHandler testcases.
var deleteDbTable_Tab = []apiCall_TC {
	{"delete table ABC missing table_name",
		deleteDbTableHandler,
		http.MethodDelete,
		`/db/_schema`,
		http.StatusBadRequest},
	{"delete table ABC empty table_name",
		deleteDbTableHandler,
		http.MethodDelete,
		`/db/_schema|table_name=`,
		http.StatusBadRequest},
	{"delete table ABCD expecting failure",
		deleteDbTableHandler,
		http.MethodDelete,
		`/db/_schema/ABCD|table_name=ABCD`,
		http.StatusBadRequest},
	{"create table ABCD expecting success",
		createDbTableHandler,
		http.MethodPost,
		`/db/_schema/ABCD|table_name=ABCD||{"resource":[{"name":"ABCD","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusCreated},
	{"delete table ABCD expecting success",
		deleteDbTableHandler,
		http.MethodDelete,
		`/db/_schema/ABCD|table_name=ABCD`,
		http.StatusOK},
	{"delete table ABCD expecting failure",
		deleteDbTableHandler,
		http.MethodDelete,
		`/db/_schema/ABCD|table_name=ABCD`,
		http.StatusBadRequest},
}

// the deleteDbTable test suite.  run all deleteDbTable testcases.
func Test_deleteDbTable(t *testing.T) {
	apiCalls_Runner(t, "deleteDbTable_Tab", deleteDbTable_Tab)
}

// ----- unit tests for createDbTablesHandler()

// table of createDbTables testcases.
var createDbTables_Tab = []apiCall_TC {
	{"create tables w/ invalid table_name",
		createDbTablesHandler,
		http.MethodPost,
		`/db/_schema|||{"resource":[{"name":"GHI.JKL","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create tables w/ malformed body",
		createDbTablesHandler,
		http.MethodPost,
		`/db/_schema|||{"resource":[{"name":"GHI","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name"}]}}`,
		http.StatusBadRequest},
	{"create tables GHI and JKL expecting success",
		createDbTablesHandler,
		http.MethodPost,
		`/db/_schema/GHI|table_name=GHI||{"resource":[{"name":"GHI","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]},{"name":"JKL","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusCreated},
	{"create tables GHI expecting failure",
		createDbTablesHandler,
		http.MethodPost,
		`/db/_schema/GHI|table_name=GHI||{"resource":[{"name":"GHI","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create tables JKL expecting failure",
		createDbTablesHandler,
		http.MethodPost,
		`/db/_schema|table_name=JKL||{"resource":[{"name":"JKL","fields":[{"name":"id","properties":["primary","int32"]},{"name":"uri","properties":[]},{"name":"name","properties":[]}]}]}`,
		http.StatusBadRequest},
	{"create record in GHI",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table|table_name=GHI||{"Records":[{"Keys":["name","uri"],"Values":["xyz-abc4","abc-xyz4"]}]}`,
		http.StatusCreated},
	{"get record 1 in GHI",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=GHI&id=1|fields=name,uri`,
		http.StatusOK},
	{"create record in JKL",
		createDbRecordsHandler,
		http.MethodPost,
		`/db/_table|table_name=JKL||{"Records":[{"Keys":["name","uri"],"Values":["xyz-abc4","abc-xyz4"]}]}`,
		http.StatusCreated},
	{"get record 1 in JKL",
		getDbRecordHandler,
		http.MethodGet,
		`/db/_table/tabname|table_name=JKL&id=1|fields=name,uri`,
		http.StatusOK},
}

// the createDbTable test suite.  run all createDbTable testcases.
func Test_createDbTables(t *testing.T) {
	apiCalls_Runner(t, "createDbTables_Tab", createDbTables_Tab)
}
