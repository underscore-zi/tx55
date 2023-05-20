package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/:userid/profile", UpdateProfile)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/user/:userid/emblem", UpdateEmblem)
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/user/:userid/connections", ListUserConnections)
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/user/search_ip/:ip", SearchByIP)
}

type ArgsUpdateProfile struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

// UpdateProfile godoc
// @Summary      Update In-Game Profile
// @Description  Update game user profile data, such as display name and password
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Param        body     body  ArgsUpdateProfile  true  "User profile data"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserJSON}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/user/{user_id}/profile [post]
// @Security ApiKeyAuth
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

	db := c.MustGet("db").(*gorm.DB)
	user.ID = restapi.ParamAsUint(c, "userid", 0)
	if args.Password != "" {
		// If the password is being changed the User object needs to be filled in (username) for hash generation
		if db.Model(&user).First(&user).Error != nil {
			restapi.Error(c, 400, "invalid user id")
			return
		}
	}

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

	if tx := db.Model(&user).Updates(&user); tx.Error != nil {
		c.MustGet("logger").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": user.ID,
			"admin_id":    FetchUserID(c),
		}).Error("failed to update game user")
		restapi.Error(c, 500, tx.Error.Error())
		return
	}

	db.Model(&user).First(&user)
	restapi.Success(c, restapi.ToUserJSON(&user))
}

type ArgsUpdateEmblem struct {
	HasEmblem  bool   `json:"has_emblem"`
	EmblemText string `json:"emblem_text"`
}

// UpdateEmblem godoc
// @Summary      Update Emblem
// @Description  Update game user's emblem status and text
// @Tags         AdminLogin
// @Accept       json
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Param        body     body  ArgsUpdateEmblem  true  "Emblem Data"
// @Success      200  {object}  restapi.ResponseJSON{}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/user/{user_id}/emblem [post]
// @Security ApiKeyAuth
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

	userID := restapi.ParamAsUint(c, "userid", 0)
	if userID == 0 {
		restapi.Error(c, 400, "invalid user id")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if tx := db.Model(&models.User{
		Model: gorm.Model{ID: userID},
	}).Updates(map[string]interface{}{
		"has_emblem":  args.HasEmblem,
		"emblem_text": emblemText,
	}); tx.Error != nil {
		c.MustGet("logger").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": userID,
			"admin_id":    c.MustGet("admin_id"),
		}).Error("failed to update game user emblem")
	}
	restapi.Success(c, nil)
}

type ResponseConnectionsJSON struct {
	FirstUsed time.Time `json:"first_used"`
	LastUsed  time.Time `json:"last_used"`
	Remote    string    `json:"remote"`
	Local     string    `json:"local"`
}

// ListUserConnections godoc
// @Summary      List Connection History
// @Description  Lists all public/private IPs the user has connected to the MGO1 server using
// @Tags         AdminLogin
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  restapi.ResponseJSON{data=[]ResponseConnectionsJSON}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/user/{user_id}/connections [get]
// @Security ApiKeyAuth
func ListUserConnections(c *gin.Context) {
	if !CheckPrivilege(c, PrivSearchByIP) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}

	uid := restapi.ParamAsUint(c, "userid", 0)
	if uid == 0 {
		restapi.Error(c, 400, "invalid userid")
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var connections []models.Connection
	if tx := db.Where("user_id = ?", uid).Order("updated_at desc").Find(&connections); tx.Error != nil {
		c.MustGet("logger").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_user": uid,
			"admin_id":    c.MustGet("admin_id"),
		}).Error("failed to list user connections")
		restapi.Error(c, 500, tx.Error.Error())
		return
	}

	var out []ResponseConnectionsJSON
	for _, conn := range connections {
		remoteAddr := conn.RemoteAddr
		localAddr := conn.LocalAddr

		if !CheckPrivilege(c, PrivFullIPs) {
			remoteAddr = remoteAddr[:strings.LastIndex(remoteAddr, ".")+1] + "xxx"
			localAddr = localAddr[:strings.LastIndex(localAddr, ".")+1] + "xxx"
		}

		out = append(out, ResponseConnectionsJSON{
			FirstUsed: conn.CreatedAt,
			LastUsed:  conn.UpdatedAt,
			Remote:    fmt.Sprintf("%s:%d", remoteAddr, conn.RemotePort),
			Local:     fmt.Sprintf("%s:%d", localAddr, conn.LocalPort),
		})
	}

	restapi.Success(c, out)
}

type ResponseSearchByIPJSON struct {
	User        restapi.UserJSON        `json:"user"`
	Connections ResponseConnectionsJSON `json:"connection"`
}

// SearchByIP godoc
// @Summary      Search for Users by IP history
// @Description  Lists all users who have connected using the given IP. The IP parameter can be a prefix to search for.
// @Tags         AdminLogin
// @Produce      json
// @Param        ip  path  string  true  "IPv4"
// @Success      200  {object}  restapi.ResponseJSON{data=[]ResponseSearchByIPJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/user/search_ip/{ip} [get]
// @Security ApiKeyAuth
func SearchByIP(c *gin.Context) {
	if !CheckPrivilege(c, PrivSearchByIP) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
	canSeeFullIPs := CheckPrivilege(c, PrivFullIPs)

	ip := c.Param("ip")
	if ip == "" {
		restapi.Error(c, 400, "missing ip")
		return
	}

	if !canSeeFullIPs {
		dotCount := strings.Count(ip, ".")
		if dotCount > 3 {
			restapi.Error(c, 400, "invalid IP")
			return
		} else if dotCount == 3 {
			// We have a full-ip, but a user that isn't allowed to see full IPs
			// so we need to limit their search to the first 3 octets
			ip = ip[:strings.LastIndex(ip, ".")+1]
		}
	}
	ip += "%" // Add wildcard to end of IP

	db := c.MustGet("db").(*gorm.DB)
	var connections []models.Connection
	if tx := db.Where("remote_addr LIKE ? or local_addr LIKE ?", ip, ip).Joins("User").Order("updated_at desc").Find(&connections); tx.Error != nil {
		c.MustGet("logger").(logrus.FieldLogger).WithError(tx.Error).WithFields(logrus.Fields{
			"target_ip": ip,
			"admin_id":  FetchUser(c),
		}).Error("failed to search by IP")
		restapi.Error(c, 500, tx.Error.Error())
		return
	}

	var out []ResponseSearchByIPJSON
	for _, conn := range connections {
		if !canSeeFullIPs {
			conn.RemoteAddr = conn.RemoteAddr[:strings.LastIndex(conn.RemoteAddr, ".")+1] + "xxx"
			conn.LocalAddr = conn.LocalAddr[:strings.LastIndex(conn.LocalAddr, ".")+1] + "xxx"
		}

		out = append(out, ResponseSearchByIPJSON{
			User: *restapi.ToUserJSON(&conn.User),
			Connections: ResponseConnectionsJSON{
				FirstUsed: conn.CreatedAt,
				LastUsed:  conn.UpdatedAt,
				Remote:    fmt.Sprintf("%s:%d", conn.RemoteAddr, conn.RemotePort),
				Local:     fmt.Sprintf("%s:%d", conn.LocalAddr, conn.LocalPort),
			},
		})
	}
	restapi.Success(c, out)
}
