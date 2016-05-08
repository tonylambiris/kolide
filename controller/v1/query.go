package v1

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/controller/helpers"
	"github.com/mephux/kolide/model"
	"github.com/mephux/kolide/shared/osquery"
	"github.com/mephux/kolide/shared/querycontrol"
)

// Query all or some nodes
func Query(c *gin.Context) {
	timeout := configuration.Server.QueryTimeout.Duration

	query := osquery.Query{}

	if err := json.Unmarshal(helpers.GetBody(c), &query); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	nodes, err := model.AllNodes()

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	batch := querycontrol.NewBatchQuery(query.Sql, nodes)

	// spew.Dump(query.Timeout.Duration, timeout)

	results := batch.Run(timeout)

	helpers.JsonResp(c, 200, gin.H{
		"results": results,
		"error":   nil,
	})
}
