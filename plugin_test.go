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

func confGet_Checker(cx *testContext, gs getStringer, tc *confGet_TC) {
	res := confGet(gs, tc.name, tc.defval)
	cx.assertEqual(tc.xval, res, "result")
}

func Test_confGet(t *testing.T) {
	cx := newTestContext(t, "confGet_Tab")
	gs := mockGetStringer{fakeConfData}
	for _, tc := range confGet_Tab {
		confGet_Checker(cx, gs, &tc)
		cx.bump()
	}
}

// ----- unit tests for initDB()

func Test_initDB(t *testing.T) {
	cx := newTestContext(t, "initDB")
	x, err := initDB(dbName)
	if !cx.assertErrorNil(err, "error ret") {
		return
	}
	cx.assertTrue(x.handle != nil, "handle should not be nil")
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

func registerHandler_Checker(cx *testContext,
		service *mockApiService,
		tc callApiMethod_TC) {
	path := basePath + tc.descStr
	fp := service.hfmap[path]
	if !cx.assertTrue(fp != nil, "handler should not be nil") {
		return
	}

	r, _ := http.NewRequest(tc.verb, path, strings.NewReader(""))
	w := httptest.NewRecorder()

	// make the call
	fp(w, r)

	// check the recorded response
	cx.assertEqual(tc.xcode, w.Code, "w.Code")
}

func Test_registerHandlers(t *testing.T) {
	service := newMockApiService()
	registerHandlers(service, fakeApiTable)

	cx := newTestContext(t, "callApiMethod_Tab")
	// check that the expected paths were in fact registered.
	for _, desc := range callApiMethod_Tab {
		registerHandler_Checker(cx, service, desc)
		cx.bump()
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
	cx := newTestContext(t, "realInitPlugin")
	gsi := mockGetStringer{}
	fmi := mockForModuler{}
	hfi := newMockApiService()
	_, err := realInitPlugin(gsi, fmi, *hfi)
	cx.assertErrorNil(err, "returned error")
}
