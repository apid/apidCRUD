package apidCRUD

// ----- types for parameter record and response structures

// NumChangedResponse is the response data for API deleteDbRecord and others.
type NumChangedResponse struct {
	NumChanged int64
}

// ErrorResponse is the response data for API errors.
type ErrorResponse struct {
	Code int
	Message string
}

// BodyRecord is the body data for APIs that create or update database records.
type BodyRecord struct {
	Records []KVRecord
}

// KVRecord represents record data in requests, used in multiple APIs.
type KVRecord struct {
	Keys []string
	Values []interface{}
}

// RecordsResponse is the type for multiple get*Record* APIs.
type RecordsResponse struct {
	Records []*KVRecord
}

// IdsResponse is the type returned by createDbRecords .
type IdsResponse struct {
	Ids []int64
}

// TablesResponse is the type returned by getDbTables.
type TablesResponse struct {
	Names []string
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

// TableSchemas is the type used to specify multiple tables to be created.
type TableSchemas struct {
	Resource []TableSchema
}

// SchemasResponse is the response format for table creation.
type SchemasResponse struct {
	Names []string
}
