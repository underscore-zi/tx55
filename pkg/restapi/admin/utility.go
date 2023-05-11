package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckPrivilege should be the first function calls by any handler that requires a particular privilege
func CheckPrivilege(c *gin.Context, p Privilege) bool {
	user := FetchUser(c)
	return user.HasPrivilege(p)
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
