package session

import (
	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kolide/kolide/shared/osquery"
)

var (
	queryChannel = make(map[int64]chan *osquery.Query)
)

// AddQuery retuns the session qeries context
func AddQuery(c *gin.Context, query *osquery.Query) {
	user := User(c)

	if user == nil {
		log.Error("session query error: unable to find user")
	}

	if _, ok := queryChannel[user.Id]; !ok {
		queryChannel[user.Id] = make(chan *osquery.Query)
	}

	queryChannel[user.Id] <- query
}

// Queries to send
func Queries(c *gin.Context) chan *osquery.Query {
	user := User(c)

	if user == nil {
		return nil
	}

	if _, ok := queryChannel[user.Id]; !ok {
		log.Error("CREATE NEW")
		queryChannel[user.Id] = make(chan *osquery.Query)
	}

	return queryChannel[user.Id]
}
