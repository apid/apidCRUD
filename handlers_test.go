package apidCRUD

import (
	"testing"
	"fmt"
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

func mkSQLRow_Checker(t *testing.T, i int, N int) {
	fname := "mkSQLRow"
	res := mkSQLRow(N)
	if len(res) != N {
		t.Errorf("#%d: %s(%d) failed", i, fname, N)
		return
	}
	for _, v := range res {
		_, ok := v.(*sql.RawBytes)
		if !ok {
			t.Errorf("#%d: %s(%d) sql conversion error", i, fname, N)
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
	fname := "notImplemented"
	xcode := http.StatusNotImplemented
	res := notImplemented()
	if res.code != xcode {
		t.Errorf("%s returned code %d; expected %d",
			fname, res.code, xcode)
	}
	if res.data == nil {
		t.Errorf("%s returned nil error; expected non-nil", fname)
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

func genListInterface(form string, N int) []interface{} {
	ret := make([]interface{}, N)
	for i := 0; i < N; i++ {
		ret[i] = fmt.Sprintf(form, i)
	}
	return ret
}

func sqlValues_Checker(t *testing.T, form string, N int) {
	fname := "validateSQLValues"
	values := genListInterface(form, N)
	err := validateSQLValues(values)
	if err != nil {
		t.Errorf("%s(...) failed on length=%d", fname, N)
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
	fname := "validateSQLKeys"
	values := genList(form, N)
	err := validateSQLKeys(values)
	if xsucc != (err == nil) {
		msg := "true"
		if err != nil {
			msg = err.Error()
		}
		t.Errorf(`%s("%s"...)=%s; expected %t`,
			fname, form, msg, xsucc)
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
	fname := "nstring"
	res := nstring(s, n)
	rlist := strings.Split(res, ",")
	if n == 0 {
		// this must be handled as a special case
		// because strings.Split() returns a list of length 1
		// on empty string.
		if res != "" {
			t.Errorf(`%s("%s",%d)="%s"; expected ""`,
				fname, s, n, res)
		}
		return
	} else if n != len(rlist) {
		t.Errorf(`%s("%s",%d)="%s" failed split test`,
			fname, s, n, res)
		return
	}
	for _, v := range rlist {
		if v != s {
			t.Errorf(`%s("%s",%d) bad item "%s"`,
				fname, s, n, v)
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
	fname := "errorRet"
	err := fmt.Errorf("%s", msg)
	res := errorRet(code, err)
	if code != res.code {
		t.Errorf(`#%d: %s returned (%d,); expected %d`,
			i, fname, res.code, code)
		return
	}
	eresp, ok := res.data.(ErrorResponse)
	if !ok {
		t.Errorf(`#%d: %s ErrorResponse conversion error`, i, fname)
		return
	}
	if code != eresp.Code {
		t.Errorf(`#%d: %s ErrorResponse.Code=%d; expected %d`,
			i, fname, eresp.Code, code)
		return
	}
	if msg != eresp.Message {
		t.Errorf(`#%d: %s ErrorResponse.Message="%s"; expected "%s"`,
			i, fname, eresp.Message, msg)
	}
}

func Test_errorRet(t *testing.T) {
	for i, tc := range errorRet_Tab {
		errorRet_Checker(t, i, tc.code, tc.msg)
	}
}

// ----- unit tests for mkIdClause()

func fakeParams(paramstr string) map[string]string {
	ret := map[string]string{}
	if paramstr == "" {
		return ret
	}
	var name, value string
	strlist := strings.Split(paramstr, "&")
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
	slist := strings.Split(s, ",")
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

func mkIdClause_Checker(t *testing.T, i int, tc idclause_TC) {
	fname := "mkIdClause"
	params := fakeParams(tc.paramstr)
	res, idlist, err := mkIdClause(params)
	if tc.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s([%s]) returned status=[%s]; expected [%t]`,
			i, fname, tc.paramstr, msg, tc.xsucc)
		return
	}
	if err != nil {
		return
	}
	if tc.xres != res {
		t.Errorf(`#%d: %s([%s]) returned "%s"; expected "%s"`,
			i, fname, tc.paramstr, res, tc.xres)
	}

	resids, err := idListToA(idlist)
	if err != nil {
		t.Errorf(`#%d: %s idListToA error "%s"`, i, fname, err)
	}
	if tc.xids != resids {
		t.Errorf(`#%d: %s([%s]) idlist=[%s]; expected [%s]`,
			i, fname, tc.paramstr, resids, tc.xids)
	}
}

func Test_mkIdClause(t *testing.T) {
	for i, tc := range idclause_Tab {
		mkIdClause_Checker(t, i, tc)
	}
}

// ----- unit tests for mkIdClauseUpdate()

var mkIdClauseUpdate_Tab = []idclause_TC {
	{ "id_field=id&id=123", "WHERE id = 123", "", true },
	{ "id_field=id&ids=123", "WHERE id in (123)", "", true },
	{ "id_field=id&ids=123,456", "WHERE id in (123,456)", "", true },
	{ "id_field=id", "", "", true },
}

func mkIdClauseUpdate_Checker(t *testing.T, i int, tc idclause_TC) {
	fname := "mkIdClauseUpdate"
	params := fakeParams(tc.paramstr)
	res, err := mkIdClauseUpdate(params)
	if tc.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s([%s]) returned status=[%s]; expected [%t]`,
			i, fname, tc.paramstr, msg, tc.xsucc)
		return
	}
	if err != nil {
		return
	}
	if tc.xres != res {
		t.Errorf(`#%d: %s([%s]) returned "%s"; expected "%s"`,
			i, fname, tc.paramstr, res, tc.xres)
	}
}

func Test_mkIdClauseUpdate(t *testing.T) {
	for i, tc := range mkIdClauseUpdate_Tab {
		mkIdClauseUpdate_Checker(t, i, tc)
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

func idTypesToInterface_Checker(t *testing.T, i int, tc string) {
	fname := "idTypesToInterface"
	alist := strings.Split(tc, ",")
	if tc == "" {
		alist = []string{}
	}
	res := idTypesToInterface(alist)
	str, err := idListToA(res)
	if err != nil {
		t.Errorf(`#%d: %s idListToA error "%s"`, i, fname, err)
	}
	if str != tc {
		t.Errorf(`#%d: %s("%s") = "%s"; expected "%s"`,
			i, fname, tc, str, tc)
		return
	}
}

func Test_idTypesToInterface(t *testing.T) {
	for i, tc := range idTypesToInterface_Tab {
		idTypesToInterface_Checker(t, i, tc)
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
func mkSelectString_Checker(t *testing.T, i int, tc mkSelectString_TC) {
	fname := "mkSelectString"
	params := fakeParams(tc.paramstr)
	// fmt.Printf("in %s_Checker, params=%s\n", fname, params)
	res, idlist, err := mkSelectString(params)
	if tc.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s returned status [%s]; expected [%t]`,
			i, fname, msg, tc.xsucc)
		return
	}
	if err != nil {
		return
	}
	if tc.xres != res {
		t.Errorf(`#%d: %s returned "%s"; expected "%s"`,
			i, fname, res, tc.xres)
		return
	}
	ids, err := idListToA(idlist)
	if err != nil {
		t.Errorf(`#%d: %s idListToA error "%s"`, i, fname, err)
	}
	if tc.xids != ids {
		t.Errorf(`#%d: %s returned ids "%s"; expected "%s"`,
			i, fname, ids, tc.xids)
	}
}

func Test_mkSelectString(t *testing.T) {
	for i, tc := range mkSelectString_Tab {
		mkSelectString_Checker(t, i, tc)
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

func getBodyRecord_Checker(t *testing.T, testno int, tc getBodyRecord_TC) {
	fname := "getBodyRecord"

	rdr := strings.NewReader(tc.data)
	req, _ := http.NewRequest(http.MethodPost, "/xyz", rdr)

	tckeys := strings.Split(tc.keys, "&")
	tcvalues := strings.Split(tc.values, "&")
	nkeys := len(tckeys)

	body, err := getBodyRecord(mkApiHandlerArg(req, nil))
	if err != nil {
		t.Errorf("#%d: %s([%s]) failed, error=%s",
			testno, fname, tc.data, err)
	}
	records := body.Records
	nrecs := len(records)

	if nkeys != nrecs {
		t.Errorf(`#%d: %s returned Records length=%d; expected %d`,
			testno, fname, nrecs, nkeys)
	}
	for j := 0; j < nrecs; j++ {
		rec := records[j]
		keystr := strings.Join(rec.Keys, ",")
		if tckeys[j] != keystr {
			t.Errorf(`#%d %s Record[%d] keys=%s; expected %s`,
				testno, fname, j, keystr, tckeys[j])
		}
		valstr := strings.Join(unmaskStrings(rec.Values), ",")
		if tcvalues[j] != valstr {
			t.Errorf(`#%d %s Record[%d] values=%s; expected %s`,
				testno, fname, j, valstr, tcvalues[j])
		}
	}
}

func Test_getBodyRecord(t *testing.T) {
	for testno, tc := range getBodyRecord_Tab {
		getBodyRecord_Checker(t, testno, tc)
	}
}

// ----- unit tests for convTableNames() and grabNameField()

type convTableNames_TC struct {
	names string
}

var convTableNames_Tab = []convTableNames_TC {
	{""},
	{"a"},
	{"a,b"},
	{"abc,def,ghi"},
}

// mimicTableNamesQuery() returns an object that mimics the return from
// the query to the "tables" table.
func mimicTableNamesQuery(names []string) []*KVRecord {
	N := len(names)
	ret := make([]*KVRecord, N)
	for i := 0; i < N; i++ {
		Keys := []string{"name"}
		Values := []interface{}{names[i]}
		ret[i] = &KVRecord{Keys, Values}
	}
	return ret
}

func convTableNames_Checker(t *testing.T, testno int, tc convTableNames_TC) {
	fname := "convTableNames"
	names := mySplit(tc.names, ",")
	obj := mimicTableNamesQuery(names)
	// fmt.Printf("obj=%s\n", obj)
	res, err := convTableNames(obj)
	if err != nil {
		t.Errorf("#%d: %s([%s]) returned error", testno, fname, tc.names)
		return
	}
	resJoin := strings.Join(res, ",")
	if tc.names != resJoin {
		t.Errorf(`#%d: %s([%s]) = "%s"; expected "%s"`,
			testno, fname, tc.names, resJoin, tc.names)
	}
}

func Test_convTableNames(t *testing.T) {
	for testno, tc := range convTableNames_Tab {
		convTableNames_Checker(t, testno, tc)
	}
}

func Test_convTableNames_bad(t *testing.T) {
	fn := "convTableNames"

	// create a good object, then munge it to force error
	names := []string{"abc", "def"}
	obj := mimicTableNamesQuery(names)
	vals := obj[0].Values
	vals[0] = Test_convTableNames_bad  // junk that can't be converted

	_, err := convTableNames(obj)
	if err == nil {
		t.Errorf("%s call succeeded; expected error", fn)
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

func validateRecords_Checker(t *testing.T, testno int, tc validateRecords_TC) {
	fname := "validateRecords"
	records := mkRecords(tc.desc)
	res := validateRecords(records)
	if tc.xsucc != (res == nil) {
		t.Errorf(`#%d: %s([%s]) = [%s]; expected %t`,
			testno, fname, tc.desc, errRep(nil), tc.xsucc)
	}
}

func Test_validateRecords(t *testing.T) {
	for testno, tc := range validateRecords_Tab {
		validateRecords_Checker(t, testno, tc)
	}
}

// ----- unit tests for convValues()

// inputs and outputs for one convValues testcase.
type convValues_TC struct {
	arg string
}

// table of convValues testcases.
var convValues_Tab = []convValues_TC {
	{ "" },
	{ "abc" },
	{ "abc,def" },
	{ "abc,def,ghi" },
}

func strToSQLValues(arg string) []interface{} {
	args := mySplit(arg, ",")
	N := len(args)
	ret := make([]interface{}, N)
	for i, s := range args {
		rb := sql.RawBytes(s)
		ret[i] = &rb;
	}
	return ret
}

// return something that can't be converted by convValues().
func mkIllegalValues() []interface{} {
	ret := make([]interface{}, 1)
	val := 555
	ret[0] = &val
	return ret
}

// run one testcase for function convValues.
func convValues_Checker(t *testing.T, testno int, tc convValues_TC) {
	fname := "convValues"
	argInter := strToSQLValues(tc.arg)
	err := convValues(argInter)
	if err != nil {
		t.Errorf(`#%d: %s([%s]) failed [%s]`,
			testno, fname, tc.arg, err)
	}
	argStrings := unmaskStrings(argInter)
	resultStr := strings.Join(argStrings, ",")
	if tc.arg != resultStr {
		t.Errorf(`#%d: %s("%s")="%s"; expected "%s"`,
			testno, fname, tc.arg, resultStr, tc.arg)
	}
}

// main test suite for convValues().
func Test_convValues(t *testing.T) {
	for testno, tc := range convValues_Tab {
		convValues_Checker(t, testno, tc)
	}
}

// test suite for testing error return.
func Test_convValues_illegal(t *testing.T) {
	fname := "convValues"
	vals := mkIllegalValues()
	err := convValues(vals)
	if err == nil {
		t.Errorf(`%s on illegal value failed to return error`, fname)
	}
}

// ----- unit tests for getDbResourcesHandler

type apiCall_TC struct {
	verb string
	path string
	query string
	body string
	xcode int
}

var getDbResources_Tab = []apiCall_TC {
	{http.MethodGet, "/db", "", "", http.StatusNotImplemented},
}

func apiCall_Checker(t *testing.T,
		testno int,
		f apiHandler,
		fname string,
		tc apiCall_TC) {
	rdr := strings.NewReader(tc.body)
	url := fmt.Sprintf("%s?%s", tc.path, tc.query)
	req, _ := http.NewRequest(tc.verb, url, rdr)
	arg := mkApiHandlerArg(req, nil)
	res := f(arg)
	if tc.xcode != res.code {
		t.Errorf(`#%d: %s(%s,%s) = %d; expected %d`,
			testno, fname, tc.verb, url,
			res.code, tc.xcode)
	}
}

func apiCalls_Runner(t *testing.T, f apiHandler, tab []apiCall_TC) {
	fname := getFunctionName(f)
	for testno, tc := range tab {
		apiCall_Checker(t, testno, f, fname, tc)
	}
}

func Test_getDbResourcesHandler(t *testing.T) {
	apiCalls_Runner(t, getDbResourcesHandler, getDbResources_Tab)
}

// ----- unit tests for getDbSchemasHandler()

var getDbSchemas_Tab = []apiCall_TC {
	{http.MethodGet, "/db/_schema", "", "", http.StatusNotImplemented},
}

func Test_getDbSchemasHandler(t *testing.T) {
	apiCalls_Runner(t, getDbSchemasHandler, getDbSchemas_Tab)
}

// ----- unit tests for createDbTableHandler()

var createDbTable_Tab = []apiCall_TC {
	{http.MethodPost, "/db/_schema", "", "", http.StatusNotImplemented},
}

func Test_getDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, createDbTableHandler, createDbTable_Tab)
}

// ----- unit tests for updateDbTablesHandler()

var updateDbTables_Tab = []apiCall_TC {
	{http.MethodPatch, "/db/_schema", "", "", http.StatusNotImplemented},
}

func Test_updateDbTablesHandler(t *testing.T) {
	apiCalls_Runner(t, updateDbTablesHandler, updateDbTables_Tab)
}

// ----- unit tests for describeDbTableHandler()

var describeDbTable_Tab = []apiCall_TC {
	{http.MethodGet, "/db/_schema/tabname", "", "", http.StatusNotImplemented},
}

func Test_describeDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, describeDbTableHandler, describeDbTable_Tab)
}

// ----- unit tests for createDbTablesHandler()

var createDbTables_Tab = []apiCall_TC {
	{http.MethodPost, "/db/_schema/tabname", "", "", http.StatusNotImplemented},
}

func Test_createDbTablesHandler(t *testing.T) {
	apiCalls_Runner(t, createDbTablesHandler, createDbTables_Tab)
}

// ----- unit tests for deleteDbTableHandler()

var deleteDbTable_Tab = []apiCall_TC {
	{http.MethodDelete, "/db/_schema/tabname", "", "", http.StatusNotImplemented},
}

func Test_deleteDbTableHandler(t *testing.T) {
	apiCalls_Runner(t, deleteDbTableHandler, deleteDbTable_Tab)
}

// ----- unit tests for describeDbField()

var describeDbField_Tab = []apiCall_TC {
	{http.MethodDelete, "/db/_schema/tabname", "", "", http.StatusNotImplemented},
}

func Test_describeDbFieldHandler(t *testing.T) {
	apiCalls_Runner(t, describeDbFieldHandler, describeDbField_Tab)
}

// ----- unit tests for getDbRecord()

var getDbRecord_Tab = []apiCall_TC {
	/*
	{http.MethodGet, "/db/_table/tabname",
		"table_name=bundles&id=123", "", http.StatusBadRequest},
	{http.MethodGet, "/db/_table/tabname",
		"table_name=bundles&id=1", "", http.StatusOK},
	{http.MethodGet, "/db/_table/tabname",
		"table_name=bundles&id=3", "", http.StatusOK},
	 */
}

/*
func Test_getDbRecordHandler(t *testing.T) {
	var err error
	db, err = fakeInitDb()
	if err != nil {
		t.Errorf(`fakeInitDb failed: [%s]`, err)
	}
	apiCalls_Runner(t, getDbRecordHandler, getDbRecord_Tab)
}
 */

// ----- unit tests for createDbRecords()

var createDbRecords_Tab = []apiCall_TC {
	/*
	{http.MethodPost, "/db/_table/tabname", "table_name=bundles",
		`{"Records":[{"Keys":["name","uri"],"Values":["abc1","xyz1"]}]}`,
		http.StatusBadRequest},
	{http.MethodPost, "/db/_table/tabname", "table_name=bundles",
		`{"Records":[{"Keys":["name","uri"],"Values":["abc2","xyz2"]}]}`,
		http.StatusOK},
	{http.MethodPost, "/db/_table/tabname", "table_name=bundles",
		`{"Records":[{"Keys":["name","uri"],"Values":["abc3","xyz3"]}]}`,
		http.StatusOK},
	{http.MethodPost, "/db/_table/tabname", "table_name=bundles",
		`{"Records":[{"Keys":["name","uri"],"Values":["abc4","xyz4"]}]}`,
		http.StatusOK},
	 */
}

/*
func Test_createDbRecordsHandler(t *testing.T) {
	var err error
	db, err = fakeInitDb()
	if err != nil {
		t.Errorf(`fakeInitDb failed: [%s]`, err)
	}
	apiCalls_Runner(t, createDbRecordsHandler, createDbRecords_Tab)
}
 */
