package apidCRUD

import (
	"testing"
	"fmt"
	"strings"
	"net/http"
)

func mkRequest(path string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, path, nil)
}

func mkExtReq(path string) (*extReq, error) {
	req, err := mkRequest(path)
	if err != nil {
		return nil, err
	}
	return newExtReq(req, validators)
}

// ---- unit tests for extReqNew()

func Test_newExtReq(t *testing.T) {
	fn := "newExtReq"
	xr, err := mkExtReq("/apid/db")
	if err != nil {
		t.Errorf("%s failure: %s", fn, err)
		return
	}
	if xr == nil {
		t.Errorf("%s returned nil", fn)
	}
}

// ----- unit tests for getParam()

func getParamHelper(t *testing.T,
		paramName string,
		val string) (string, error) {
	path := fmt.Sprintf("/apid/db?%s=%s", paramName, val)
	xr, err := mkExtReq(path)
	if err != nil {
		return "", nil
	}
	return xr.getParam(paramName)
}

func Test_getParam(t *testing.T) {

	// test getParam on id values
	run_validator(t,
		func(val string) (string, error) {
			return getParamHelper(t, "id", val)
		},
		validate_id_Tab)

	// test getParam on ids values
	run_validator(t,
		func(val string) (string, error) {
			return getParamHelper(t, "ids", val)
		},
		validate_ids_Tab)

	// test getParam on id_field values
	run_validator(t,
		func(val string) (string, error) {
			return getParamHelper(t, "id_field", val)
		},
		validate_id_field_Tab)

	// test getParam on a field with no validator
	run_validator(t,
		func(val string) (string, error) {
			return getParamHelper(t, "nofield", val)
		},
		validate_nofield_Tab)
}

// ----- unit tests for fetchParams()

type fetchParams_TC struct {
	arg string	// query params to use in call
	xsucc bool	// expected success
}

var fetchParams_Tab = []fetchParams_TC {
	{ "id=123", true },
	{ "id=123&ids=123,456", true },
	{ "id=1&fields=a,b,c", true },
	{ "junk=1&fields=a,b,c", false },
}

func fetchParamsHelper(qp string) (map[string]string, error) {
	qplist := strings.Split(qp, "&")
	names := make([]string, len(qplist))
	for i, parm := range qplist {
		nv := strings.SplitN(parm, "=", 2)
		names[i] = nv[0]
	}

	req, err := mkRequest("/api/db?" + qp)
	if err != nil {
		vmap := map[string]string{}
		return vmap, err
	}

	vmap, err := fetchParams(req, names...)
	if err != nil {
		return vmap, err
	}

	// check that the map has the expected number of keys
	nvmap := len(vmap)
	nnames := len(names)
	if nvmap != nnames {
		err := fmt.Errorf("map has %d entries, expected %d",
				nvmap, nnames)
		return vmap, err
	}

	// check that each expected name is there
	for _, name := range names {
		_, ok := vmap[name]
		if !ok {
			err := fmt.Errorf("map does not have %s", name)
			return vmap, err
		}
	}

	return vmap, nil
}

func call_fetchParams(t *testing.T, i int, qp string, xsucc bool) {
	_, err := fetchParamsHelper(qp)
	if xsucc != (err == nil) {
		msg := "true"
		if err != nil {
			msg = err.Error()
		}
		t.Errorf(`#%d: fetchParams("%s")=(%s), expected (%t)`,
			i, qp, msg, xsucc)
	}
}

func Test_fetchParams(t *testing.T) {
	for i, test := range fetchParams_Tab {
		call_fetchParams(t, i, test.arg, test.xsucc)
	}
}
