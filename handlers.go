package apidCRUD

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"database/sql"
)

// dbType is intended to encapsulate the database handle type.
type dbType struct {
	handle *sql.DB
}

// badStat is a convenience constant, the http status for a bad request.
const badStat = http.StatusBadRequest

// dbName is the name of the database that is implicitly used in these APIs.
var dbName = "apidCRUD.db"

// DEBUG enables debugging printouts
var DEBUG = true

// db is our global database handle
var db dbType

// ----- plain old handlers that are compatible with the apiHandler type.

// getDbResourcesHandler handles GET requests on /db
func getDbResourcesHandler(req *http.Request) (int, interface{}) {
	// not sure what this should do
	return notImplemented()
}

// getDbTablesHandler handles GET requests on /db/_table
func getDbTablesHandler(req *http.Request) (int, interface{}) {
	// the "tables" table is our convention, not maintained by sqlite.

	idlist := []interface{}{}
	qstring := "select name from tables;"
	result, err := myselect(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err)
	}

	// convert from query format to simple list of names
	ret := make([]string, len(result))
	for i, tab := range result {
		obj := (*tab)["name"]
		pname, ok := obj.(*string)
		if !ok {
			return errorRet(badStat, fmt.Errorf("conversion error"))
		}
		ret[i] = *pname
	}

	return http.StatusOK, TablesResponse{Resource: ret}
}

// createDbRecordsHandler() handles POST requests on /db/_table/{table_name} .
func createDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name")
	if err != nil {
		return errorRet(badStat, err)
	}

	jrec, err := getJsonRecord(req)
	if err != nil {
		return badStat, err
	}
	log.Debugf("jrec = %s", jrec)

	resources := jrec.Resource
	idlist := make([]int, 0, len(resources))
	log.Debugf("... idlist = %s", idlist)

	for _, rec := range resources {
		log.Debugf("rec = (%T) %s", rec, rec)
		keys := rec.Keys
		values := rec.Values
		err := validate_sql_keys(keys)
		if err != nil {
			return badStat, err
		}
		err = validate_sql_values(values)
		if err != nil {
			return badStat, err
		}
		id, err := write_rec(db, params["table_name"], keys, values)
		if err != nil {
			return badStat, err
		}
		idlist = append(idlist, int(id))
	}

	return http.StatusOK, RecordIds{Ids: idlist}
}

// getDbRecordsHandler() handles GET requests on /db/_table/{table_name} .
func getDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "fields", "id_field", "ids", "limit", "offset")
	if err != nil {
		return errorRet(badStat, err)
	}

	return get_common(params)
}

// getDbRecordHandler() handles GET requests on /db/_table/{table_name}/{id} .
func getDbRecordHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "id", "fields", "id_field")
	if err != nil {
		return errorRet(badStat, err)
	}
	params["limit"] = strconv.Itoa(1)
	params["offset"] = strconv.Itoa(0)

	return get_common(params)
}

// updateDbRecordsHandler() handles PATCH requests on /db/_table/{table_name} .
func updateDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err)
	}

	return update_common(req, params)
}

// updateDbRecordHandler() handles PATCH requests on /db/_table/{table_name}/{id} .
func updateDbRecordHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err)
	}
	return update_common(req, params)
}

// deleteDbRecordsHandler handles DELETE requests on /db/_table/{table_name} .
func deleteDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err)
	}

	return del_common(params)
}

// deleteDbRecordHandler handles DELETE requests on /db/_table/{table_name}/{id} .
func deleteDbRecordHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err)
	}
	return del_common(params)
}

// getDbSchemasHandler handles GET requests on /db/_schema .
func getDbSchemasHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// createDbTableHandler handles POST requests on /db/_schema .
func createDbTableHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// replaceDbTables handles PUT requests on /db/_schema .
func replaceDbTablesHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// updateDbTables handles PATCH requests on /db/_schema .
func updateDbTablesHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// describeDbTableHandler handles GET requests on /db/_schema/{table_name} .
func describeDbTableHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// createDbTablesHandler handles POST requests on /db/_schema/{table_name} .
func createDbTablesHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// deleteDbTableHandler handles DELETE requests on /db/_schema/{table_name} .
func deleteDbTableHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// describeDbFieldHandler handles GET requests on /db/_schema/{table_name} .
func describeDbFieldHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// ----- misc support functions

// errorRet is called by apiHandler routines to pass back the code/data
// pair appropriate to the given code and error object.
func errorRet(code int, err error) (int, interface{}) {
	return code, ErrorResponse{code, err.Error()}
}

func initDB() {
	localdb, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	// assign the global db variable
	db = dbType{handle: localdb}
}

// mkVmap() takes a list of keys (string) and a list of values.
// the values are *sql.RawBytes as interface{}.
// and returns a map from key to corresponding value.
// these map values are *string as interface{}.
func mkVmap(keys []string, values []interface{}) (*map[string]interface{}, error) {
	N := len(keys)
	if N != len(values) {
		return nil, fmt.Errorf("nkeys different from nvalues")
	}
	ret := make(map[string]interface{}, N)
	for i := 0; i < N; i++ {
		// convert from sql.RawBytes to string
		rbp, ok := values[i].(*sql.RawBytes)
		if !ok {
			return &ret, fmt.Errorf("sql conversion error")
		}
		s := string(*rbp)
		ret[keys[i]] = interface{}(&s)
	}
	return &ret, nil
}

func mkSqlRow(N int) []interface{} {
	ret := make([]interface{}, N)
	for i := 0; i < N; i++ {
		ret[i] = new(sql.RawBytes)
	}
	return ret
}

func printRow(rownum int, colnames []string, vals []interface{}) {
	log.Debugf("row #%d: ", rownum)
	sep := ""
	N := len(colnames)
	for i := 0; i < N; i++ {
		var v string
		rbp, ok := vals[i].(*sql.RawBytes)
		if !ok {
			v = "?"
		} else {
			v = string(*rbp)
		}
		log.Debugf("%s%s=%s", sep, colnames[i], v)
		sep = " | "
	}
	log.Debugf("\n")
}

func myselect(db dbType, qstring string, ivals []interface{}) ([]*map[string]interface{}, error) {
	log.Debugf("query = %s", qstring)

	ret := make([]*map[string]interface{}, 0, maxRecs)

	rownum := 0
	rows, err := db.handle.Query(qstring, ivals...)
	if err != nil {
		return ret, err
	}

	// ensure rows gets closed at end
	defer func() {
		_ = rows.Close()
	}()

	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		return ret, err
	}
	log.Debugf("cols = %s", cols)

	for rows.Next() {
		rownum ++

		vals := mkSqlRow(len(cols))
		err := rows.Scan(vals...)
		if err != nil {
			return ret, fmt.Errorf("scan error at rownum %d", rownum)
		}

		if DEBUG {
			printRow(rownum, cols, vals)
		}

		m, err := mkVmap(cols, vals)
		if err != nil {
			return ret, err
		}
		ret = append(ret, m)
	}

	if rows.Err() != nil {
		return ret, fmt.Errorf("rows ended with error at rownum %d", rownum)
	}

	return ret, nil
}

// convert a list of strings to a list of interface{}.
func strToInterface(vals []string) []interface{} {
	ret := make([]interface{}, len(vals))
	for i, v := range vals {
		ret[i] = interface{}(v)
	}
	return ret
}

func atoIdType(idstr string) int64 {
	id, _ := strconv.ParseInt(idstr, idTypeRadix, idTypeBits)
	return id
}

// convert a list of strings to a list database id's disguised as interface{}.
func idTypeToInterface(vals []string) []interface{} {
	ret := make([]interface{}, len(vals))
	for i, v := range vals {
		ret[i] = interface{}(atoIdType(v))
	}
	return ret
}


// return a string with n comma-separated copies of the given string s
func nstring(s string, n int) string {
	ret := make([]string, n, n)
	for i := 0; i < n; i++ {
		ret[i] = s
	}
	return strings.Join(ret, ",")
}

func write_rec(db dbType, tabname string, keys []string, values []string) (idType, error) {
	NORET := idType(-1)
	nkeys := len(keys)
	nvalues := len(values)
	if nkeys != nvalues {
		return NORET, fmt.Errorf("number of keys must equal number of values")
	}

	keystr := strings.Join(keys, ",")
	placestr := nstring("?", nvalues)

	qstring := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tabname, keystr, placestr)
	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}

	ivalues := strToInterface(values)

	log.Debugf("qstring = %s\n", qstring)
	result, err := stmt.Exec(ivalues...)
	if err != nil {
		return NORET, err
	}
	// fmt.Debugf("result=%s", result)

	lastid, err := result.LastInsertId()
	if err != nil {
		return NORET, err
	}
	log.Debugf("lastid = %d\n", lastid)
	nrecs, err := result.RowsAffected()
	if err != nil {
		return NORET, err
	}
	log.Debugf("rowsaffected = %d\n", nrecs)
	return idType(lastid), nil
}

func del_common(params map[string]string) (int, interface{}) {
	nc, err := del_recs(db, params)
	if err != nil {
		return errorRet(badStat, err)
	}

	return http.StatusOK, DeleteResponse{nc}
}

func del_recs(db dbType, params map[string]string) (int, error) {
	NORET := -1
	idclause, idlist, err := idclause_setup(params)
	if err != nil {
		return NORET, err
	}
	qstring := fmt.Sprintf("DELETE FROM %s %s;",
		params["table_name"],
		idclause)
	log.Debugf("qstring = %s\n", qstring)

	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}

	result, err := stmt.Exec(idlist...)
	if err != nil {
		return NORET, err
	}
	// log.Debugf("result=%s", result)

	lastid, err := result.LastInsertId()
	if err != nil {
		return NORET, err
	}
	log.Debugf("lastid = %d\n", lastid)

	ra, err := result.RowsAffected()
	if err != nil {
		return NORET, err
	}
	log.Debugf("ra = %d\n", ra)

	if int(ra) != len(idlist) {
		return NORET, fmt.Errorf("mismatch in number of affected records")
	}
	return int(ra), nil
}

func validate_sql_keys(keys []string) error {
	return nil    // no error for now
}

func validate_sql_values(values []string) error {
	return nil    // no error for now
}


func notImplemented() (int, interface{}) {
	code := http.StatusNotImplemented
	return errorRet(code, fmt.Errorf("API not implemented yet"))
}

func getJsonRecord(req *http.Request) (jsonRecord, error) {
	jrec := jsonRecord{}
	err := getBody(req, &jrec)
	return jrec, err
}

func getBody(req *http.Request, jrec *jsonRecord) error {
        err := json.NewDecoder(req.Body).Decode(jrec)
        if err != nil {
                log.Errorf("JSON Response Data not parsable: %v", err)
                return err
        }
	return err
}

func idclause_setup(params map[string]string) (string, []interface{}, error) {
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		idlist := []interface{}{atoIdType(id)}
		placestr := "?"
		idclause := fmt.Sprintf("WHERE %s = %s", id_field, placestr)
		return idclause, idlist, nil
	}
	ids, ok := params["ids"]
	if ok && ids != "" {
		idstrings := strings.Split(ids, ",")
		idlist := idTypeToInterface(idstrings)
		placestr := nstring("?", len(idlist))
		idclause := fmt.Sprintf("WHERE %s in (%s)", id_field, placestr)
		return idclause, idlist, nil
	}

	// no id and no ids implies everything matches.
	return "", []interface{}{}, nil
}

func mk_idclause(params map[string]string) string {
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		return fmt.Sprintf("WHERE %s = %s", id_field, id)
	}
	ids, ok := params["ids"]
	if ok && ids != "" {
		return fmt.Sprintf("WHERE %s in (%s)", id_field, ids)
	}

	// no id and no ids implies everything matches.
	return ""
}

func update_rec(db dbType,
		params map[string]string,
		jrec jsonRecord) (int, error) {
	NORET := -1
	dbrec := jrec.Resource[0]
	keylist := dbrec.Keys
	keystr := strings.Join(keylist, ",")
	placestr := nstring("?", len(keylist))

	qstring := fmt.Sprintf("UPDATE %s SET (%s) = (%s) %s;",
			params["table_name"],
			keystr,
			placestr,
			mk_idclause(params))

	log.Debugf("qstring = %s", qstring)
	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}
	ivals := strToInterface(dbrec.Values)
	result, err := stmt.Exec(ivals...)
	if err != nil {
		return NORET, err
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return NORET, err
	}
	return int(ra), nil
}

func get_common(params map[string]string) (int, interface{}) {
	idclause, idlist, err := idclause_setup(params)
	if err != nil {
		return errorRet(badStat, err)
	}

	qstring := fmt.Sprintf("SELECT %s FROM %s %s LIMIT %s OFFSET %s;",
		params["fields"],
		params["table_name"],
		idclause,
		params["limit"],
		params["offset"])
	result, err := myselect(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err)
	}

	if len(result) == 0 {
		return errorRet(badStat, fmt.Errorf("no matching record"))
	}

	return http.StatusOK, GetRecordResponse{Record:result}
}

func update_common(req *http.Request, params map[string]string) (int, interface{}) {
	jrecs, err := getJsonRecord(req)
	if err != nil {
		return errorRet(badStat, err)
	}

	ra, err := update_rec(db, params, jrecs)
	if err != nil {
		return errorRet(badStat, err)
	}
	return http.StatusOK, DeleteResponse{ra}
}
