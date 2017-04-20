package apidCRUD

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

type apiHandler func(*http.Request) (int, interface{})

type verbMap struct {
	path string
	methods map[string]apiHandler
}

type apiDesc struct {
	path string
	verb string
	handler apiHandler
}

type apiWiring struct {
	pathsMap map[string]verbMap
}

func newApiWiring(n int) (*apiWiring) {
	pm := make(map[string]verbMap)
	return &apiWiring{pm}
}

var basePath = "/apid"

var descTable = []apiDesc{
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
	{ "/db/_schema", http.MethodPut, replaceDbTablesHandler },
	{ "/db/_schema", http.MethodPatch, updateDbTablesHandler },
	{ "/db/_schema/{table_name}", http.MethodGet, describeDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodPost, createDbTableHandler },
	{ "/db/_schema/{table_name}", http.MethodDelete, deleteDbTableHandler },
	{ "/db/_schema/{table_name}/{field_name}", http.MethodGet, describeDbFieldHandler },
}

func initWiring() (*apiWiring) {
	apiws := newApiWiring(len(descTable))
	for _, b := range(descTable) {
		apiws.addWiring(basePath + b.path, b.verb, b.handler)
	}
	return apiws
}

func (apiws *apiWiring) addWiring(path string, verb string, handler apiHandler) {
	vmap, ok := apiws.pathsMap[path]
	if !ok {
		vmap = verbMap{path: path, methods: map[string]apiHandler{}}
		apiws.pathsMap[path] = vmap
	}
	vmap.methods[verb] = handler
}

func getFunc(vmap verbMap, verb string) (apiHandler, error) {
	verbFunc, ok := vmap.methods[verb]
	if !ok {
		return nil, fmt.Errorf("internal wiring error for %s on %s",
			verb, vmap.path)
	}
	return verbFunc, nil
}

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
