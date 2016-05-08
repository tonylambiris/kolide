package v1

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/model"
	"github.com/mephux/kolide/shared/httputil"
	"github.com/mephux/kolide/shared/token"
)

type tokenPayload struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Expires int64  `json:"expires_in,omitempty"`
}

// Auth route
func Auth(c *gin.Context) {

	if c.Request.Method == "DELETE" {
		httputil.DelCookie(c.Writer, c.Request,
			"user_session")

		c.Redirect(303, "/")
		return
	}

	in := &tokenPayload{}

	err := c.Bind(in)

	if err != nil {
		log.Error(err)
		c.Redirect(303, "/")
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	log.Info("Login Request: ", email)

	user, err := model.UserLogin(email, password)

	if err != nil {
		c.Redirect(303, "/")

		// helpers.JsonResp(c, 404, gin.H{
		// "token": nil,
		// "error": errors.New("unknown user"),
		// })

		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	token := token.New(token.SessToken, user.Email, user)

	tokenstr, err := token.SignExpires(user.Hash, exp)

	if err != nil {
		log.Error(err)

		log.Errorf("cannot create token for %s. %s", user.Email, err.Error())
		c.Redirect(303, "/login?error=internal_error")
		return
	}

	httputil.SetCookie(c.Writer, c.Request,
		"user_session", tokenstr)

	// helpers.JsonResp(c, 200, gin.H{
	// "token": tokenstr,
	// "error": nil,
	// })

	c.Redirect(303, "/")
}
