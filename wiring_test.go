package apidCRUD

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"reflect"
	"fmt"
	"encoding/json"
)

const (
	abcGetRet = 123
	abcPostRet = 456
	xyzPutRet = 789
)

// ----- unit tests for initWiring()

// a dummy handler, returns abcGetRet.
func abcGetHandler(req *http.Request) (int, interface{}) {
	return abcGetRet, ""
}

// a dummy handler, returns abcPostRet.
func abcPostHandler(req *http.Request) (int, interface{}) {
	return abcPostRet, ""
}

// a dummy handler, returns xyzPutRet.
func xyzPutHandler(req *http.Request) (int, interface{}) {
	return xyzPutRet, ""
}

var fakeApiTable = []apiDesc {
	{ "/abc", http.MethodGet, abcGetHandler },
	{ "/abc", http.MethodPost, abcPostHandler },
	{ "/xyz", http.MethodPut, xyzPutHandler },
}

// countPaths returns the number of unique paths in the given tab.
func countPaths(tab []apiDesc) int {
	paths := map[string]int{}
	N := len(tab)
	for i := 0; i < N; i++ {
		paths[tab[i].path] = 1
	}
	return len(paths)
}

// ----- unit tests for InitWiring, NewApiWiring, GetMaps

func Test_NewApiWiring(t *testing.T) {
	fn := "NewApiWiring"
	ws := NewApiWiring("", []apiDesc{})
	if ws == nil {
		t.Errorf("%s failed", fn)
		return
	}
}

func Test_GetMaps(t *testing.T) {
	fn := "GetMaps"
	ws := NewApiWiring("", []apiDesc{})
	maps := ws.GetMaps();
	if len(maps) != 0 {
		t.Errorf("%s unexpectedly nonempty", fn)
	}
}

func Test_addApi(t *testing.T) {
	fn := "addApi"
	ws := NewApiWiring("", fakeApiTable)
	maps := ws.GetMaps()
	N := countPaths(fakeApiTable)
	wslen := len(maps)
	if N != wslen {
		t.Errorf("%s maps length=%d; expected %d", fn, wslen, N)
	}
}

// ----- unit tests for CallApiMethod()

type CallApiMethod_TC struct {
	path string
	verb string
	xcode int
}

var CallApiMethod_Tab = []CallApiMethod_TC {
	{ "/abc", http.MethodGet, abcGetRet },
	{ "/abc", http.MethodPost, abcPostRet },
	{ "/xyz", http.MethodPut, xyzPutRet },
}

func CallApiMethod_Checker(t *testing.T, i int, ws *ApiWiring, test CallApiMethod_TC) {
	fn := "CallApiMethod"
	vmap, ok := ws.pathsMap[test.path]
	if !ok {
		t.Errorf(`#%d: %s bad path "%s"`, i, fn, test.path)
		return
	}
	code, _ := CallApiMethod(vmap, test.verb, nil)
	if test.xcode != code {
		t.Errorf(`#%d: %s("%s","%s")=%d; expected %d`,
			i, fn, test.path, test.verb, code, test.xcode)
	}
}

func Test_CallApiMethod(t *testing.T) {
	ws := NewApiWiring("", fakeApiTable)
	for i, test := range CallApiMethod_Tab {
		CallApiMethod_Checker(t, i, ws, test)
	}
}

// ----- unit tests for dispatch()

func dispatch_Checker(t *testing.T, i int, ws *ApiWiring, test CallApiMethod_TC) {
	fn := "dispatch"

	vmap, ok := ws.pathsMap[test.path]
	if !ok {
		t.Errorf(`#%d: %s bad path "%s"`, i, fn, test.path)
		return
	}

	rdr := strings.NewReader("")
	req, _ := http.NewRequest(test.verb, test.path, rdr)
	w := httptest.NewRecorder()

	dispatch(vmap, w, req)
	code := w.Code
	if test.xcode != code {
		t.Errorf(`#%d: %s("%s","%s") code=%d; expected %d`,
			i, fn, test.path, test.verb, code, test.xcode)
	}
}

func Test_dispatch(t *testing.T) {
	ws := NewApiWiring("", fakeApiTable)
	for i, test := range CallApiMethod_Tab {
		dispatch_Checker(t, i, ws, test)
	}
}

// ----- unit tests for convData()

func errRep(err error) string {
	if err == nil {
		return "true"
	}
	return err.Error()
}

type convData_TC struct {
	idata interface{}
	xbytes []byte
	xsucc bool
}

var erdata = ErrorResponse{567, "junk"}

var erjson = `{"Code":567,"Message":"junk"}`

var convData_Tab = []convData_TC {
	{"abc", []byte("abc"), true},
	{[]byte("xyz"), []byte("xyz"), true},
	{erdata, []byte(erjson), true},
}

func convData_Checker(t *testing.T, i int, test convData_TC) {
	fn := "convData"
	res, err := convData(test.idata)
	if test.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s returned status=[%s]; expected %t`,
			i, fn, msg, test.xsucc)
	}
	if err != nil {
		// if the actual call failed, nothing more can be checked.
		return
	}
	if ! reflect.DeepEqual(test.xbytes, res) {
		t.Errorf(`#%d: %s returned data=[%s]; expected [%s]`,
			i, fn, res, test.xbytes)
	}
}

func Test_convData(t *testing.T) {
	for i, test := range convData_Tab {
		convData_Checker(t, i, test)
	}
}

// ----- unit tests for writeErrorResponse()

type writeErrorResponse_TC struct {
	msg string
	xcode int
}

var writeErrorResponse_Tab = []writeErrorResponse_TC {
	{ "wxyz", http.StatusInternalServerError },
	{ "abcd", http.StatusInternalServerError },
}

func writeErrorResponse_Checker(t *testing.T, i int, test writeErrorResponse_TC) {
	fn := "writeErrorResponse"
	w := httptest.NewRecorder()
	err := fmt.Errorf("%s", test.msg)
	writeErrorResponse(w, err)
	if test.xcode != w.Code {
		t.Errorf(`#%d: %s wrote code=%d; expected %d`,
			i, fn, w.Code, test.xcode)
		return
	}
	body := w.Body.Bytes()
	erec := &ErrorResponse{}
	json.Unmarshal(body, erec)
	if test.xcode != erec.Code {
		t.Errorf(`#%d: %s ErrorResponse code=%d; expected %d`,
			i, fn, erec.Code, test.xcode)
		return
	}
	if test.msg != erec.Message {
		t.Errorf(`#%d: %s ErrorResponse msg="%s"; expected "%s"`,
			i, fn, erec.Message, test.msg)
	}
}

func Test_writeErrorResponse(t *testing.T) {
	for i, test := range writeErrorResponse_Tab {
		writeErrorResponse_Checker(t, i, test)
	}
}
