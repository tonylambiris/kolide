package jwtmw

import (
	"errors"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kolide/kolide/router/middleware/session"
)

// Auth is gin middleware to check for the jwt session and
// validate that this information is correct.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := session.User(c)

		if user == nil {
			c.AbortWithError(401, errors.New("unable to locate user"))
		}

		_, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(user.Hash))
			return b, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
	}
}
