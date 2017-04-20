// package main in apidCRUD is a test version of apid
// built with only the apidCRUD plugin.
package main

import (
	"fmt"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	_ "github.com/30x/apidCRUD"
)

const DEBUG = false

func main() {
	dprintf("in main...\n")

	// initialize apid using default services
	dprintf("before Initialize...\n")
	apid.Initialize(factory.DefaultServicesFactory())

	dprintf("before Log...\n")
	log := apid.Log()

	// call all initialization functions on all registered plugins
	dprintf("before InitializePlugins...\n")
	apid.InitializePlugins()

	// start client API listener
	api := apid.API()
	err := api.Listen()

	// if we got here, an error occurred
	config := apid.Config()
	port := config.GetString("api_port")
	log.Fatalf("Error. Is something already running on port %d? %s", port, err)
}

func dprintf(format string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(format, args...)
	}
}
