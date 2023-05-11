package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "POST", "/admin/login", Login, ArgsLogin{}, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/change_password", ChangePassword, ArgsChangePassword{}, nil)
}

type ArgsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var args ArgsLogin
	var user User
	adminDB := c.MustGet("adminDB").(*gorm.DB)

	err := c.BindJSON(&args)
	if err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	adminDB.Model(&user).Joins("Role").First(&user, "username = ?", args.Username)

	if user.ID == 0 || !user.CheckPassword([]byte(args.Password)) {
		restapi.Error(c, 400, "invalid credentials")
		return
	}

	session := sessions.Default(c)
	session.Clear()
	session.Set("admin_id", user.ID)
	_ = session.Save()

	restapi.Success(c, user)
}

type ArgsChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	RequirePrivilege(c, PrivNone)

	var args ArgsChangePassword
	if err := c.BindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	adminUser := FetchUser(c)
	if !adminUser.CheckPassword([]byte(args.OldPassword)) {
		restapi.Error(c, 400, "invalid password")
		return
	}

	adminDB := c.MustGet("adminDB").(*gorm.DB)
	updates := User{
		ID:       adminUser.ID,
		Password: args.NewPassword,
	}

	if tx := adminDB.Save(&updates); tx.Error != nil {
		log := c.MustGet("logger").(logrus.FieldLogger)
		log.WithError(tx.Error).Error("failed to update admin user password")
		restapi.Error(c, 500, "Database Error")
		return
	}

	restapi.Success(c, nil)
}
