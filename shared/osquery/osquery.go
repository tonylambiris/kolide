package osquery

import (
	"encoding/json"
	"time"
)

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

// OsqueryTimestamp is a type alias used for parsing osquery's timestamp format
type OsqueryTimestamp time.Time

// UnmarshalJSON is the custom parsing logic for the osquery timestamp format
func (ot *OsqueryTimestamp) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	var t time.Time
	t, err = time.Parse("Mon Jan 2 15:04:05 2006 MST", string(b))
	*ot = OsqueryTimestamp(t)
	return
}

// LogResultType is the request json result set
type LogResultType struct {
	Name      string      `json:"name"`
	Id        string      `json:"hostIdentifier"`
	UnixTime  string      `json:"unixTime"`
	Timestamp OsqueryTimestamp   `json:"calendarTime"`
	Results   interface{} `json:"columns"`
	Action    string      `json:"action"`
}

// LogStatusReq is the request for log output
type LogStatusReq struct {
	Data []LogStatusType `json:"data"`
	Key  string          `json:"node_key"`
	Type string          `json:"log_type"`
}

// LogHeader holds the common data for this type
type LogHeader struct {
	Key  string `json:"node_key"`
	Type string `json:"log_type"`
}

// LogReq wrapper
type LogReq struct {
	LogHeader
	Data json.RawMessage `json:"data"`
}
