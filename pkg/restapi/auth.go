package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
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

func Login(c *gin.Context) {
	var args struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var user models.User
	db := c.MustGet("db").(*gorm.DB)

	err := c.BindJSON(&args)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	db.Model(&user).First(&user, "username = ?", args.Username)

	if user.ID == 0 || !user.CheckRawPassword([]byte(args.Password)) {
		Error(c, 400, "invalid credentials")
		return
	}

	session := sessions.Default(c)
	session.Clear()
	session.Set("user_id", user.ID)
	_ = session.Save()

	success(c, toUserJSON(&user))
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	success(c, nil)
}
