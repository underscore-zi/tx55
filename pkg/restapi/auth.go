package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

func init() {
	Register(AuthLevelNone, "POST", "/login", Login)
	Register(AuthLevelNone, "GET", "/logout", Logout)
}

type ArgsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Login to a GameUser account
// @Description  Login to a specific GameUser account using the in-game credentials
// @Tags         GameUserLogin
// @Accept       json
// @Produce      json
// @Param 	     body  body  restapi.ArgsLogin  true  "Body"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserJSON{}{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Router       /login [post]
func Login(c *gin.Context) {
	var args ArgsLogin
	var user models.User
	db := c.MustGet("db").(*gorm.DB)

	err := c.BindJSON(&args)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	db.Model(&user).First(&user, "username LIKE ?", args.Username)

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

// Logout godoc
// @Summary      Logout of a GameUser or Admin session
// @Description  Logout of any currently existing session
// @Tags         GameUserLogin,AdminLogin
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserJSON{}{}}
// @Router       /logout [post]
// @Security ApiKey	Auth
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	Success(c, nil)
}
