package apidCRUD

import (
	"testing"
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

// mockGetStringer is the interface expected by getConf()
type mockGetStringer struct {
	data map[string]string
}

func (self mockGetStringer) GetString(name string) string {
	return self.data[name]
}

var fakeConfData = map[string]string{"there": "yes"}

func confGet_Checker(t *testing.T, i int, gs GetStringer, test *confGet_TC) {
	fn := "confGet"
	res := confGet(gs, test.name, test.defval)
	if test.xval != res {
		t.Errorf(`#%d: %s("%s","%s")="%s"; expected "%s"`,
			i, fn, test.name, test.defval, res, test.xval)
	}
}

func Test_getConf(t *testing.T) {
	gs := mockGetStringer{fakeConfData}
	for i, test := range confGet_Tab {
		confGet_Checker(t, i, gs, &test)
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
