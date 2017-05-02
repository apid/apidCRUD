package apidCRUD

import (
	"testing"
	"strings"
	"net/http"
	"net/http/httptest"
	"github.com/30x/apid-core"
	)

// ----- unit tests for confGet()

type confGet_TC struct {
	name string
	defval string
	xval string
}

var confGet_Tab = []confGet_TC {
	{"there", "no", "yes"},		// this key is present in
	{"not-there", "no", "no"},	// this key is not present
}

// mockGetStringer is compatible with the interface expected by confGet().
type mockGetStringer struct {
	data map[string]string
}

func (gs mockGetStringer) GetString(name string) string {
	return gs.data[name]
}

var fakeConfData = map[string]string{"there": "yes"}

func confGet_Checker(t *testing.T, i int, gs getStringer, tc *confGet_TC) {
	fn := "confGet"
	res := confGet(gs, tc.name, tc.defval)
	if tc.xval != res {
		t.Errorf(`#%d: %s("%s","%s")="%s"; expected "%s"`,
			i, fn, tc.name, tc.defval, res, tc.xval)
	}
}

func Test_confGet(t *testing.T) {
	gs := mockGetStringer{fakeConfData}
	for i, tc := range confGet_Tab {
		confGet_Checker(t, i, gs, &tc)
	}
}

// ----- unit tests for initDB()

func Test_initDB(t *testing.T) {
	x, err := initDB(dbName)
	if err != nil {
		t.Errorf(`initDB() error %s`, err)
		return
	}
	if x.handle == nil {
		t.Errorf(`initDB() returned nil handle`)
	}
}

// ----- unit tests for initConfig()

func Test_initConfig(t *testing.T) {
	// this just proves it can be called without crashing
	gs := mockGetStringer{fakeConfData}
	initConfig(gs)
}

// ----- unit tests for registerHandlers() and addPath()

type mockApiService struct {
	hfmap map[string]http.HandlerFunc
}

func newMockApiService() *mockApiService {
	fmap := make(map[string]http.HandlerFunc)
	return &mockApiService{fmap}
}

func (service mockApiService) HandleFunc(path string,
				hf http.HandlerFunc) apid.Route {
	// record the handler that is being registered.
	service.hfmap[path] = hf
	return nil
}

func registerHandler_Checker(t *testing.T,
		testno int,
		service *mockApiService,
		tc callApiMethod_TC) {
	fn := "registerHandlers"
	path := basePath + tc.descStr
	fp := service.hfmap[path]
	if fp == nil {
		t.Errorf("%s handler for %s is nil", fn, path)
		return
	}

	r, _ := http.NewRequest(tc.verb, path, strings.NewReader(""))
	w := httptest.NewRecorder()

	// make the call
	fp(w, r)

	// check the recorded response
	if tc.xcode != w.Code {
		t.Errorf(`#%d: %s returned code=%d; expected %d`,
			testno, fn, w.Code, tc.xcode)
		return
	}
}

func Test_registerHandlers(t *testing.T) {
	service := newMockApiService()
	registerHandlers(service, fakeApiTable)

	// check that the expected paths were in fact registered.
	for testno, desc := range callApiMethod_Tab {
		registerHandler_Checker(t, testno, service, desc)
	}
}

// ----- unit tests for initPlugin()

type mockForModuler struct {
	name string
}

func (fmi mockForModuler) ForModule(name string) apid.LogService {
	fmi.name = name
	return apid.Log()
}

func Test_realInitPlugin(t *testing.T) {
	fn := "realInitPlugin"
	gsi := mockGetStringer{}
	fmi := mockForModuler{}
	hfi := newMockApiService()
	_, err := realInitPlugin(gsi, fmi, *hfi)
	if err != nil {
		t.Errorf(`%s returned error [%s]`, fn, err)
	}
}
