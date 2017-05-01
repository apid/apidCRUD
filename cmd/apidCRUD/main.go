// package main in apidCRUD is a test version of apid
// built with only the apidCRUD plugin.
package main

import (
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	_ "github.com/30x/apidCRUD"
)

const DEBUG = false

func main() {
	// initialize apid using default services
	apid.Initialize(factory.DefaultServicesFactory())

	log := apid.Log()

	// call all initialization functions on all registered plugins
	apid.InitializePlugins()

	// start client API listener
	api := apid.API()
	err := api.Listen()

	// if we got here, an error occurred
	config := apid.Config()
	api_listen := config.GetString("api_listen")
	log.Fatalf("api.Listen() on %s returned [%s]", api_listen, err)
}
