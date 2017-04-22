package apidCRUD

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
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

// ----- unit tests for CallFunc()

type CallFunc_TC struct {
	path string
	verb string
	xcode int
}

var CallFunc_Tab = []CallFunc_TC {
	{ "/abc", http.MethodGet, abcGetRet },
	{ "/abc", http.MethodPost, abcPostRet },
	{ "/xyz", http.MethodPut, xyzPutRet },
}

func CallFunc_Checker(t *testing.T, i int, ws *apiWiring, test CallFunc_TC) {
	fn := "CallFunc"
	vmap, ok := ws.pathsMap[test.path]
	if !ok {
		t.Errorf(`#%d: %s bad path "%s"`, i, fn, test.path)
		return
	}
	code, _ := CallFunc(vmap, test.verb, nil)
	if test.xcode != code {
		t.Errorf(`#%d: %s("%s","%s")=%d; expected %d`,
			i, fn, test.path, test.verb, code, test.xcode)
	}
}

func Test_CallFunc(t *testing.T) {
	ws := NewApiWiring("", fakeApiTable)
	for i, test := range CallFunc_Tab {
		CallFunc_Checker(t, i, ws, test)
	}
}

// ----- unit tests for dispatch()

func dispatch_Checker(t *testing.T, i int, ws *apiWiring, test CallFunc_TC) {
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
	for i, test := range CallFunc_Tab {
		dispatch_Checker(t, i, ws, test)
	}
}
