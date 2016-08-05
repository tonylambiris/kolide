package session

import (
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kolide/kolide/model"
	"github.com/kolide/kolide/shared/token"
)

// User retuns the session user context
func User(c *gin.Context) *model.User {
	v, ok := c.Get("user")
	if !ok {
		return nil
	}
	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

// Token will return the session token context
func Token(c *gin.Context) *token.Token {
	v, ok := c.Get("token")
	if !ok {
		return nil
	}
	u, ok := v.(*token.Token)
	if !ok {
		return nil
	}
	return u
}

// SetUser will extract the user data from the session
// and add the user to the gin context. This function will also
// do some limited validation.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *model.User

		t, err := token.ParseRequest(c.Request, func(t *token.Token) (string, error) {
			var err error
			user, err = model.GetUserByEmail(t.Text)
			return user.Hash, err
		})

		if err == nil {
			c.Set("user", user)

			// if this is a session token (ie not the API token)
			// this means the user is accessing with a web browser,
			// so we should implement CSRF protection measures.
			if t.Kind == token.SessToken {
				err = token.CheckCsrf(c.Request, func(t *token.Token) (string, error) {
					return user.Hash, nil
				})
				// if csrf token validation fails, exit immediately
				// with a not authorized error.
				if err != nil {
					log.Error("CSRF Fail")
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}
		}

		c.Next()
	}
}

// MustAdmin will force the current user session to be admin
// and report is value respectfully
func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)

		switch {
		case user == nil:
			c.AbortWithStatus(http.StatusUnauthorized)
			// c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
		case user.Admin == false:
			c.AbortWithStatus(http.StatusForbidden)
			// c.HTML(http.StatusForbidden, "401.html", gin.H{})
		default:
			c.Next()
		}

	}
}

// MustUser will force the current user session to exsit
// and be valid
func MustUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User(c)

		switch {
		case user == nil:
			c.AbortWithStatus(http.StatusUnauthorized)
			// c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
		default:
			c.Next()
		}

	}
}
