package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/config"
	"github.com/mephux/kolide/controller/v1"
	"github.com/mephux/kolide/router/middleware/session"
)

// Register v1 route handlers
func Register(e *gin.Engine, configuration *config.Config) {

	e.POST("/authorize", v1.Auth)
	e.DELETE("/authorize", v1.Auth)

	g := e.Group("/api/v1")

	v1.Init(configuration)

	g.POST("/osquery/enroll", v1.OSQEnroll)
	g.POST("/osquery/config", v1.OSQConfig)
	g.POST("/osquery/read", v1.OSQRead)
	g.POST("/osquery/write", v1.OSQWrite)
	g.POST("/osquery/log", v1.OSQLog)

	// everything after this will require auth
	g.Use(session.MustUser())

	g.POST("/query", v1.Query)

	// saved queries
	g.GET("/saved-queries", v1.SavedQueries)
	g.DELETE("/saved-queries/:id", v1.DeleteSavedQuery)
	g.POST("/saved-queries", v1.CreateSavedQuery)

	// nodes
	g.DELETE("/nodes/:key", v1.DeleteNode)
	// g.GET("/nodes", v1.Nodes)
	// g.PUT("/nodes/:id", v1.UpdateNode)

	// websocket
	g.GET("/websocket", v1.Websocket)
}
