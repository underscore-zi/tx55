package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

func init() {
	Register(AuthLevelNone, "POST", "/login", Login, ArgsLogin{}, UserJSON{})
	Register(AuthLevelNone, "GET", "/logout", Logout, nil, nil)
}

type ArgsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var args ArgsLogin
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

	Success(c, ToUserJSON(&user))
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	Success(c, nil)
}
