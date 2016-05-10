package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/controller/helpers"
	"github.com/mephux/kolide/model"
)

// Node route
func Node(c *gin.Context) {
	key := c.Param("key")

	node, err := model.FindNodeByNodeKey(key)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"node":  node,
		"error": nil,
	})
}

// UpdateNode route
func UpdateNode(c *gin.Context) {
	key := c.Param("key")

	node, err := model.FindNodeByNodeKey(key)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	name := c.PostForm("name")
	// categoryId := c.PostForm("category_id")

	if len(name) > 0 {
		node.Name = name
	}

	if err := node.Update(); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"node":  node,
		"error": nil,
	})
}

// DeleteNode route
func DeleteNode(c *gin.Context) {
	key := c.Param("key")

	q, err := model.FindNodeByNodeKey(key)

	if err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	if err := q.Delete(); err != nil {
		helpers.JsonError(c, 500, err)
		return
	}

	helpers.JsonResp(c, 200, gin.H{
		"node":  q,
		"error": nil,
	})
}
