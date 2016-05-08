package helpers

import "github.com/gin-gonic/gin"

// Error return index.html for react/pushState
func Error(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}
