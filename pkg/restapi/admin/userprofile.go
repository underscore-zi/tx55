package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/profile", UpdateProfile, ArgsUpdateProfile{}, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/emblem", UpdateEmblem, ArgsUpdateEmblem{}, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/user/:userid/connections", ListUserIPs, nil, []ConnectionInfoJSON{})
}

type ArgsUpdateProfile struct {
	UserID      uint
	DisplayName string
	Password    string
}

func UpdateProfile(c *gin.Context) {
	if !CheckPrivilege(c, PrivUpdateProfiles) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

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
	if !CheckPrivilege(c, PrivUpdateProfiles) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
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

type ConnectionInfoJSON struct {
	Remote string
	Local  string
}

func ListUserIPs(c *gin.Context) {
	if !CheckPrivilege(c, PrivSearchByIP) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

	uidParam := c.Param("userid")
	if uidParam == "" {
		restapi.Error(c, 400, "missing userid")
		return
	}

	uid, err := strconv.Atoi(uidParam)
	if err != nil {
		restapi.Error(c, 400, "invalid userid")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var connections []models.Connection
	if tx := db.Where("user_id = ?", uid).Find(&connections); tx.Error != nil {
		c.MustGet("log").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": uid,
			"admin_id":    c.MustGet("admin_id"),
		}).Error("failed to list user connections")
		restapi.Error(c, 500, tx.Error.Error())
		return
	}

	var out []ConnectionInfoJSON
	for _, conn := range connections {
		remoteAddr := conn.RemoteAddr
		localAddr := conn.LocalAddr

		if !CheckPrivilege(c, PrivFullIPs) {
			remoteAddr = remoteAddr[:strings.LastIndex(remoteAddr, ".")+1] + "xxx"
			localAddr = localAddr[:strings.LastIndex(localAddr, ".")+1] + "xxx"
		}

		out = append(out, ConnectionInfoJSON{
			Remote: fmt.Sprintf("%s:%d", remoteAddr, conn.RemotePort),
			Local:  fmt.Sprintf("%s:%d", localAddr, conn.LocalPort),
		})
	}

	restapi.Success(c, out)
}
