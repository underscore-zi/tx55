package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequireAPIKey just requires that the X-API-TOKEN header is set to some value.
// The particular value does not matter as the intent is purely to determine if
// the call is has the ability to set custom headers on a request.
func RequireAPIKey(c *gin.Context) {
	apiKey := c.GetHeader("X-API-TOKEN")
	if apiKey == "" {
		Error(c, 401, "unauthorized")
		return
	}
	c.Next()
}

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

func ProvideContextVar(name string, val any) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(name, val)
		c.Next()
	}
}
