package apidCRUD

import (
	// "fmt"
	"net/http"
	"encoding/json"
	"github.com/30x/apid-core"
)

var (
	log apid.LogService
	apiws *apiWiring
)

// initPlugin() is called by the apid-core startup
func initPlugin(services apid.Services) (apid.PluginData, error) {
	log = services.Log().ForModule(pluginData.Name)
	log.Printf("in initPlugin")

	initConfig()

	initDB()

	apiws = initWiring()

	registerHandlers(services.API())

	return pluginData, nil
}

func registerHandlers(service apid.APIService) {
	for path, id := range apiws.pathsMap {
		addHandler(service, apiws, path, id)
	}
}

func dispatch(apiws *apiWiring, pathid int, w http.ResponseWriter, req *http.Request) {
	log.Debugf("in dispatch: method=%s path=%s", req.Method, req.URL.Path)
	defer func() {
		req.Body.Close()
	}()

	verbFunc, err := apiws.getFunc(pathid, req.Method)
	if err != nil {
		errorResponse(w, err)
		return
	}

	code, data := verbFunc(req)

	rawdata, err := convData(data)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(code)
	w.Write(rawdata)

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

func addHandler(service apid.APIService, apiws *apiWiring, path string, pathid int) {
	service.HandleFunc(path,
		func(w http.ResponseWriter, r *http.Request) {
			dispatch(apiws, pathid, w, r)
		})
}

func errorResponse(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := err.Error()
	data, _ := convData(ErrorResponse{code,msg})

        w.WriteHeader(code)
        w.Write(data)

        log.Errorf("error handling API request: %s", msg)
}

func conf_get(cfg apid.ConfigService, vname string, defval string) string {
	ret := cfg.GetString(vname)
	if ret == "" {
		return defval
	} else {
		return ret
	}
}

// initConfig() sets up some global configuration parameters for this plugin.
func initConfig() {
	cfg := apid.Config()

	dbName := conf_get(cfg, "apidCRUD_db_name", "apidCRUD.db")
	log.Debugf("cfg_db_name = %s", dbName)

	base_path := conf_get(cfg, "apidCRUD_base_path", "/apid")
	log.Debugf("base_path = %s", base_path)
}
