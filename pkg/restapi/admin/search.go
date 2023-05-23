package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/search/:name/:page", SearchPlayer)
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/ip/:ip/:page", SearchIP)
}

// SearchPlayer godoc
// @Summary      Search User by Display Name and Username
// @Description  Find users by portions of their display name or username. The Display Name in this result will be in the form of `Username (Display Name)`.
// @Tags         AdminLogin
// @Produce      json
// @Param        name  path  string  true  "Name"
// @Param        page  path  string  false  "Page"
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.UserJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/search/{name}/{page} [get]
func SearchPlayer(c *gin.Context) {
	var limit = 50
	l := c.MustGet("logger").(*logrus.Logger)
	db := c.MustGet("db").(*gorm.DB)

	name := c.Param("name")
	if name == "" {
		restapi.Error(c, 400, "Missing name")
		return
	}

	page := restapi.ParamAsInt(c, "page", 1)
	var users []models.User
	if err := db.Where("CAST(display_name as CHAR(20) CHARACTER SET latin1) LIKE ?", "%"+name+"%").Limit(limit).Offset((page - 1) * limit).Find(&users).Error; err != nil {
		l.WithError(err).WithFields(logrus.Fields{
			"page":  page,
			"limit": limit,
			"name":  name,
		}).Error("Error searching for users")
		restapi.Error(c, 500, "Database error")
		return
	}

	var out []restapi.UserJSON
	for _, user := range users {
		user.DisplayName = []byte(fmt.Sprintf("%s (%s)", user.Username, user.DisplayName))
		out = append(out, *restapi.ToUserJSON(&user))
	}

	restapi.Success(c, out)
}

type ConnectionJSON struct {
	ID         uint             `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	LastUsed   time.Time        `json:"last_used"`
	User       restapi.UserJSON `json:"user"`
	RemoteAddr string           `json:"remote_addr"`
	LocalAddr  string           `json:"local_addr"`
}

// SearchIP godoc
// @Summary      Search User by Display Name and Username
// @Description  Find users by portions of their display name or username. The Display Name in this result will be in the form of `Username (Display Name)`.
// @Tags         AdminLogin
// @Produce      json
// @Param        name  path  string  true  "IP"
// @Param        page  path  string  false  "Page"
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.UserJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      403  {object}  restapi.ResponseJSON{data=string}
// @Router       /admin/ip/{ip}/{page} [get]
func SearchIP(c *gin.Context) {
	if !CheckPrivilege(c, PrivSearchByIP) {
		restapi.Error(c, 403, "Insufficient privileges")
		return
	}
	fullIps := CheckPrivilege(c, PrivFullIPs)

	var limit = 50
	l := c.MustGet("logger").(*logrus.Logger)
	db := c.MustGet("db").(*gorm.DB)

	page := restapi.ParamAsInt(c, "page", 1)
	ip := c.Param("ip")
	if ip == "" {
		restapi.Error(c, 400, "Missing IP")
		return
	}

	// Make sure this is just IP characters in the search
	ip = regexp.MustCompile(`^[0-9.]+$`).FindString(ip)

	// If they are not allowed to see full IPs, don't allow searches by them
	if !fullIps {
		if strings.Count(ip, ".") == 3 {
			ip = ip[:strings.LastIndex(ip, ".")]
		}
	}

	var connections []models.Connection
	q := db.Model(&models.Connection{}).Joins("User").Where("remote_addr LIKE ? or local_addr LIKE ?", ip+"%", ip+"%")
	q = q.Order("updated_at desc").Limit(limit).Offset((page - 1) * limit)
	if err := q.Find(&connections).Error; err != nil {
		l.WithError(err).WithFields(logrus.Fields{
			"page":  page,
			"limit": limit,
			"ip":    ip,
		}).Error("Error searching for connections")
		restapi.Error(c, 500, "Database error")
		return
	}

	var out []ConnectionJSON
	for _, connection := range connections {
		if !fullIps {
			connection.RemoteAddr = connection.RemoteAddr[:strings.LastIndex(connection.RemoteAddr, ".")] + ".xxx"
			connection.LocalAddr = connection.LocalAddr[:strings.LastIndex(connection.LocalAddr, ".")] + ".xxx"
		}

		connection.User.DisplayName = []byte(fmt.Sprintf("%s (%s)", connection.User.Username, connection.User.DisplayName))
		out = append(out, ConnectionJSON{
			ID:         connection.ID,
			CreatedAt:  connection.CreatedAt,
			LastUsed:   connection.UpdatedAt,
			User:       *restapi.ToUserJSON(&connection.User),
			RemoteAddr: fmt.Sprintf("%s:%d", connection.RemoteAddr, connection.RemotePort),
			LocalAddr:  fmt.Sprintf("%s:%d", connection.LocalAddr, connection.LocalPort),
		})
	}

	restapi.Success(c, out)
}
