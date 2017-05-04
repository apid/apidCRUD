package apidCRUD

import (
	"fmt"
	"strings"
	"strconv"
	"unicode"
)

// idType is an alias for the type of the database's rowid.
type idType int64

const (
	// idTypeBits is the number of bits in idType.
	idTypeBits = 64

	// idTypeRadix is the base to use when converting an id string to int.
	idTypeRadix = 10

	// types of parameters we recognize in an http request.
	paramQuery = 0		// comes from query portion of the URL.
	paramPathOnly = 1	// comes from the path portion of the URL.
	paramPathOrQuery = 2	// may come from path or query.
)

// maxRecs is the max number of results allowed in a bulk request.
var maxRecs = 1000

// the type of parameter validator function
type paramValidator func (value string) (string, error)

// map from param name to validator function
var validators = map[string]paramValidator {
	"table_name": validate_table_name,
	"fields": validate_fields,
	"id": validate_id,
	"id_field": validate_id_field,
	"ids": validate_ids,
	"limit": validate_limit,
	"offset": validate_offset,
}

// paramType tells which parameters come from where.
// by default, parameters not listed here are paramQuery.
var paramType = map[string]int {
	"table_name": paramPathOnly,
	"id": paramPathOrQuery,
}

// ----- start of functions

// fetch_param() fetches the named parameter from the Request as a string.
// the parameter must have a validator function.
// the call fails if the validator function fails.
func (harg *apiHandlerArg) getParam(name string) (string, error) {
	var val string
	var ok bool
	switch (paramType[name]) {
	case paramPathOnly:
		val = harg.pathParams[name]
	case paramPathOrQuery:
		val, ok = harg.pathParams[name]
		if !ok {
			val = harg.formValue(name)
		}
	default:
		val = harg.formValue(name)
	}
	vfunc, ok := validators[name]
	if ! ok {
		return val, fmt.Errorf("no validator for %s", name)
	}
	return vfunc(val)
}

// fetchParams() gets the named parameters from the given Request.
// the parameters may be in the path or in the query.
// each parameter must have a validator function.
// the call returns an error if a validator function fails on any parameter.
// the parameter values are returned as a map of string.
func fetchParams(harg apiHandlerArg, names ...string) (map[string]string, error) {
	ret := map[string]string{}

	// fetch and validate each named param, storing values in ret[]
	for _, name := range names {
		val, err := harg.getParam(name)
		if err != nil {
			return ret, err
		}
		ret[name] = val
	}

	return ret, nil
}

// ----- param validator functions compatible with paramValidator type

// validate_fields() is the validator for the "fields" parameter.
func validate_fields(fields string) (string, error) {
	log.Debugf("... fields = %s", fields)
	if fields == "" {
		return "*", nil
	}
	for _, f := range strings.Split(fields, ",") {
		if ! isValidIdent(f) {
			return fields, fmt.Errorf("illegal field name")
		}
	}
	return fields, nil
}

// validate_table_name() is the validator for the "table_name" parameter.
func validate_table_name(table_name string) (string, error) {
	log.Debugf("... table_name = %s", table_name)
	if table_name == "" || ! isValidIdent(table_name) {
		return table_name, fmt.Errorf("invalid table name %s", table_name)
	}
	return table_name, nil
}

// validate_id_field() is the validator for the "id_field" parameter.
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

// validate_id() validates the given string as an SQL id value.
// it must be a valid nonempty numeric string.
func validate_id(id string) (string, error) {
	log.Debugf("... id = %s", id)
	n, err := strconv.ParseInt(id, idTypeRadix, idTypeBits)
	if err != nil {
		return id, err
	}
	return idTypeToA(n), nil
}

// validate_ids() validates the given string as a comma-separated list
// of SQL id values.  each item must be a valid numeric string.
// the empty string is valid and means the empty list.
func validate_ids(ids string) (string, error) {
	log.Debugf("... ids = %s", ids)
	if ids == "" {
		// an empty list is valid.
		return ids, nil
	}

	idlist := strings.Split(ids, ",")
	nids := len(idlist)
	for k := 0; k < nids; k++ {
		// verify that each item is a valid numeric string
		n, err := strconv.ParseInt(idlist[k], idTypeRadix, idTypeBits)
		if err != nil {
			return ids, err
		}
		// store back in normalized form
		idlist[k] = idTypeToA(n)
	}

	return strings.Join(idlist, ","), nil
}

// validate_limit() checks the given string for validity as an SQL limit.
// an empty string is valid and means the default 0.
// a negative number or a number greater than maxRecs, is valid
// and means maxRecs.
func validate_limit(s string) (string, error) {
	log.Debugf("... limit = %s", s)
	if s == "" {
		s = "0"
	}
	n, err := strconv.ParseInt(s, idTypeRadix, idTypeBits)
	if err != nil {
		return s, err
	}

	imax := int64(maxRecs)
	if n <= 0 || n > imax {
		n = imax
	}
	return idTypeToA(n), nil
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
	n, err := strconv.ParseInt(s, idTypeRadix, idTypeBits)
	if err != nil {
		return s, err
	}
	return idTypeToA(n), nil
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
	if len(s) == 0 {
		return false
	}
	r := rune(s[0])
	return (r == '_' || unicode.In(r, unicode.Letter)) &&
		strings.IndexFunc(s, notIdentChar) < 0
}

// aToIdType() converts a string to idType.
// on error, return -1.  note that -1 is also a legitimate value,
// so should use this only on strings that are known to be valid.
func aToIdType(idstr string) int64 {	// nolint
	id, err := strconv.ParseInt(idstr, idTypeRadix, idTypeBits)
	if err != nil {
		return -1
	}
	return id
}

// idTypeToA() converts an idType value to a string, in the standard way.
func idTypeToA(val int64) string {
	return strconv.FormatInt(val, idTypeRadix)
}
