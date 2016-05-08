package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/version"

	log "github.com/Sirupsen/logrus"
)

// JsonResp is a wrapper for the basic api endpoint result
func JsonResp(c *gin.Context, status int, params gin.H) {

	if params["error"] != nil {
		log.Error(params["error"].(error))
		params["error"] = params["error"].(error).Error()
	}

	c.JSON(status, gin.H{
		// app version
		"version": version.Version,

		// api version
		"apiVersion": "1",

		// error:
		// resourceVersion:
		// data: obj
		"context": params,
	})
}

// JsonRaw is a wrapper for the basic api endpoint result
func JsonRaw(c *gin.Context, status int, params interface{}) {
	c.JSON(status, params)
}

// JsonError returns a json error with http status code
func JsonError(c *gin.Context, status int, err error) {
	log.Error(err)

	c.JSON(status, gin.H{
		// app version
		"version": version.Version,

		// api version
		"apiVersion": "1",

		// error:
		// resourceVersion:
		// data: obj
		"context": map[string]string{
			"error": err.Error(),
		},
	})
}
