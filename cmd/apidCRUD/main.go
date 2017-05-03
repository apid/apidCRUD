// package main in apidCRUD is a test version of apid
// built with only the apidCRUD plugin.
package main

import (
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	_ "github.com/30x/apidCRUD"
)

// main() here is a stripped-down version of the real apid main.
func main() {
	// initialize apid using default services
	apid.Initialize(factory.DefaultServicesFactory())

	// call all initialization functions on all registered plugins
	apid.InitializePlugins("xxx")

	// start client API listener
	err := apid.API().Listen()

	// if we got here, an error occurred
	config := apid.Config()
	api_listen := config.GetString("api_listen")
	apid.Log().Fatalf("api.Listen() on %s returned [%s]", api_listen, err)
}
