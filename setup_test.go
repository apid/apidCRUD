package apidCRUD

// this module does global setup for unit tests.

import (
	"testing"
	"os"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
)

var testServices = factory.DefaultServicesFactory()

// TestMain() is called by the test framework before running the tests.
// we use it to initialize the log variable.
func TestMain(m *testing.M) {
	// do this in case functions under test need to log something
	apid.Initialize(testServices)
	log = apid.Log()

	var err error
	db, err = fakeInitDB()
	if err != nil {
		panic(err.Error())
	}

	// db, _ = initDB("test.db");

	// required boilerplate
	os.Exit(m.Run())
}