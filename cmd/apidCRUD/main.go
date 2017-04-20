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
	dprintf("in main...")

	// initialize apid using default services
	dprintf("before Initialize...")
	apid.Initialize(factory.DefaultServicesFactory())

	dprintf("before Log...")
	log := apid.Log()

	// call all initialization functions on all registered plugins
	dprintf("before InitializePlugins...")
	apid.InitializePlugins()

	// start client API listener
	api := apid.API()
	err := api.Listen()

	// if we got here, an error occurred
	config := apid.Config()
	api_listen := config.GetString("api_listen")
	log.Fatalf("Error. Is something already running on %s? %s", api_listen, err)
}

func dprintf(format string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(format, args...)
		fmt.Print()
	}
}
