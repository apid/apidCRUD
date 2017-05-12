package apidCRUD

import (
	"fmt"
	"strings"
	"strconv"
	"bytes"
	"net/http"
	"encoding/json"
	"database/sql"
)

// ----- plain old handlers that are compatible with the apiHandler type.

// getDbResourcesHandler handles GET requests on /db
func getDbResourcesHandler(harg *apiHandlerArg) apiHandlerRet {
	// not sure what this should do
	return notImplemented()
}

// getDbTablesHandler handles GET requests on /db/_table
func getDbTablesHandler(harg *apiHandlerArg) apiHandlerRet {
	return tablesQuery(harg, "tables", "name")
}

// createDbRecordsHandler() handles POST requests on /db/_table/{table_name} .
func createDbRecordsHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}

	body, err := getBodyRecord(harg)
	if err != nil {
		return apiHandlerRet{badStat, err}
	}
	log.Debugf("body = %s", body)

	records := body.Records
	idlist := make([]int64, 0, len(records))
	log.Debugf("... idlist = %s", idlist)

	err = validateRecords(records)
	if err != nil {
		return apiHandlerRet{badStat, err}
	}

	for _, rec := range records {
		// log.Debugf("rec = (%T) %s", rec, rec)
		id, err := runInsert(db, params["table_name"], rec.Keys, rec.Values)
		if err != nil {
			return apiHandlerRet{badStat, err}
		}
		idlist = append(idlist, int64(id))
	}

	return apiHandlerRet{http.StatusCreated, IdsResponse{Ids: idlist}}
}

// getDbRecordsHandler() handles GET requests on /db/_table/{table_name} .
func getDbRecordsHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg,
		"table_name", "fields", "id_field", "ids", "limit", "offset")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}

	return getCommon(params)
}

// getDbRecordHandler() handles GET requests on /db/_table/{table_name}/{id} .
func getDbRecordHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg,
		"table_name", "id", "fields", "id_field")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	params["limit"] = strconv.Itoa(1)
	params["offset"] = strconv.Itoa(0)

	return getCommon(params)
}

// updateDbRecordsHandler() handles PATCH requests on /db/_table/{table_name} .
func updateDbRecordsHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	return updateCommon(harg, params)
}

// updateDbRecordHandler() handles PATCH requests on /db/_table/{table_name}/{id} .
func updateDbRecordHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	return updateCommon(harg, params)
}

// deleteDbRecordsHandler handles DELETE requests on /db/_table/{table_name} .
func deleteDbRecordsHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name", "id_field", "ids")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	return delCommon(params)
}

// deleteDbRecordHandler handles DELETE requests on /db/_table/{table_name}/{id} .
func deleteDbRecordHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name", "id", "id_field")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	return delCommon(params)
}

// getDbSchemasHandler handles GET requests on /db/_schema .
func getDbSchemasHandler(harg *apiHandlerArg) apiHandlerRet {
	return notImplemented()
}

// createDbTableHandler handles POST requests on /db/_schema/{table_name} .
func createDbTableHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	log.Debugf("... params = %v", params)
	schema, err := getBodySchemas(harg)
	if err != nil {
		return errorRet(badStat, err, "missing or ill-formed schema")
	}
	if len(schema.Resource) != 1 {
		return errorRet(badStat,
			fmt.Errorf("body schema should have exactly 1 spec"),
			"schema count")
	}
	log.Debugf("schema=%v", schema)
	tabNames := []string{params["table_name"]}
	return createTablesCommon(tabNames, schema)
}

// updateDbTables handles PATCH requests on /db/_schema .
func updateDbTablesHandler(harg *apiHandlerArg) apiHandlerRet {
	return notImplemented()
}

// describeDbTableHandler handles GET requests on /db/_schema/{table_name} .
func describeDbTableHandler(harg *apiHandlerArg) apiHandlerRet {
	return notImplemented()
}

// deleteDbTableHandler handles DELETE requests on /db/_schema/{table_name} .
func deleteDbTableHandler(harg *apiHandlerArg) apiHandlerRet {
	params, err := fetchParams(harg, "table_name")
	if err != nil {
		return errorRet(badStat, err, "after fetchParams")
	}
	cmd := fmt.Sprintf("drop table %s", params["table_name"])
	log.Debugf("cmd = %s", cmd)
	result, err := db.handle.Exec(cmd)
	if err != nil {
		return errorRet(badStat, err, "after Exec")
	}
	log.Debugf("result = %s", result)
	return apiHandlerRet{http.StatusOK, nil}
}

// describeDbFieldHandler handles GET requests on /db/_schema/{table_name} .
func describeDbFieldHandler(harg *apiHandlerArg) apiHandlerRet {
	return notImplemented()
}

// ----- misc support functions

// tablesQuery is the guts of getDbTablesHandler().
// it's easier to test with an argument.
func tablesQuery(harg *apiHandlerArg,
		tabname string,
		fieldname string) apiHandlerRet {
	// the "tables" table is our convention, not maintained by sqlite.

	idlist := []interface{}{}
	qstring := fmt.Sprintf("select %s from %s;", fieldname, tabname)
	result, err := runQuery(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err, "after runQuery")
	}
	ret, err := convTableNames(result)
	if err != nil {
		return errorRet(badStat, err, "after convTableNames")
	}

	return apiHandlerRet{http.StatusOK, TablesResponse{Names: ret}}
}

// errorRet() is called by apiHandler routines to pass back the code/data
// pair appropriate to the given code and error object.
// optionally logs a debug message along with the code and error.
func errorRet(code int, err error, dmsg string) apiHandlerRet {
	if dmsg != "" {
		log.Debugf("errorRet %d [%s], %s", code, err, dmsg)
	}
	return apiHandlerRet{code, ErrorResponse{code, err.Error()}}
}

// notImpemented() returns the code/data pair for an apiHandler
// that is not implemented.
func notImplemented() apiHandlerRet {
	return errorRet(http.StatusNotImplemented,
		fmt.Errorf("API not implemented yet"), "")
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

// queryErrorRet() passes thru the first 2 args (ret and err),
// while logging the third argument (dmsg).
func queryErrorRet(ret []*KVRecord,
		err error,
		dmsg string) ([]*KVRecord, error) {
	if dmsg != "" {
		log.Debugf("queryErrorRet [%s], %s", err, dmsg)
	}
	return ret, err
}

// runQuery() does a select query using the given query string.
// the return value is a list of the retrieved records.
func runQuery(db dbType,
		qstring string,
		ivals []interface{}) ([]*KVRecord, error) {
	log.Debugf("query = %s", qstring)
	log.Debugf("ivals = %s", ivals)

	ret := make([]*KVRecord, 0, 1)

	rows, err := db.handle.Query(qstring, ivals...)
	if err != nil {
		return queryErrorRet(ret, err, "failure after Query")
	}

	// ensure rows gets closed at end
	defer rows.Close()	// nolint

	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		return queryErrorRet(ret, err, "failure after Columns")
	}
	log.Debugf("cols = %s", cols)
	ncols := len(cols)

	for rows.Next() {
		vals := mkSQLRow(ncols)
		err = rows.Scan(vals...)
		if err != nil {
			return queryErrorRet(ret, err, "failure after Scan")
		}

		err = convValues(vals)
		if err != nil {
			return queryErrorRet(ret, err, "failure after convValues")
		}
		kvrow := KVRecord{Keys: cols, Values: vals}
		ret = append(ret, &kvrow)
		if len(ret) >= maxRecs { // safety check
			break
		}
	}

	err = rows.Err()
	if err != nil {
		return queryErrorRet(ret, err, "failure after rows.Err")
	}

	return ret, nil
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

// sqlExecResult represents the info of Result returned from sql.Exec().
type sqlExecResult struct {
	lastInsertId idType
	rowsAffected idType
}

// getExecResult() constructs an sqlExecResult from the given
// res argument, presumably obtained from calling sql.Exec.
func getExecResult(res sql.Result) sqlExecResult {
	// fmt.Debugf("result=%s", res)
	lastid, _ := res.LastInsertId()
	log.Debugf("lastid = %d", lastid)

	nrecs, _ := res.RowsAffected()
	log.Debugf("rowsaffected = %d", nrecs)

	return sqlExecResult{idType(lastid), idType(nrecs)}
}

// runInsert() inserts a record whose data is specified by the
// given keys and values.  it returns the id of the inserted record.
func runInsert(db dbType,
		tabname string,
		keys []string,
		values []interface{}) (idType, error) {
	nvalues := len(values)

	keystr := strings.Join(keys, ",")
	placestr := nstring("?", nvalues)

	qstring := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",  // nolint
		tabname, keystr, placestr)

	exres, err := runExec(db, qstring, values)
	return exres.lastInsertId, err
}

// delCommon() is the common part of record deletion APIs.
func delCommon(params map[string]string) apiHandlerRet {
	nc, err := delRecs(db, params)
	if err != nil {
		return errorRet(badStat, err, "after delRec")
	}

	return apiHandlerRet{http.StatusOK, NumChangedResponse{int64(nc)}}
}

// dbErrorRet() returns an error value on behalf of a db caller
// that normally returns an idType/error pair.
func dbErrorRet(err error) (idType, error) {
	return idType(-1), err
}

// delRecs() deletes multiple records, using parameters in the params map.
// it returns the number of records deleted.
func delRecs(db dbType, params map[string]string) (idType, error) {
	idclause, idlist := mkIdClause(params)
	if idclause == "" {
		return dbErrorRet(fmt.Errorf("deletion must specify id or ids"))
	}
	qstring := fmt.Sprintf("DELETE FROM %s %s",		// nolint
		params["table_name"],
		idclause)
	log.Debugf("qstring = %s", qstring)

	exres, err := runExec(db, qstring, idlist)
	if int(exres.rowsAffected) != len(idlist) {
		return dbErrorRet(fmt.Errorf("mismatch in rows affected"))
	}
	return exres.rowsAffected, err
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

// getBodySchemas() returns a json schema from the body of the request.
func getBodySchemas(harg *apiHandlerArg) (TableSchemas, error) {
	jrec := TableSchemas{}
	err := json.NewDecoder(harg.getBody()).Decode(&jrec)
	return jrec, err
}

// getBodyRecord() returns a json record from the body of the given request.
func getBodyRecord(harg *apiHandlerArg) (BodyRecord, error) {
	jrec := BodyRecord{}
        err := json.NewDecoder(harg.getBody()).Decode(&jrec)
	return jrec, err
}

// mkIdClause() takes the API parameters,
// and returns the implied WHERE clause that can be
// plugged in to a query string (for use with Prepare)
// and list of data items (for use with Exec).
// the params examined include id_field, id, and/or ids.
// if id is specified, that is used; otherwise ids.
// if neither id nor ids is specified, the WHERE clause is empty.
func mkIdClause(params map[string]string) (string, []interface{}) { // nolint
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		idlist := []interface{}{aToIdType(id)}
		placestr := "?"
		idclause := fmt.Sprintf("WHERE %s = %s",	// nolint
				id_field, placestr)
		return idclause, idlist
	}

	ids, ok := params["ids"]
	if ok && ids != "" {
		idstrings := strings.Split(ids, ",")
		idlist := idTypesToInterface(idstrings)
		placestr := nstring("?", len(idlist))
		idclause := fmt.Sprintf("WHERE %s in (%s)", id_field, placestr) // nolint
		return idclause, idlist
	}

	// no id and no ids implies everything matches.
	// if this is bad, caller should check.
	return "", []interface{}{}
}

// mkIdClauseUpdate() is like mkIdClause(), but for UPDATE operations.
// the difference is, for now, that the id values are formatted directly
// into the WHERE string, rather than being subbed in by Exec.
func mkIdClauseUpdate(params map[string]string) string {  // nolint
	id_field := params["id_field"]
	id, ok := params["id"]
	if ok {
		return fmt.Sprintf("WHERE %s = %s", id_field, id) // nolint
	}
	ids, ok := params["ids"]
	if ok && ids != "" {
		return fmt.Sprintf("WHERE %s in (%s)", id_field, ids) // nolint
	}

	// no id and no ids implies everything matches.
	// if this is bad, caller should check.
	return ""
}

// updateRec() updates certain fields of a given record or records,
// using parameters in the params map.
// it returns the number of records changed.
func updateRec(db dbType,
		params map[string]string,
		body BodyRecord) (idType, error) {
	dbrec := body.Records[0]
	keylist := dbrec.Keys
	keystr := strings.Join(keylist, ",")
	placestr := nstring("?", len(keylist))
	idclause := mkIdClauseUpdate(params)
	if idclause == "" {
		return dbErrorRet(fmt.Errorf("update must specify id or ids"))
	}

	qstring := fmt.Sprintf("UPDATE %s SET (%s) = (%s) %s",	// nolint
			params["table_name"],
			keystr,
			placestr,
			idclause)

	exres, err := runExec(db, qstring, dbrec.Values)
	return exres.rowsAffected, err
}

// runExec() is common code for database APIs that do
// Prepare followed by Exec followed by getting the exec results.
func runExec(db dbType,
		query string,
		values []interface{}) (sqlExecResult, error) {
	log.Debugf("query = %s", query)
	stmt, err := db.handle.Prepare(query)
	if err != nil {
		return sqlExecResult{}, err
	}
	defer stmt.Close()	// nolint
	result, err := stmt.Exec(values...)
	if err != nil {
		return sqlExecResult{}, err
	}
	return getExecResult(result), nil
}

// mkSelectString() returns the WHERE part of a selection query.
func mkSelectString(params map[string]string) (string, []interface{}) {
	idclause, idlist := mkIdClause(params)

	qstring := fmt.Sprintf("SELECT %s FROM %s %s LIMIT %s OFFSET %s", // nolint
		params["fields"],
		params["table_name"],
		idclause,
		params["limit"],
		params["offset"])

	return qstring, idlist
}

// getCommon() is common code for selection APIs.
func getCommon(params map[string]string) apiHandlerRet {
	qstring, idlist := mkSelectString(params)
	result, err := runQuery(db, qstring, idlist)
	if err != nil {
		return errorRet(badStat, err, "after runQuery")
	}

	if len(result) == 0 {
		return errorRet(badStat, fmt.Errorf("no matching record"), "")
	}

	return apiHandlerRet{http.StatusOK, RecordsResponse{Records:result}}
}

// updateCommon() is common code for update APIs.
func updateCommon(harg *apiHandlerArg, params map[string]string) apiHandlerRet {
	body, err := getBodyRecord(harg)
	if err != nil {
		return errorRet(badStat, err, "after getBodyRecord")
	}
	if len(body.Records) < 1 {
		return errorRet(badStat,
			fmt.Errorf("update: no data records in body"), "")
	}

	ra, err := updateRec(db, params, body)
	if err != nil {
		return errorRet(badStat, err, "after updateRec")
	}
	return apiHandlerRet{http.StatusOK, NumChangedResponse{int64(ra)}}
}

// convTableNames() converts the return format from runQuery()
// into a simple list of names.
func convTableNames(result []*KVRecord) ([]string, error) {
	// convert from query format to simple list of names
	ret := make([]string, len(result))
	for i, row := range result {
		str, ok := (*row).Values[0].(string)
		if !ok {
			return ret, fmt.Errorf("table name conversion error")
		}
		ret[i] = str
	}
	return ret, nil
}

// validateRecords() checks the validity of an array of KVRecord.
// returns an error if any record has an invalid key.
// no validation is done on the values except to check
// that the length matches.
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
	}
	return nil
}

// convValues() converts masked *sql.RawBytes to masked strings.
// the slice is changed in-place.
func convValues(vals []interface{}) error {
	N := len(vals)
	for i := 0; i < N; i++ {
		v := vals[i]
		rbp, ok := v.(*sql.RawBytes)
		if !ok {
			return fmt.Errorf("SQL conversion error")
		}
		vals[i] = string(*rbp)
	}
	return nil
}

func listToMap(strList []string) map[string]int {
	ret := map[string]int{}
	if strList == nil {
		return ret
	}
	for _, s := range strList {
		ret[s] = 1
	}
	return ret
}

func createTablesCommon(tabNames []string, schema TableSchemas) apiHandlerRet {
	retNames := make([]string, 0)
	for i, tab := range schema.Resource {
		log.Debugf("... i = %d, tab = %v", i, tab)
		var guts bytes.Buffer
		sep := ""
		for _, field := range tab.Fields {
			guts.WriteString(sep)
			guts.WriteString(field.Name)
			props := listToMap(field.Properties)
			if props["is_primary_key"] != 0 {
				guts.WriteString(" integer primary key autoincrement")
			} else {
				guts.WriteString(" text not null")
			}
			sep = ", "
		}
		cmd := fmt.Sprintf("create table %s(%s)",
				tabNames[i], guts.String())
		log.Debugf("cmd = %s", cmd)
		result, err := db.handle.Exec(cmd)
		if err != nil {
			return errorRet(badStat, err, "after Exec")
		}
		log.Debugf("result = %v", result)
		retNames = append(retNames, tab.Name)
	}
	return apiHandlerRet{http.StatusCreated, SchemasResponse{retNames}}
}
