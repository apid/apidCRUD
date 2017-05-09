package apidCRUD

import (
	"testing"
)

// global conf variables with values to be used during unit testing.
// used by initConfig().
var utConfData = map[string]string {
	"apidCRUD_base_path": "",
	"apidCRUD_max_recs": "7",
	"apidCRUD_db_driver": "sqlite3",
	"apidCRUD_db_name": "unit-test.db",
}

// ----- unit tests for confGet()

type confGet_TC struct {
	name string
	defval string
	xval string
}

var confGet_Tab = []confGet_TC {
	{"apidCRUD_db_name", "garbage", "unit-test.db"}, // this key is present
	{"not-there", "no", "no"},		// this key is not present
}

// mockGetStringer is compatible with the interface expected by confGet().
type mockGetStringer struct {
	data map[string]string
}

func (gs mockGetStringer) GetString(name string) string {
	return gs.data[name]
}

func confGet_Checker(cx *testContext, gs getStringer, tc *confGet_TC) {
	res := confGet(gs, tc.name, tc.defval)
	cx.assertEqual(tc.xval, res, "result")
}

func Test_confGet(t *testing.T) {
	cx := newTestContext(t, "confGet_Tab")
	gs := mockGetStringer{utConfData}
	for _, tc := range confGet_Tab {
		confGet_Checker(cx, gs, &tc)
		cx.bump()
	}
}

// ----- unit tests for initConfig()

func utInitConfig() {
	// this just proves it can be called without crashing
	gs := mockGetStringer{utConfData}
	initConfig(gs)
}
