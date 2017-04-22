package apidCRUD

import (
	"fmt"
	"net/http"
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

func NewApiWiring(basePath string, tab []apiDesc) (*apiWiring) {
	pm := make(map[string]verbMap)
	apiws := &apiWiring{pm}
	for _, b := range(tab) {
		apiws.AddApi(basePath + b.path, b.verb, b.handler)
	}
	return apiws
}

func (apiws *apiWiring) AddApi(path string, verb string, handler apiHandler) {
	vmap, ok := apiws.pathsMap[path]
	if !ok {
		vmap = verbMap{path: path, methods: map[string]apiHandler{}}
		apiws.pathsMap[path] = vmap
	}
	vmap.methods[verb] = handler
}

func (apiws *apiWiring) GetMaps() map[string]verbMap {
	return apiws.pathsMap
}

func CallFunc(vmap verbMap, verb string, req *http.Request) (int, interface{}) {
	verbFunc, ok := vmap.methods[verb]
	if !ok {
		return badStat, fmt.Errorf("internal wiring error for %s on %s",
			verb, vmap.path)
	}

	return verbFunc(req)
}
