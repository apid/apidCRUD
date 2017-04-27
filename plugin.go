package apidCRUD

import (
	"net/http"
	"github.com/30x/apid-core"
)

// getStringer is an interface that supports GetString().
type getStringer interface {
	GetString(string) string
}

type handleFuncer interface {
	HandleFunc(path string, hf http.HandlerFunc) apid.Route
}

// apiTable is the list of APIs that need to be wired up.
var apiTable = []apiDesc{
	{ "/db", http.MethodGet, getDbResourcesHandler },
	{ "/db/_table", http.MethodGet, getDbTablesHandler },
	{ "/db/_table/{table_name}", http.MethodGet, getDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodPost, createDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodDelete, deleteDbRecordsHandler },
	{ "/db/_table/{table_name}", http.MethodPatch, updateDbRecordsHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodGet, getDbRecordHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodPatch, updateDbRecordHandler },
	{ "/db/_table/{table_name}/{id}", http.MethodDelete, deleteDbRecordHandler },
	{ "/db/_schema", http.MethodGet, getDbSchemasHandler },
	{ "/db/_schema", http.MethodPost, createDbTablesHandler },
	{ "/db/_schema", http.MethodPatch, updateDbTablesHandler },
	{ "/db/_schema/{table_name}", http.MethodGet, describeDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodPost, createDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodDelete, deleteDbTableHandler },
	{ "/db/_schema/{table_name}/{field_name}", http.MethodGet, describeDbFieldHandler },
}

// initPlugin() is called by the apid-core startup.
func initPlugin(services apid.Services) (apid.PluginData, error) {
	log = services.Log().ForModule(pluginData.Name)
	log.Printf("in initPlugin")

	initConfig()

	var err error
	db, err = initDB()
	if err != nil {
		return pluginData, err
	}

	registerHandlers(services.API(), apiTable)

	return pluginData, nil
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
			pathDispatch(vmap, w, apiHandlerArg{r})
		})
}

// confGet() returns the config value of the named string,
// or if there is no configured value, the given default value.
func confGet(cfg getStringer, vname string, defval string) string {
	ret := cfg.GetString(vname)
	if ret == "" {
		return defval
	}
	return ret
}

// initConfig() sets up some global configuration parameters for this plugin.
func initConfig() {
	cfg := apid.Config()

	dbName := confGet(cfg, "apidCRUD_db_name", "apidCRUD.db")
	log.Debugf("apidCRUD_db_name = %s", dbName)

	base_path := confGet(cfg, "apidCRUD_base_path", "/apid")
	log.Debugf("apidCRUD_base_path = %s", base_path)
}
