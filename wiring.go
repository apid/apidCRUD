package apidCRUD

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type apiHandler func(*http.Request) (int, interface{})

// type VerbMap maps each wired verb for a given path, to its handler function.
type VerbMap struct {
	path string
	methods map[string]apiHandler
}

type apiDesc struct {
	path string
	verb string
	handler apiHandler
}

// type ApiWiring is the current state of the API wiring.
type ApiWiring struct {
	pathsMap map[string]VerbMap
}

// NewApiWiring returns an API configuration, after adding the APIs
// from the given table.  basePath is prefixed to the paths
// specified in the table items.
func NewApiWiring(basePath string, tab []apiDesc) (*ApiWiring) {
	pm := make(map[string]VerbMap)
	apiws := &ApiWiring{pm}
	for _, b := range(tab) {
		apiws.AddApi(basePath + b.path, b.verb, b.handler)
	}
	return apiws
}

// AddApi() configures the wiring for one path and verb to their handler.
func (apiws *ApiWiring) AddApi(path string, verb string, handler apiHandler) {
	vmap, ok := apiws.pathsMap[path]
	if !ok {
		vmap = VerbMap{path: path, methods: map[string]apiHandler{}}
		apiws.pathsMap[path] = vmap
	}
	vmap.methods[verb] = handler
}

// GetMaps returns the configured path-to-VerbMap mapping,
// for possible range iteration.
func (apiws *ApiWiring) GetMaps() map[string]VerbMap {
	return apiws.pathsMap
}

// CallApiMethod() calls the handler that was configured for
// the given VerbMap and verb.
func CallApiMethod(vmap VerbMap, verb string, req *http.Request) (int, interface{}) {
	verbFunc, ok := vmap.methods[verb]
	if !ok {
		return badStat, fmt.Errorf("internal wiring error for %s on %s",
			verb, vmap.path)
	}

	return verbFunc(req)
}

// dispatch() is the general handler for all our APIs.
// it is called indirectly thru a closure function that
// supplies the vmap argument.
func dispatch(vmap VerbMap, w http.ResponseWriter, req *http.Request) {
	log.Debugf("in dispatch: method=%s path=%s", req.Method, req.URL.Path)
	defer func() {
		_ = req.Body.Close()
	}()

	code, data := CallApiMethod(vmap, req.Method, req)

	rawdata, err := convData(data)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(rawdata)

	log.Debugf("in dispatch: code=%d", code)
}

func convData(data interface{}) ([]byte, error) {
	switch data := data.(type) {
	case []byte:
		return data, nil
	case string:
		return []byte(data), nil
	default: // json conversion
		return json.Marshal(data)
	}
}

// writeErrorResponse() writes to the ResponseWriter,
// the given error's message, and logs it.
func writeErrorResponse(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := err.Error()
	data, _ := convData(ErrorResponse{code,msg})

        w.WriteHeader(code)
        _, _ = w.Write(data)

        log.Errorf("error handling API request: %s", msg)
}
