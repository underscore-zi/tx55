package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/restapi"
)

// RequirePrivilege should be the first function calls by any handler that requires a particular privilege
func RequirePrivilege(c *gin.Context, p Privilege) {
	user := FetchUser(c)
	if !user.HasPrivilege(p) {
		restapi.Error(c, 401, "unauthorized")
		return
	}
	c.Next()
}

// FetchUser will grab the user+role from the database once and cache it in the request context, returning the cached
// object on subsequent calls
func FetchUser(c *gin.Context) *User {
	if u, exists := c.Get("admin_user"); exists {
		return u.(*User)
	}

	adminDB := c.MustGet("adminDB").(*gorm.DB)
	session := sessions.Default(c)

	var u User
	u.ID = session.Get("admin_id").(uint)
	adminDB.Model(&u).Joins("Role").First(&u)
	c.Set("admin_user", &u)
	return &u
}
