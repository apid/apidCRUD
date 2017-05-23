package apidCRUD

// ----- types for parameter record and response structures

// field tags are used to change the case of JSON keys
// created by json.Unmarshal().

// NumChangedResponse is the response data for API deleteDbRecord and others.
type NumChangedResponse struct {
	NumChanged int64 `json:"numChanged"`
	Kind string	`json:"kind"`
}

// ErrorResponse is the response data for API errors.
type ErrorResponse struct {
	Code int	`json:"code"`
	Message string	`json:"message"`
	Kind string	`json:"kind"`
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

// KVResponse represents data records returned from an API call.
type KVResponse struct {
	Keys []string	`json:"keys"`
	Values []interface{} `json:"values"`
	Kind string	`json:"kind"`
	Self string	`json:"self"`
}

// RecordsResponse is the type for multiple get*Record* APIs.
type RecordsResponse struct {
	Records []*KVResponse `json:"records"`
	Kind string	`json:"kind"`
}

// IdsResponse is the type returned by createDbRecords .
type IdsResponse struct {
	Ids []int64	`json:"ids"`
	Kind string	`json:"kind"`
}

// TablesResponse is the type returned by getDbTables.
type TablesResponse struct {
	Names []string	`json:"names"`
	Kind string	`json:"kind"`
	Self string	`json:"self"`
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
	Schema string	`json:"schema"`
	Kind string	`json:"kind"`
	Self string	`json:"self"`
}

// ServiceResponse is the response format for the describeService API.
type ServiceResponse struct {
	Description string `json:"resource"`
	Kind string	`json:"kind"`
	Self string	`json:"self"`
}
