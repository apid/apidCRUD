package apidCRUD

// this module defines the data structures and functions
// that translate from apid's model of a handler function
// (where a single function deals with all verbs on a given path),
// with a higher-level concept of an api handler (where one
// function deals with a single verb on a given path).

import (
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"io"
	"sort"
)

// apiHandlerRet is the return type from an apiHandler function.
type apiHandlerRet struct {
	code int
	data interface{}
}

// apiHandlerArg is the type of the parameter to an apiHandler function.
type apiHandlerArg struct {
	req *http.Request
	pathParams map[string]string
	err error
}

// apiHandler is the type an API handler function.
type apiHandler func(*apiHandlerArg) apiHandlerRet

// type verbMap maps each wired verb for a given path, to its handler function.
type verbMap struct {
	path string
	methods map[string]apiHandler
}

// type apiDesc describes the wiring for one API.
type apiDesc struct {
	path string
	verb string
	handler apiHandler
}

// type apiWiring is the state needed to dispatch an incoming API call.
type apiWiring struct {
	pathsMap map[string]verbMap
}

// newApiWiring returns an API configuration, after adding the APIs
// from the given table.  basePath is prefixed to the paths
// specified in the table items.
func newApiWiring(basePath string, tab []apiDesc) (*apiWiring) {  // nolint
	pm := make(map[string]verbMap)
	apiws := &apiWiring{pm}
	for _, b := range(tab) {
		apiws.AddApi(basePath + b.path, b.verb, b.handler)
	}
	return apiws
}

// AddApi() configures the wiring for one path and verb to their handler.
func (apiws *apiWiring) AddApi(path string,
		verb string,
		handler apiHandler) { // nolint
	vmap, ok := apiws.pathsMap[path]
	if !ok {
		vmap = verbMap{path: path, methods: map[string]apiHandler{}}
		apiws.pathsMap[path] = vmap
	}
	vmap.methods[verb] = handler
}

// GetMaps returns the configured path-to-verbMap mapping,
// for possible range iteration.
func (apiws *apiWiring) GetMaps() map[string]verbMap {
	return apiws.pathsMap
}

// callApiMethod() calls the handler that was configured for
// the given verbMap and verb.
func callApiMethod(vmap verbMap, verb string, harg *apiHandlerArg) apiHandlerRet {
	verbFunc, ok := vmap.methods[verb]
	if !ok {
		return apiHandlerRet{http.StatusMethodNotAllowed,
			fmt.Errorf(`No handler for %s on %s`, verb, vmap.path)}
	}

	return verbFunc(harg)
}

// allowedMethods() returns a sorted list of the names of the
// methods allowed by a verbMap.
func allowedMethods(vmap verbMap) []string {
	ret := make([]string, len(vmap.methods))
	i := 0
	for verb := range vmap.methods {
		ret[i] = verb
		i++
	}
	sort.Strings(ret)
	return ret
}

// pathDispatch() is the general handler for all our APIs.
// it is called indirectly thru a closure function that
// supplies the vmap argument.
func pathDispatch(vmap verbMap, w http.ResponseWriter, harg *apiHandlerArg) {
	log.Debugf("in pathDispatch: method=%s path=%s",
		harg.req.Method, harg.req.URL.Path)
	defer func() {
		_ = harg.bodyClose()
	}()

	res := callApiMethod(vmap, harg.req.Method, harg)
	if res.code == http.StatusMethodNotAllowed {
		w.Header().Set("Allowed",
			strings.Join(allowedMethods(vmap), ","))
	}

	rawdata, err := convData(res.data)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(res.code)
	w.Write(rawdata)	// nolint

	log.Debugf("in pathDispatch: code=%d", res.code)
}

// convData() converts the interface{} data part returned by
// an apiHandler function.  the return value is a byte slice.
// if the data is json, it essentially gets ascii-fied.
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
	data, _ := convData(ErrorResponse{code,msg,"ErrorResponse"})

        w.WriteHeader(code)
        _, _ = w.Write(data)

        log.Errorf("error handling API request: %s", msg)
}

// ----- methods for apiHandlerArg

// formValue() is a wrapper for http.Request.FormValue().
// it returns the value, if any, of the named parameter from
// the query portion of the request's URL (or body params).
func (harg *apiHandlerArg) formValue(name string) string {
	return harg.req.FormValue(name)
}

// getBody() is an accessor for http.Request.Body .
func (harg *apiHandlerArg) getBody() io.ReadCloser {
	return harg.req.Body
}

// bodyClose() is a wrapper for http.Request.Body.Close()
func (harg *apiHandlerArg) bodyClose() error {
	return harg.req.Body.Close()
}

// mkApiHandlerArg() takes an http.Request and a map of path variables,
// and returns the corresponding apiHandlerArg object.
func mkApiHandlerArg(req *http.Request,
		pathParams map[string]string) *apiHandlerArg {
	err := req.ParseForm()
	return &apiHandlerArg{req, pathParams, err}
}
