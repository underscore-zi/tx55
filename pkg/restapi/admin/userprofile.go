package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/profile", UpdateProfile, ArgsUpdateProfile{}, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/emblem", UpdateEmblem, ArgsUpdateEmblem{}, restapi.UserJSON{})
}

type ArgsUpdateProfile struct {
	UserID      uint
	DisplayName string
	Password    string
}

func UpdateProfile(c *gin.Context) {
	RequirePrivilege(c, PrivUpdateProfiles)

	var args ArgsUpdateProfile
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	var user models.User
	var err error

	user.ID = args.UserID
	if args.DisplayName != "" {
		user.DisplayName, err = iso8859.EncodeAsBytes(args.DisplayName)
		if err != nil {
			restapi.Error(c, 400, "display name contains invalid characters")
			return
		}
	}

	if args.Password != "" {
		if len(args.Password) < 3 {
			restapi.Error(c, 400, "password too short (minimum 3 characters)")
			return
		}
		user.Password, err = iso8859.Encode(args.Password)
		if err != nil {
			restapi.Error(c, 400, "password contains invalid characters")
			return
		}
	}

	db := c.MustGet("db").(*gorm.DB)
	if tx := db.Model(&user).Updates(user); tx.Error != nil {
		c.MustGet("log").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": user.ID,
			"admin_id":    c.MustGet("admin_id"),
		}).Error("failed to update game user")
		restapi.Error(c, 500, tx.Error.Error())
		return
	}

	db.Model(&user).First(&user)
	restapi.Success(c, restapi.ToUserJSON(&user))
}

type ArgsUpdateEmblem struct {
	UserID     uint   `json:"userid" binding:"required"`
	HasEmblem  bool   `json:"has_emblem" binding:"required"`
	EmblemText string `json:"emblem_text" binding:"required"`
}

func UpdateEmblem(c *gin.Context) {
	RequirePrivilege(c, PrivUpdateProfiles)
	var args ArgsUpdateEmblem

	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	emblemText, err := iso8859.EncodeAsBytes(args.EmblemText)
	if err != nil {
		restapi.Error(c, 400, "emblem text contains invalid characters")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if tx := db.Model(&models.User{
		Model: gorm.Model{ID: args.UserID},
	}).Updates(map[string]interface{}{
		"has_emblem":  args.HasEmblem,
		"emblem_text": emblemText,
	}); tx.Error != nil {
		c.MustGet("log").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": args.UserID,
			"admin_id":    c.MustGet("admin_id"),
		}).Error("failed to update game user emblem")
	}
	restapi.Success(c, nil)
}
