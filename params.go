package apidCRUD

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"unicode"
	"github.com/30x/apid-core"
)

// the type of parameter validator function
type param_validator func (value string) (string, error)

// maxRecs is the maximum number of results allowed in a single bulk request.
const maxRecs = 1000

// map from param name to validator function
var param_validators = map[string]param_validator {
	"table_name": validate_table_name,
	"fields": validate_fields,
	"id": validate_id,
	"id_field": validate_id_field,
	"ids": validate_ids,
	"limit": validate_limit,
	"offset": validate_offset,
}


// ----- start of functions

func fetch_params(req *http.Request, names ...string) (map[string]string, error) {
	ret := map[string]string{}

	// make the query params available via FormValue().
	err := req.ParseForm()
	if err != nil {
		return ret, err
	}

	// make path params available thru path_params[]
	path_params := apid.API().Vars(req)

	// fetch and validate each named param, storing values in ret[]
	for _, name := range names {
		val, err := fetch_param(req, path_params, name)
		if err != nil {
			return ret, err
		}
		ret[name] = val
	}

	return ret, nil
}

func fetch_param(req *http.Request,
		path_params map[string]string,
		name string) (string, error) {
	switch (name) {
	case "table_name": // table_name comes from path_params[] only
		return validate_table_name(path_params[name])
	case "id":	// id may come from path_params[] or FormValue()
		id, ok := path_params[name]
		if !ok {
			id = req.FormValue(name)
		}
		return validate_id(id)
	default:	// param comes from FormValue()
		val := req.FormValue(name)
		vfunc, ok := param_validators[name]
		if ! ok {
			return val, fmt.Errorf("no validator for %s", name)
		}
		return vfunc(val)
	}
}

// ----- param validator functions compatible with param_validator type

func validate_fields(fields string) (string, error) {
	log.Debugf("... fields = %s", fields)
	if fields == "" {
		return "*", nil
	} else {
		for _, f := range strings.Split(fields, ",") {
			if ! isValidIdent(f) {
				return fields, fmt.Errorf("illegal field name")
			}
		}
		return fields, nil
	}
}

func validate_table_name(table_name string) (string, error) {
	log.Debugf("... table_name = %s", table_name)
	if table_name == "" || ! isValidIdent(table_name) {
		return table_name, fmt.Errorf("invalid table name %s", table_name)
	}
	return table_name, nil
}

func validate_id_field(id_field string) (string, error) {
	log.Debugf("... id_field = %s", id_field)
	if id_field == "" {
		id_field = "id"
	}
	if ! isValidIdent(id_field) {
		return "", fmt.Errorf("invalid id_field %s", id_field)
	}
	return id_field, nil
}

func validate_id(id string) (string, error) {
	log.Debugf("... id = %s", id)
	if !isValidIdent(id) {
		return id, fmt.Errorf("invalid id")
	}
	return id, nil
}

func validate_ids(ids string) (string, error) {
	log.Debugf("... ids = %s", ids)
	if ids == "" {
		// an empty list is valid.
		return ids, nil
	}
	idlist := strings.Split(ids, ",")
	for _, id := range idlist {
		// verify that each string is a valid number
		_, err := strconv.Atoi(id)
		if err != nil {
			return ids, err
		}
	}
	return strings.Join(idlist, ","), nil
}

func validate_limit(s string) (string, error) {
	log.Debugf("... limit = %s", s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return s, err
	}
	if n <= 0 || n > maxRecs {
		n = maxRecs
	}
	return strconv.Itoa(n), nil
}

// validate_offset() checks the given string for validity as an SQL offset.
// the empty string is valid and means the default 0.
// a nonempty string must be a number.
// if the input string is valid, a string is returned, and the error is nil.
// an invalid string will result in the error being non-nil.
func validate_offset(s string) (string, error) {
	log.Debugf("... offset = %s", s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return s, err
	}
	if n <= 0 {
		n = 0
	}
	return strconv.Itoa(n), nil
}

// ----- misc validation support functions

// notIdentChar() returns true iff the given rune is not valid in an
// SQL identifier.
func notIdentChar(r rune) bool {
	return !(r == '_' ||
		unicode.In(r, unicode.Digit, unicode.Letter))
}

// isValidIdent() returns true iff the given string is considered a valid
// field identifier in SQL.  s must be nonempty, and must contain only
// chars that from the valid set (notIdentChar).
func isValidIdent(s string) bool {
	return len(s) > 0 && strings.IndexFunc(s, notIdentChar) < 0
}

