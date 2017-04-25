package apidCRUD

import (
	"testing"
	"net/http"
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

// mockGetStringer is the interface expected by confGet()
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
	x, err := initDB()
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
	initConfig()
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

func Test_registerHandlers(t *testing.T) {
	fn := "registerHandlers"
	service := newMockApiService()
	registerHandlers(service, apiTable)

	// check that the expected paths were in fact registered.
	for _, desc := range apiTable {
		path := basePath + desc.path
		if service.hfmap[path] == nil {
			t.Errorf("%s handler for %s is nil", fn, path)
		}
	}
}
