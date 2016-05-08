package osquery

import "time"

// KeyReq is the normal request with node_key
type KeyReq struct {
	Key     string `json:"node_key"`
	Address string `json:"address"`
}

// EnrollReq is the request from oquery for enrollment
type EnrollReq struct {
	Secret  string `json:"enroll_secret"`
	Key     string `json:"host_identifier"`
	Address string `json:"address"`
}

// Queries is a type for osquery query
type Queries map[string]interface{}

// QueryType for query json objects
type QueryType map[string]string

// WriteReq is the json body for the write request
type WriteReq struct {
	Queries Queries `json:"queries"`
	Key     string  `json:"node_key"`
}

// ReadResp is the response for read requests
type ReadResp struct {
	Queries QueryType `json:"queries"`
	Invalid bool      `json:"node_invalid"`
}

// LogStatusType is the request json log
type LogStatusType struct {
	Severity string `json:"severity"`
	Filename string `json:"filename"`
	Line     string `json:"line"`
	Message  string `json:"message"`
}

// LogResultType is the request json result set
type LogResultType struct {
	Name      string      `json:"name"`
	Id        string      `json:"hostIdentifier"`
	UnixTime  string      `json:"unixTime"`
	Timestamp time.Time   `json:"calendarTime"`
	Results   interface{} `json:"columns"`
	Action    string      `json:"action"`
}

// LogStatusReq is the request for log output
type LogStatusReq struct {
	Data []LogStatusType `json:"data"`
	Key  string          `json:"node_key"`
	Type string          `json:"log_type"`
}

// LogResultReq is the request for json request logs
type LogResultReq struct {
	Data []LogResultType `json:"data"`
	Key  string          `json:"node_key"`
	Type string          `json:"log_type"`
}
