package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GameLoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user_id")
	if user == nil {
		Error(c, 401, "unauthorized")
		return
	} else {
		c.Next()
	}
}
func AdminLoginRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("admin_id")
	if user == nil {
		Error(c, 401, "unauthorized")
		return
	} else {
		c.Next()
	}
}
