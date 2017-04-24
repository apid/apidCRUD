package apidCRUD

// ----- types for parameter record and response structures

// NumChangedResponse is the response data for the API deleteDbRecord and others.
type NumChangedResponse struct {
	NumChanged int64
}

// ErrorResponse is the response data for API errors.
type ErrorResponse struct {
	Code int
	Message string
}

type BodyRecord struct {
	Records []KVRecord
}

// KVRecord represents record data in requests, used in multiple APIs.
type KVRecord struct {
	Keys []string
	Values []string
}

// RecordsResponse is the type for multiple get*Record* APIs.
type RecordsResponse struct {
	Records interface{}
}

// IdsResponse is the type returned by createDbRecords .
type IdsResponse struct {
	Ids []int64
}

// TablesResponse is the type returned by getDbTables.
type TablesResponse struct {
	Names []string
}
