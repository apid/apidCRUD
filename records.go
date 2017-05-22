package apidCRUD

// ----- types for parameter record and response structures

// NumChangedResponse is the response data for API deleteDbRecord and others.
type NumChangedResponse struct {
	NumChanged int64
	Kind string
}

// ErrorResponse is the response data for API errors.
type ErrorResponse struct {
	Code int
	Message string
	Kind string
}

// KVRecord represents record data in requests, used in multiple APIs.
type KVRecord struct {
	Keys []string
	Values []interface{}
}

// BodyRecord is the body data for APIs that create or update database records.
type BodyRecord struct {
	Records []KVRecord
}

type KVResponse struct {
	Keys []string
	Values []interface{}
	Kind string
	Self string
}

// RecordsResponse is the type for multiple get*Record* APIs.
type RecordsResponse struct {
	Records []*KVResponse
	Kind string
}

// IdsResponse is the type returned by createDbRecords .
type IdsResponse struct {
	Ids []int64
	Kind string
}

// TablesResponse is the type returned by getDbTables.
type TablesResponse struct {
	Names []string
	Kind string
}

// FieldSchema is the type used to specify a field in a table.
type FieldSchema struct {
	Name string
	Properties []string
}

// TableSchema is the type used to describe one table to be created.
type TableSchema struct {
	Fields []FieldSchema
}

// SchemaResponse is the response format for table creation.
type SchemaResponse struct {
	Schema string
	Kind string
	Self string
}
