package controller

import "github.com/gin-gonic/gin"

// Error return index.html for react/pushState
func Error(c *gin.Context) {
	c.HTML(404, "index.html", gin.H{})
}
