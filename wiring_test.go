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
	xyzPatchRet = http.StatusMethodNotAllowed
)

// ----- unit tests for initWiring()

// a dummy handler, returns abcGetRet.
func abcGetHandler(harg *apiHandlerArg) apiHandlerRet {
	return apiHandlerRet{abcGetRet, ""}
}

// a dummy handler, returns abcPostRet.
func abcPostHandler(harg *apiHandlerArg) apiHandlerRet {
	return apiHandlerRet{abcPostRet, ""}
}

// a dummy handler, returns xyzPutRet.
func xyzPutHandler(harg *apiHandlerArg) apiHandlerRet {
	return apiHandlerRet{xyzPutRet, ""}
}

// a dummy handler, returns a value that causes convData() to fail
func badHandler(harg *apiHandlerArg) apiHandlerRet {
	return apiHandlerRet{http.StatusInternalServerError, badconv}
}

var fakeApiTable = []apiDesc {	// nolint
	{ "/abc", http.MethodGet, abcGetHandler },
	{ "/abc", http.MethodPost, abcPostHandler },
	{ "/xyz", http.MethodPut, xyzPutHandler },
	{ "/xyz", http.MethodDelete, badHandler },
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
	cx := newTestContext(t, "newApiWiring")
	ws := newApiWiring("", []apiDesc{})
	if ws == nil {
		cx.Errorf("failed")
		return
	}
}

func Test_GetMaps(t *testing.T) {
	cx := newTestContext(t, "GetMaps")
	ws := newApiWiring("", []apiDesc{})
	maps := ws.GetMaps();
	if len(maps) != 0 {
		cx.Errorf("unexpectedly nonempty")
	}
}

func Test_addApi(t *testing.T) {
	cx := newTestContext(t, "addApi")
	ws := newApiWiring("", fakeApiTable)
	maps := ws.GetMaps()
	N := countPaths(fakeApiTable)
	wslen := len(maps)
	if N != wslen {
		cx.Errorf("maps length=%d; expected %d", wslen, N)
	}
}

// ----- unit tests for callApiMethod()

type callApiMethod_TC struct {
	descStr string
	verb string
	xcode int
}

// test cases for callApiMethod
var callApiMethod_Tab = []callApiMethod_TC {
	{ "/abc", http.MethodGet, abcGetRet },
	{ "/abc", http.MethodPost, abcPostRet },
	{ "/xyz", http.MethodPut, xyzPutRet },
	{ "/xyz", http.MethodPatch, xyzPatchRet },
	{ "/xyz", http.MethodDelete, http.StatusInternalServerError },
}

func callApiMethod_Checker(cx *testContext, ws *apiWiring, tc callApiMethod_TC) {
	vmap, ok := ws.pathsMap[tc.descStr]
	if !ok {
		cx.Errorf(`bad path "%s"`, tc.descStr)
		return
	}
	res := callApiMethod(vmap, tc.verb, parseHandlerArg(tc.verb, tc.descStr))
	if tc.xcode != res.code {
		cx.Errorf(`("%s","%s")=%d; expected %d`,
			tc.verb, tc.descStr, res.code, tc.xcode)
	}
}

func Test_callApiMethod(t *testing.T) {
	cx := newTestContext(t, "callApiMethod_Tab", "callApiMethod")
	ws := newApiWiring("", fakeApiTable)
	for _, tc := range callApiMethod_Tab {
		callApiMethod_Checker(cx, ws, tc)
		cx.bump()
	}
}

// ----- unit tests for pathDispatch()

func pathDispatch_Checker(cx *testContext, ws *apiWiring, tc callApiMethod_TC) {
	vmap, ok := ws.pathsMap[tc.descStr]
	if !ok {
		cx.Errorf(`bad path "%s"`, tc.descStr)
		return
	}

	rdr := strings.NewReader("")
	req, _ := http.NewRequest(tc.verb, tc.descStr, rdr)
	w := httptest.NewRecorder()

	pathDispatch(vmap, w, mkApiHandlerArg(req, nil))
	code := w.Code
	if tc.xcode != code {
		cx.Errorf(`("%s","%s") code=%d; expected %d`,
			tc.verb, tc.descStr, code, tc.xcode)
	}
}

func Test_pathDispatch(t *testing.T) {
	cx := newTestContext(t, "pathDispatch_Tab", "pathDispatch")
	ws := newApiWiring("", fakeApiTable)
	for _, tc := range callApiMethod_Tab {
		pathDispatch_Checker(cx, ws, tc)
		cx.bump()
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

var badconv = func() { }	// cause convData to choke.

var convData_Tab = []convData_TC {
	{"abc", []byte("abc"), true},
	{[]byte("xyz"), []byte("xyz"), true},
	{erdata, []byte(erjson), true},
	{badconv, []byte(""), false},
}

func convData_Checker(cx *testContext, tc convData_TC) {
	res, err := convData(tc.idata)
	if tc.xsucc != (err == nil) {
		msg := errRep(err)
		cx.Errorf(`returned status=[%s]; expected %t`,
			msg, tc.xsucc)
	}
	if err != nil {
		// if the actual call failed, nothing more can be checked.
		return
	}
	if ! reflect.DeepEqual(tc.xbytes, res) {
		cx.Errorf(`returned data=[%s]; expected [%s]`,
			res, tc.xbytes)
	}
}

func Test_convData(t *testing.T) {
	cx := newTestContext(t, "convData_Tab", "convData")
	for _, tc := range convData_Tab {
		convData_Checker(cx, tc)
		cx.bump()
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

func writeErrorResponse_Checker(cx *testContext, tc writeErrorResponse_TC) {
	w := httptest.NewRecorder()
	err := fmt.Errorf("%s", tc.msg)
	writeErrorResponse(w, err)
	if tc.xcode != w.Code {
		cx.Errorf(`wrote code=%d; expected %d`,
			w.Code, tc.xcode)
		return
	}
	body := w.Body.Bytes()
	erec := &ErrorResponse{}
	_ = json.Unmarshal(body, erec)
	if tc.xcode != erec.Code {
		cx.Errorf(`ErrorResponse code=%d; expected %d`,
			erec.Code, tc.xcode)
		return
	}
	if tc.msg != erec.Message {
		cx.Errorf(`ErrorResponse msg="%s"; expected "%s"`,
			erec.Message, tc.msg)
	}
}

func Test_writeErrorResponse(t *testing.T) {
	cx := newTestContext(t, "writeErrorResponse_Tab", "writeErrorResponse")
	for _, tc := range writeErrorResponse_Tab {
		writeErrorResponse_Checker(cx, tc)
		cx.bump()
	}
}
