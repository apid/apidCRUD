package apidCRUD

import (
	"testing"
	"fmt"
	"strings"
	"os"
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

// ----- unit tests for mkIdClause()

func fakeParams(paramstr string) map[string]string {
	ret := map[string]string{}
	if paramstr == "" {
		return ret
	}
	strlist := strings.Split(paramstr, "&")
	var name, value string
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

func aToIdList(s string) []interface{} {
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

func idListToA(idlist []interface{}) string {
	alist := make([]string, len(idlist))
	for i, ival := range idlist {
		val := ival.(int64)
		alist[i] = idTypeToA(val)
	}
	return strings.Join(alist, ",")
}

func mkIdClause_Checker(t *testing.T, i int, test idclause_TC) {
	fn := "mkIdClause"
	params := fakeParams(test.paramstr)
	res, idlist, err := mkIdClause(params)
	if test.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s([%s]) returned status=[%s]; expected [%t]`,
			i, fn, test.paramstr, msg, test.xsucc)
		return
	}
	if err != nil {
		return
	}
	if test.xres != res {
		t.Errorf(`#%d: %s([%s]) returned "%s"; expected "%s"`,
			i, fn, test.paramstr, res, test.xres)
	}

	resids := idListToA(idlist)
	if test.xids != resids {
		t.Errorf(`#%d: %s([%s]) idlist=[%s]; expected [%s]`,
			i, fn, test.paramstr, resids, test.xids)
	}
}

func Test_mkIdClause(t *testing.T) {
	for i, test := range idclause_Tab {
		mkIdClause_Checker(t, i, test)
	}
}

// ----- unit tests for mkIdClauseUpdate()

var mkIdClauseUpdate_Tab = []idclause_TC {
	{ "id_field=id&id=123", "WHERE id = 123", "", true },
	{ "id_field=id&ids=123", "WHERE id in (123)", "", true },
	{ "id_field=id&ids=123,456", "WHERE id in (123,456)", "", true },
	{ "id_field=id", "", "", true },
}

func mkIdClauseUpdate_Checker(t *testing.T, i int, test idclause_TC) {
	fn := "mkIdClauseUpdate"
	params := fakeParams(test.paramstr)
	res, err := mkIdClauseUpdate(params)
	if test.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s([%s]) returned status=[%s]; expected [%t]`,
			i, fn, test.paramstr, msg, test.xsucc)
		return
	}
	if err != nil {
		return
	}
	if test.xres != res {
		t.Errorf(`#%d: %s([%s]) returned "%s"; expected "%s"`,
			i, fn, test.paramstr, res, test.xres)
	}
}

func Test_mkIdClauseUpdate(t *testing.T) {
	for i, test := range mkIdClauseUpdate_Tab {
		mkIdClauseUpdate_Checker(t, i, test)
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

func idTypesToInterface_Checker(t *testing.T, i int, test string) {
	fn := "idTypesToInterface"
	alist := strings.Split(test, ",")
	if test == "" {
		alist = []string{}
	}
	res := idTypesToInterface(alist)
	str := idListToA(res)
	if str != test {
		t.Errorf(`#%d: %s("%s") = "%s"; expected "%s"`,
			i, fn, test, str, test)
		return
	}
}

func Test_idTypesToInterface(t *testing.T) {
	for i, test := range idTypesToInterface_Tab {
		idTypesToInterface_Checker(t, i, test)
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
		"SELECT a FROM T WHERE id = ? LIMIT 1 OFFSET 0;",
		"456", true},
	{"table_name=T&id_field=id&id=123,456&fields=a,b,c&limit=1&offset=0",
		"SELECT a,b,c FROM T WHERE id = ? LIMIT 1 OFFSET 0;",
		"123,456", true},
}

// run one test case
func mkSelectString_Checker(t *testing.T, i int, test mkSelectString_TC) {
	fn := "mkSelectString"
	params := fakeParams(test.paramstr)
	res, idlist, err := mkSelectString(params)
	if test.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s returned status [%s]; expected [%t]`,
			i, fn, msg, test.xsucc)
		return
	}
	if err != nil {
		return
	}
	if test.xres != res {
		t.Errorf(`#%d: %s returned "%s"; expected "%s"`,
			i, fn, res, test.xres)
		return
	}
	ids := idListToA(idlist)
	if test.xids != ids {
		t.Errorf(`#%d: %s returned ids "%s"; expected "%s"`,
			i, fn, ids, test.xids)
	}
}

func Test_mkSelectString(t *testing.T) {
	for i, test := range mkSelectString_Tab {
		mkSelectString_Checker(t, i, test)
	}
}
