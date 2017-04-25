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

// ----- unit tests for InitWiring, newApiWiring, GetMaps

func Test_newApiWiring(t *testing.T) {
	fn := "newApiWiring"
	ws := newApiWiring("", []apiDesc{})
	if ws == nil {
		t.Errorf("%s failed", fn)
		return
	}
}

func Test_GetMaps(t *testing.T) {
	fn := "GetMaps"
	ws := newApiWiring("", []apiDesc{})
	maps := ws.GetMaps();
	if len(maps) != 0 {
		t.Errorf("%s unexpectedly nonempty", fn)
	}
}

func Test_addApi(t *testing.T) {
	fn := "addApi"
	ws := newApiWiring("", fakeApiTable)
	maps := ws.GetMaps()
	N := countPaths(fakeApiTable)
	wslen := len(maps)
	if N != wslen {
		t.Errorf("%s maps length=%d; expected %d", fn, wslen, N)
	}
}

// ----- unit tests for callApiMethod()

type callApiMethod_TC struct {
	path string
	verb string
	xcode int
}

var callApiMethod_Tab = []callApiMethod_TC {
	{ "/abc", http.MethodGet, abcGetRet },
	{ "/abc", http.MethodPost, abcPostRet },
	{ "/xyz", http.MethodPut, xyzPutRet },
}

func callApiMethod_Checker(t *testing.T, i int, ws *apiWiring, tc callApiMethod_TC) {
	fn := "callApiMethod"
	vmap, ok := ws.pathsMap[tc.path]
	if !ok {
		t.Errorf(`#%d: %s bad path "%s"`, i, fn, tc.path)
		return
	}
	code, _ := callApiMethod(vmap, tc.verb, nil)
	if tc.xcode != code {
		t.Errorf(`#%d: %s("%s","%s")=%d; expected %d`,
			i, fn, tc.path, tc.verb, code, tc.xcode)
	}
}

func Test_callApiMethod(t *testing.T) {
	ws := newApiWiring("", fakeApiTable)
	for i, tc := range callApiMethod_Tab {
		callApiMethod_Checker(t, i, ws, tc)
	}
}

// ----- unit tests for pathDispatch()

func pathDispatch_Checker(t *testing.T, i int, ws *apiWiring, tc callApiMethod_TC) {
	fn := "pathDispatch"

	vmap, ok := ws.pathsMap[tc.path]
	if !ok {
		t.Errorf(`#%d: %s bad path "%s"`, i, fn, tc.path)
		return
	}

	rdr := strings.NewReader("")
	req, _ := http.NewRequest(tc.verb, tc.path, rdr)
	w := httptest.NewRecorder()

	pathDispatch(vmap, w, req)
	code := w.Code
	if tc.xcode != code {
		t.Errorf(`#%d: %s("%s","%s") code=%d; expected %d`,
			i, fn, tc.path, tc.verb, code, tc.xcode)
	}
}

func Test_pathDispatch(t *testing.T) {
	ws := newApiWiring("", fakeApiTable)
	for i, tc := range callApiMethod_Tab {
		pathDispatch_Checker(t, i, ws, tc)
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

func convData_Checker(t *testing.T, i int, tc convData_TC) {
	fn := "convData"
	res, err := convData(tc.idata)
	if tc.xsucc != (err == nil) {
		msg := errRep(err)
		t.Errorf(`#%d: %s returned status=[%s]; expected %t`,
			i, fn, msg, tc.xsucc)
	}
	if err != nil {
		// if the actual call failed, nothing more can be checked.
		return
	}
	if ! reflect.DeepEqual(tc.xbytes, res) {
		t.Errorf(`#%d: %s returned data=[%s]; expected [%s]`,
			i, fn, res, tc.xbytes)
	}
}

func Test_convData(t *testing.T) {
	for i, tc := range convData_Tab {
		convData_Checker(t, i, tc)
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

func writeErrorResponse_Checker(t *testing.T, i int, tc writeErrorResponse_TC) {
	fn := "writeErrorResponse"
	w := httptest.NewRecorder()
	err := fmt.Errorf("%s", tc.msg)
	writeErrorResponse(w, err)
	if tc.xcode != w.Code {
		t.Errorf(`#%d: %s wrote code=%d; expected %d`,
			i, fn, w.Code, tc.xcode)
		return
	}
	body := w.Body.Bytes()
	erec := &ErrorResponse{}
	_ = json.Unmarshal(body, erec)
	if tc.xcode != erec.Code {
		t.Errorf(`#%d: %s ErrorResponse code=%d; expected %d`,
			i, fn, erec.Code, tc.xcode)
		return
	}
	if tc.msg != erec.Message {
		t.Errorf(`#%d: %s ErrorResponse msg="%s"; expected "%s"`,
			i, fn, erec.Message, tc.msg)
	}
}

func Test_writeErrorResponse(t *testing.T) {
	for i, tc := range writeErrorResponse_Tab {
		writeErrorResponse_Checker(t, i, tc)
	}
}
