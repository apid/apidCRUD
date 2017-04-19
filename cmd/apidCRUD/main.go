package main

import (
	"fmt"
	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	_ "github.com/30x/apidCRUD"
)

func main() {
	fmt.Printf("in main...\n")

	fmt.Printf("before Initialize...\n")

	// initialize apid using default services
	apid.Initialize(factory.DefaultServicesFactory())

	fmt.Printf("before Log...\n")
	log := apid.Log()

	fmt.Printf("before InitializePlugins...\n")
	// this will call all initialization functions on all registered plugins
	apid.InitializePlugins()

	fmt.Printf("before Config...\n")
	// print the base url to the console
	config := apid.Config()
	basePath := ""
	port := config.GetString("api_port")
	log.Print()
	log.Printf("apidCRUD API is at: http://localhost:%s%s", port, basePath)
	log.Print()

	// start client API listener
	api := apid.API()
	err := api.Listen() // doesn't return if no error
	log.Fatalf("Error. Is something already running on port %d? %s", port, err)
}
