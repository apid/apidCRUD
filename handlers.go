package apidCRUD

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"database/sql"
)

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
	result, err := runQuery(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err)
	}
	ret, err := convTableNames(result)
	if err != nil {
		return errorRet(badStat, err)
	}

	return http.StatusOK, TablesResponse{Names: ret}
}

// createDbRecordsHandler() handles POST requests on /db/_table/{table_name} .
func createDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name")
	if err != nil {
		return errorRet(badStat, err)
	}

	body, err := getBodyRecord(req)
	if err != nil {
		return badStat, err
	}
	log.Debugf("body = %s", body)

	records := body.Records
	idlist := make([]int64, 0, len(records))
	log.Debugf("... idlist = %s", idlist)

	err = validateRecords(records)
	if err != nil {
		return badStat, err
	}

	for _, rec := range records {
		// log.Debugf("rec = (%T) %s", rec, rec)
		id, err := runInsert(db, params["table_name"], rec.Keys, rec.Values)
		if err != nil {
			return badStat, err
		}
		idlist = append(idlist, int64(id))
	}

	return http.StatusCreated, IdsResponse{Ids: idlist}
}

// getDbRecordsHandler() handles GET requests on /db/_table/{table_name} .
func getDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req,
		"table_name", "fields", "id_field", "ids", "limit", "offset")
	if err != nil {
		return errorRet(badStat, err)
	}

	return getCommon(params)
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

	return getCommon(params)
}

// updateDbRecordsHandler() handles PATCH requests on /db/_table/{table_name} .
func updateDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err)
	}
	return updateCommon(req, params)
}

// updateDbRecordHandler() handles PATCH requests on /db/_table/{table_name}/{id} .
func updateDbRecordHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err)
	}
	return updateCommon(req, params)
}

// deleteDbRecordsHandler handles DELETE requests on /db/_table/{table_name} .
func deleteDbRecordsHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err)
	}

	return delCommon(params)
}

// deleteDbRecordHandler handles DELETE requests on /db/_table/{table_name}/{id} .
func deleteDbRecordHandler(req *http.Request) (int, interface{}) {
	params, err := fetchParams(req, "table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err)
	}
	return delCommon(params)
}

// getDbSchemasHandler handles GET requests on /db/_schema .
func getDbSchemasHandler(req *http.Request) (int, interface{}) {
	return notImplemented()
}

// createDbTableHandler handles POST requests on /db/_schema .
func createDbTableHandler(req *http.Request) (int, interface{}) {
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

// errorRet() is called by apiHandler routines to pass back the code/data
// pair appropriate to the given code and error object.
func errorRet(code int, err error) (int, interface{}) {
	return code, ErrorResponse{code, err.Error()}
}

// notImpemented() returns the code/data pair for an apiHandler
// that is not implemented.
func notImplemented() (int, interface{}) {
	code := http.StatusNotImplemented
	return errorRet(code, fmt.Errorf("API not implemented yet"))
}

// initDB() initializes the global db variable
func initDB() (dbType, error) {
	h, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return dbType{}, err
	}

	// assign the global db variable
	return dbType{handle: h}, nil
}

// mkVmap() takes a list of keys (string) and a list of values.
// the values are *sql.RawBytes as interface{}.
// and returns a map from key to corresponding value.
// these map values are *string as interface{}.
// this is somewhat like python's zip function.
func mkVmap(keys []string, values []interface{}) (map[string]interface{}, error) {
	N := len(keys)
	if N != len(values) {
		return nil, fmt.Errorf("nkeys different from nvalues")
	}
	ret := make(map[string]interface{}, N)
	for i := 0; i < N; i++ {
		// convert from sql.RawBytes to string
		rbp, ok := values[i].(*sql.RawBytes)
		if !ok {
			return ret, fmt.Errorf("sql conversion error")
		}
		s := string(*rbp)
		ret[keys[i]] = interface{}(&s)
	}
	return ret, nil
}

// mkSQLRow() returns a list of interface{} of the given length,
// each element is actually a pointer to sql.RawBytes .
func mkSQLRow(N int) []interface{} {
	ret := make([]interface{}, N)
	for i := 0; i < N; i++ {
		ret[i] = new(sql.RawBytes)
	}
	return ret
}

// runQuery() does a select query using the given query string.
// the return value is a list of the retrieved records.
func runQuery(db dbType,
		qstring string,
		ivals []interface{}) ([]map[string]interface{}, error) {
	log.Debugf("query = %s", qstring)

	ret := make([]map[string]interface{}, 0, maxRecs)

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

		vals := mkSQLRow(len(cols))
		err := rows.Scan(vals...)
		if err != nil {
			return ret, fmt.Errorf("scan error at rownum %d", rownum)
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

// strListToInterfaces() converts a list of strings to a list of interface{}.
func strListToInterfaces(vals []string) []interface{} {
	ret := make([]interface{}, len(vals))
	for i, v := range vals {
		ret[i] = interface{}(v)
	}
	return ret
}

// idTypesToInterface() convert a list of strings to
// a list of database id's (of idType) disguised as interface{}.
func idTypesToInterface(vals []string) []interface{} {
	ret := make([]interface{}, len(vals))
	for i, v := range vals {
		ret[i] = interface{}(aToIdType(v))
	}
	return ret
}

// nstring() returns a string with n comma-separated copies of
// the given string s.
func nstring(s string, n int) string {
	ret := make([]string, n, n)
	for i := 0; i < n; i++ {
		ret[i] = s
	}
	return strings.Join(ret, ",")
}

// runInsert() inserts a record whose data is specified by the
// given keys and values.  it returns the id of the inserted record.
func runInsert(db dbType, tabname string, keys []string, values []string) (idType, error) {
	NORET := idType(-1)

	nkeys := len(keys)
	nvalues := len(values)
	if nkeys != nvalues {
		return NORET, fmt.Errorf("number of keys must equal number of values")
	}

	keystr := strings.Join(keys, ",")
	placestr := nstring("?", nvalues)

	qstring := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tabname, keystr, placestr)

	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}
	defer func() {
		stmt.Close()
	}()

	ivalues := strListToInterfaces(values)

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

// delCommon() is the common part of record deletion APIs.
func delCommon(params map[string]string) (int, interface{}) {
	nc, err := delRecs(db, params)
	if err != nil {
		return errorRet(badStat, err)
	}

	return http.StatusOK, NumChangedResponse{nc}
}

// delRecs() deletes multiple records, using parameters in the params map.
// it returns the number of records deleted.
func delRecs(db dbType, params map[string]string) (int64, error) {
	NORET := int64(-1)
	idclause, idlist, err := mkIdClause(params)
	if err != nil {
		return NORET, err
	}
	if idclause == "" {
		return NORET,
			fmt.Errorf("id or ids must be specified")
	}
	qstring := fmt.Sprintf("DELETE FROM %s %s",
		params["table_name"],
		idclause)
	log.Debugf("qstring = %s\n", qstring)

	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}
	defer func() {
		stmt.Close()
	}()

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
	return ra, nil
}

// validateSQLKeys() checks an array of key names,
// returning a non-nil error if anything is found that
// would not be a valid SQL key.
func validateSQLKeys(keys []string) error {
	for _, k := range keys {
		if !isValidIdent(k) {
			return fmt.Errorf("invalid key %s", k)
		}
	}
	return nil
}

// validateSQLValues() checks an array of string values,
// returning a non-nil error if anything is found that
// would not be a valid SQL value.
func validateSQLValues(values []string) error {
	return nil    // no error for now
}

// getBodyRecord() returns a json record from the body of the given request.
func getBodyRecord(req *http.Request) (BodyRecord, error) {
	jrec := BodyRecord{}
        err := json.NewDecoder(req.Body).Decode(&jrec)
	return jrec, err
}

// mkIdClause() takes the API parameters,
// and returns the implied WHERE clause that can be
// plugged in to a query string (for use with Prepare)
// and list of data items (for use with Exec).
// the params examined include id_field, id, and/or ids.
// if id is specified, that is used; otherwise ids.
// if neither id nor ids is specified, the WHERE clause is empty.
func mkIdClause(params map[string]string) (string, []interface{}, error) { // nolint
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		idlist := []interface{}{aToIdType(id)}
		placestr := "?"
		idclause := fmt.Sprintf("WHERE %s = %s", id_field, placestr)
		return idclause, idlist, nil
	}

	ids, ok := params["ids"]
	if ok && ids != "" {
		idstrings := strings.Split(ids, ",")
		idlist := idTypesToInterface(idstrings)
		placestr := nstring("?", len(idlist))
		idclause := fmt.Sprintf("WHERE %s in (%s)", id_field, placestr)
		return idclause, idlist, nil
	}

	// no id and no ids implies everything matches.
	// if this is bad, caller should check.
	return "", []interface{}{}, nil
}

// mkIdClauseUpdate() is like mkIdClause(), but for UPDATE operations.
// the difference is, for now, that the id values are formatted directly
// into the WHERE string, rather than being subbed in by Exec.
// don't allow the case where neither id nor ids is specified.
func mkIdClauseUpdate(params map[string]string) (string, error) {  // nolint
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		return fmt.Sprintf("WHERE %s = %s", id_field, id), nil
	}
	ids, ok := params["ids"]
	if ok && ids != "" {
		return fmt.Sprintf("WHERE %s in (%s)", id_field, ids), nil
	}

	// no id and no ids implies everything matches.
	// if this is bad, caller should check.
	return "", nil
}

// updateRec() updates certain fields of a given record or records,
// using parameters in the params map.
// it returns the number of records changed.
func updateRec(db dbType,
		params map[string]string,
		body BodyRecord) (int64, error) {
	NORET := int64(-1)
	dbrec := body.Records[0]
	keylist := dbrec.Keys
	keystr := strings.Join(keylist, ",")
	placestr := nstring("?", len(keylist))
	idclause, err := mkIdClauseUpdate(params)
	if err != nil {
		return NORET, err
	}
	if idclause == "" {
		return NORET,
			fmt.Errorf("id or ids must be specified")
	}

	qstring := fmt.Sprintf("UPDATE %s SET (%s) = (%s) %s",
			params["table_name"],
			keystr,
			placestr,
			idclause)

	log.Debugf("qstring = %s", qstring)
	stmt, err := db.handle.Prepare(qstring)
	if err != nil {
		return NORET, err
	}
	defer func() {
		stmt.Close()
	}()
	ivals := strListToInterfaces(dbrec.Values)
	result, err := stmt.Exec(ivals...)
	if err != nil {
		return NORET, err
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return NORET, err
	}
	return ra, nil
}

// mkSelectString() returns the WHERE part of a selection query.
func mkSelectString(params map[string]string) (string, []interface{}, error) {
	idclause, idlist, err := mkIdClause(params)
	if err != nil {
		return idclause, idlist, err
	}

	qstring := fmt.Sprintf("SELECT %s FROM %s %s LIMIT %s OFFSET %s",
		params["fields"],
		params["table_name"],
		idclause,
		params["limit"],
		params["offset"])

	return qstring, idlist, nil
}

// getCommon() is common code for selection APIs.
func getCommon(params map[string]string) (int, interface{}) {
	qstring, idlist, err := mkSelectString(params)
	if err != nil {
		return errorRet(badStat, err)
	}
	result, err := runQuery(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err)
	}

	if len(result) == 0 {
		return errorRet(badStat, fmt.Errorf("no matching record"))
	}

	return http.StatusOK, RecordsResponse{Records:result}
}

// updateCommon is common code for update APIs.
func updateCommon(req *http.Request, params map[string]string) (int, interface{}) {
	body, err := getBodyRecord(req)
	if err != nil {
		return errorRet(badStat, err)
	}
	if len(body.Records) < 1 {
		return errorRet(badStat,
			fmt.Errorf("update: no data records in body"))
	}

	ra, err := updateRec(db, params, body)
	if err != nil {
		return errorRet(badStat, err)
	}
	return http.StatusOK, NumChangedResponse{ra}
}

// convTableNames() converts the return format from runQuery()
// into a simple list of names.
func convTableNames(result []map[string]interface{}) ([]string, error) {
	// convert from query format to simple list of names
	ret := make([]string, len(result))
	for i, row := range result {
		name, err := grabNameField(row)
		if err != nil {
			return ret, err
		}
		ret[i] = name
	}
	return ret, nil
}

// grabNameField() grabs the name field from this row object.
func grabNameField(row map[string]interface{}) (string, error) {
	obj := row["name"]
	pname, ok := obj.(*string)
	if !ok {
		return "", fmt.Errorf("sql table name conversion error")
	}
	return *pname, nil
}

// validateRecords() checks the validity of an array of KVRecord.
// returns an error if any record has an invalid key or value.
func validateRecords(records []KVRecord) error {
	for i, rec := range records {
		// log.Debugf("rec = (%T) %s", rec, rec)
		keys := rec.Keys
		values := rec.Values
		if len(keys) != len(values) {
			return fmt.Errorf("Record %d nkeys != nvalues", i)
		}
		err := validateSQLKeys(keys)
		if err != nil {
			return err
		}
		err = validateSQLValues(values)
		if err != nil {
			return err
		}
	}
	return nil
}
