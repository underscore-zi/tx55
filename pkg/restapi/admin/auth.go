package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "POST", "/admin/login", Login)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/change_password", ChangePassword)
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/whoami", WhoAmI)
}

type ArgsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Admin Login
// @Description  Login to an administrative session
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        body     body  ArgsLogin  true  "Account credentials"
// @Success      200  {object}  restapi.ResponseJSON{data=User}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/login [post]
// @Security ApiKeyAuth
func Login(c *gin.Context) {
	var args ArgsLogin
	var user User
	adminDB := c.MustGet("adminDB").(*gorm.DB)

	err := c.BindJSON(&args)
	if err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	if err := adminDB.Model(&user).Joins("Role").First(&user, "username = ?", args.Username).Error; err != nil {
		logrus.WithError(err).Error("failed to fetch admin user")
		restapi.Error(c, 500, "database error")
		return
	}

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

// ChangePassword godoc
// @Summary      Change Password
// @Description  Change the password for the currently logged in administrative user
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        body     body  ArgsChangePassword  true  "User profile data"
// @Success      200  {object}  restapi.ResponseJSON{}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/change_password [post]
// @Security ApiKeyAuth
func ChangePassword(c *gin.Context) {
	if !CheckPrivilege(c, PrivNone) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

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
	adminUser.Password = args.NewPassword

	if tx := adminDB.Model(&adminUser).Updates(adminUser); tx.Error != nil {
		log := c.MustGet("logger").(logrus.FieldLogger)
		log.WithError(tx.Error).Error("failed to update admin user password")
		restapi.Error(c, 500, "Database Error")
		return
	}

	restapi.Success(c, nil)
}

// WhoAmI godoc
// @Summary      Profile of Current Admin User
// @Description  Get the profile and role of the current administrative user
// @Tags         AdminLogin
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{date=User}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/whoami [post]
// @Security ApiKeyAuth
func WhoAmI(c *gin.Context) {
	a := FetchUser(c)
	restapi.Success(c, a)
}
