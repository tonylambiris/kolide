package v1

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/controller/helpers"
	"github.com/mephux/kolide/model"
	"github.com/mephux/kolide/shared/osquery"
	"github.com/mephux/kolide/shared/querycontrol"

	log "github.com/Sirupsen/logrus"
)

// OSQEnroll allows new osquery nodes to checkin
func OSQEnroll(c *gin.Context) {
	req := osquery.EnrollReq{}

	if err := json.Unmarshal(helpers.GetBody(c), &req); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	// log.Info("[ENROLL REQ] ", req)
	// log.Info("Attempting to validate enroll key")

	if req.Secret != configuration.Server.EnrollSecret {

		helpers.JsonRaw(c, 200, gin.H{
			"node_invalid": true,
		})

		return
	}

	req.Address = strings.Split(c.ClientIP(), ":")[0]
	node, err := model.CreateOrUpdateNode(&req)

	if err != nil {
		log.Errorf("[osquery enroll/%s]: %s", req.Key, err)

		helpers.JsonRaw(c, 200, gin.H{
			"node_invalid": true,
		})

		return
	}

	helpers.JsonRaw(c, 200, gin.H{
		"node_key":     node.Key,
		"node_invalid": false,
	})
}

// OSQRead receives a post request for the node id
// and returns json queries
func OSQRead(c *gin.Context) {
	req := osquery.KeyReq{}

	if err := json.Unmarshal(helpers.GetBody(c), &req); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	node, err := model.FindNodeByRequest(c, &req)

	if err != nil {
		log.Errorf("[osquery read/%s]: %s", req.Key, err)

		helpers.JsonRaw(c, 200, gin.H{
			"node_invalid": true,
		})

		return
	}

	queryResp := querycontrol.Control.PendingQueries(node)

	if queryResp == nil {
		queryResp = &osquery.ReadResp{
			Queries: make(osquery.QueryType),
			Invalid: false,
		}
	}

	helpers.JsonRaw(c, 200, queryResp)
}

// OSQWrite receives json query results and returns
// node enroll validation
func OSQWrite(c *gin.Context) {
	req := osquery.WriteReq{}

	if err := json.Unmarshal(helpers.GetBody(c), &req); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	node, err := model.FindNodeByRequest(c, &osquery.KeyReq{
		Key:     req.Key,
		Address: strings.Split(c.ClientIP(), ":")[0],
	})

	if err != nil {
		log.Errorf("[osquery write/%s]: %s", req.Key, err)

		helpers.JsonRaw(c, 200, gin.H{
			"node_invalid": true,
		})

		return
	}

	querycontrol.Control.AddResponse(node, &req)

	helpers.JsonRaw(c, 200, gin.H{
		"node_invalid": false,
	})
}

// OSQConfig receives a json node id and returns
// a osqueryd configuration json object
func OSQConfig(c *gin.Context) {
	req := osquery.KeyReq{}

	if err := json.Unmarshal(helpers.GetBody(c), &req); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	// log.Info("[CONFIG REQ] NODE ID::", req.Key)

	_, err := model.FindNodeByRequest(c, &req)

	if err != nil {
		log.Errorf("[osquery config/%s]: %s", req.Key, err)

		helpers.JsonRaw(c, 200, gin.H{
			"node_invalid": true,
		})

		return
	}

	dat, err := ioutil.ReadFile("shared/osqueryd-example-config.json")

	if err != nil {
		log.Error(err)
	}

	var osqueryConfig interface{}

	if err := json.Unmarshal(dat, &osqueryConfig); err != nil {
		log.Error(err)
		return
	}

	helpers.JsonRaw(c, 200, osqueryConfig)
}

// OSQLog receives a osquery log json object
func OSQLog(c *gin.Context) {
	req := osquery.LogReq{}

	d := helpers.GetBody(c)

	if err := json.Unmarshal(d, &req); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	if req.Type == "result" {
		var result []osquery.LogResultType

		if err := json.Unmarshal(req.Data, &result); err != nil {
			helpers.JsonError(c, 500, err)
			return
		}

		log.Info(result)

	} else if req.Type == "status" {
		var status []osquery.LogStatusType

		if err := json.Unmarshal(req.Data, &status); err != nil {
			helpers.JsonError(c, 500, err)
			return
		}

		log.Info(status)

	} else {
		log.Error("unknown log type: ", req.Type)
	}

	helpers.JsonRaw(c, 200, gin.H{
		"node_invalid": false,
	})
}
