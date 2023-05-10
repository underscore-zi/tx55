package gameweb

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

type ArgsRegister struct {
	Username    string `form:"name"`
	Password    string `form:"passwd"`
	DisplayName string `form:"pname"`
}

func RegisterAccount(c *gin.Context) {
	var args struct {
		Username    string `form:"name"`
		Password    string `form:"passwd"`
		DisplayName string `form:"pname"`
	}

	if err := c.ShouldBind(&args); err != nil {
		c.String(400, "Invalid arguments")
		return
	}

	var newUser models.User
	newUser.Username = []byte(args.Username)
	newUser.Password = args.Password
	newUser.DisplayName = []byte(args.DisplayName)

	db := c.MustGet("db").(*gorm.DB)
	if tx := db.Create(&newUser); tx.Error != nil {
		log.WithError(tx.Error).Error("Failed to save user")
		c.String(500, "Database error")
		return
	}

	log.WithFields(log.Fields{
		"id":           newUser.ID,
		"username":     newUser.Username,
		"display_name": newUser.DisplayName,
	}).Info("Registered user")

	c.String(200, "0")
}

func DeleteAccount(c *gin.Context) {
	var args struct {
		Username string `form:"name"`
		Password string `form:"passwd"`
	}
	if err := c.ShouldBind(&args); err != nil {
		c.String(400, "Invalid arguments")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if tx := db.Where("username = ?", args.Username).First(&user); tx.Error != nil {
		c.String(404, "User not found")
		return
	}
	if !user.CheckRawPassword([]byte(args.Password)) {
		c.String(404, "User not found")
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.String(500, "Database error")
		return
	}

	log.WithFields(log.Fields{
		"id":       user.ID,
		"username": user.Username,
	}).Info("Deleted user")

	c.String(200, "0")
}

func ChangePassword(c *gin.Context) {
	var args struct {
		Username string `form:"name"`
		Password string `form:"passwd"`
		NewPass  string `form:"pswdnew"`
	}
	if err := c.ShouldBind(&args); err != nil {
		c.String(400, "Invalid arguments")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if tx := db.Where("username = ?", args.Username).First(&user); tx.Error != nil {
		c.String(404, "User not found")
		return
	}
	if !user.CheckRawPassword([]byte(args.Password)) {
		c.String(404, "User not found")
		return
	}

	user.Password = args.NewPass
	if err := db.Save(&user).Error; err != nil {
		c.String(500, "Database error")
		return
	}

	log.WithFields(log.Fields{
		"id":       user.ID,
		"username": user.Username,
	}).Info("Changed password")

	c.String(200, "0")
}
