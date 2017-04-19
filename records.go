package apidCRUD

// ----- types for parameter record and response structures

// DeleteResponse is the response data for the API deleteDbRecord and others.
type DeleteResponse struct {
	NumChanged int
}

// ErrorResponse is the response data for API errors.
type ErrorResponse struct {
	Code int
	Message string
}

// jsonRecord is the body data and/or response data for multiple APIs.
type jsonRecord struct {
	Resource []dbRecord
}

// dbRecord is an element type used in multiple APIs.
type dbRecord struct {
	Keys []string
	Values []string
}

// GetRecordResponse is the type for multiple get*Record* APIs.
type GetRecordResponse struct {
	Ids []string
	Record interface{}
}

// RecordIds is the type returned by createDbRecords .
type RecordIds struct {
	Ids []int
}

// TablesResponse is the type returned by getDbTables.
type TablesResponse struct {
	Resource []string
}
