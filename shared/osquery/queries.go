package osquery

import "github.com/mephux/kolide/config"

// Query holds user session requested queries
type Query struct {
	All     bool            `json:"all" form:"all" binding:"required"`
	Nodes   []string        `json:"nodes" form:"nodes"`
	Sql     string          `json:"sql" form:"sql" binding:"required"`
	Timeout config.Duration `json:"timeout"`
}
