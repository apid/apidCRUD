package apidCRUD

import (
	"net/http"
	"github.com/30x/apid-core"
)

// ----- narrowed-down versions of apid service interfaces

// these narrowed interfaces make testing easier,
// by making it easier to hand craft a simple mock interface
// that can be used as an argument to pieces of code under test.

// handleFuncer interface provides the HandleFunc() method.
// narrowed from apid.APISerivce.
type handleFuncer interface {
	HandleFunc(path string, hf http.HandlerFunc) apid.Route
}

// forModuler interface proviees the ForModule() method.
// narrowed from apid.LogService.
type forModuler interface {
	ForModule(name string) apid.LogService
}

// ----- apiTable is the list of APIs that need to be wired up.
var apiTable = []apiDesc {
	{ "/db", http.MethodGet, getDbResourcesHandler },
	{ "/db/_table", http.MethodGet, getDbTablesHandler },
	{ "/db/_table/{table_name}", http.MethodGet, getDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodPost, createDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodDelete, deleteDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodPatch, updateDbRecordsHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodGet, getDbRecordHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodPatch, updateDbRecordHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodDelete, deleteDbRecordHandler },
	{ "/db/_schema/{table_name}", http.MethodPost, createDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodGet, describeDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodDelete, deleteDbTableHandler },
}

// ----- functions go below this line

// initPlugin() is called by the apid InitializePlugins().
// just calls realInitPlugin() which has been designed to simplify unit testing.
func initPlugin(services apid.Services) (apid.PluginData, error) {
	return realInitPlugin(services.Config(), services.Log(), services.API())
}

// realInitPlugin() drives miscellaneous plugin-specific setup activities,
// then returns apidCRUD's pluginData.
//	reads in the plugin-specific configuration data.
//	sets the log variable.
//	sets the db variable.
//	registers the API handlers.
func realInitPlugin(gsi getStringer,
		fmi forModuler,
		hfi handleFuncer) (apid.PluginData, error) {
	initConfig(gsi)
	log = fmi.ForModule(pluginData.Name)	// NOTE: non-local var
	registerHandlers(hfi, apiTable)
	log.Infof("in apidCRUD realInitPlugin")

	var err error
	db, err = initDB(dbName)		// NOTE: non-local var

	return pluginData, err
}

// registerHandlers() register all our handlers with the given service.
func registerHandlers(service handleFuncer, tab []apiDesc) {
	ws := newApiWiring(basePath, tab)
	maps := ws.GetMaps()
	for path, vmap := range maps {
		addPath(service, path, vmap)
	}
}

// addPath() registers the given path with the given service,
// so that it will be handled indirectly by pathDispatch().
// when an API call is made on this path, the vmap argument from
// this context will be suppllied, along with the w and r arguments
// passed in by the service framework.
func addPath(service handleFuncer, path string, vmap verbMap) {
	service.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			pathDispatch(vmap, w, mkApiHandlerArg(r, getPathParams(r)))
		})
}

// getPathParams() returns a map of path params from the given http request.
func getPathParams(req *http.Request) map[string]string {
	return apid.API().Vars(req)
}
