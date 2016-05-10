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

	var nodes []*model.Node

	var err error

	if query.All {
		nodes, err = model.AllNodes(&model.AllNodeOptions{
			OnlyEnabled: true,
		})

		if err != nil {
			helpers.JsonError(c, 500, err)
			return
		}

	} else if len(query.Nodes) > 0 {

		for _, n := range query.Nodes {
			if node, err := model.FindNodeByNodeKey(n); err != nil {
				continue
			} else {
				if node.Enabled {
					nodes = append(nodes, node)
				}
			}

		}
	}

	batch := querycontrol.NewBatchQuery(query.Sql, nodes)

	results := batch.Run(timeout)

	helpers.JsonResp(c, 200, gin.H{
		"results": results,
		"error":   nil,
	})
}
