package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "GET", "/admin/bans/list", ListBans, nil, []BanJSON{})
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/bans/update", UpdateBans, ArgsUpdateBan{}, nil)
}

type BanJSON struct {
	ID        uint             `json:"id"`
	User      restapi.UserJSON `json:"user"`
	CreatedBy string           `json:"created_by"`
	UpdatedBy string           `json:"updated_by"`
	BanType   string           `json:"ban_type"`
	Reason    string           `json:"reason"`
	ExpiresAt time.Time        `json:"expires_at"`
}

func ListBans(c *gin.Context) {
	if !CheckPrivilege(c, PrivReadBans) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	var bans []models.Ban
	q := db.Joins("User").Where("expires_at > ?", time.Now())
	if err := q.Find(&bans).Error; err != nil {
		c.MustGet("logger").(*logrus.Logger).WithError(err).Error("Error getting bans list")
		restapi.Error(c, 500, "Error getting bans list")
	}

	var out []BanJSON
	for _, ban := range bans {
		b := BanJSON{
			ID:        ban.ID,
			User:      *restapi.ToUserJSON(&ban.User),
			BanType:   ban.Type.String(),
			CreatedBy: ban.CreatedBy,
			UpdatedBy: ban.UpdatedBy,
			Reason:    ban.Reason,
			ExpiresAt: ban.ExpiresAt,
		}
		out = append(out, b)
	}
	restapi.Success(c, out)
}

type ArgsUpdateBan struct {
	BanID     uint      `json:"ban_id"`
	BanType   string    `json:"ban_type" binding:"required"`
	UserID    uint      `json:"user_id" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

func UpdateBans(c *gin.Context) {
	if !CheckPrivilege(c, PrivUpdateBans) {
		restapi.Error(c, 403, "insufficient privileges")
		return
	}
	adminUser := FetchUser(c)
	db := c.MustGet("db").(*gorm.DB)

	var args ArgsUpdateBan
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, "Invalid arguments")
		return
	}

	updatedBan := models.Ban{
		Reason:    args.Reason,
		ExpiresAt: args.ExpiresAt,
		UserID:    args.UserID,
	}

	switch args.BanType {
	case models.IPBan.String():
		updatedBan.Type = models.IPBan
	case models.UserBan.String():
		updatedBan.Type = models.UserBan
	default:
		restapi.Error(c, 400, "Invalid ban type")
	}

	if args.BanID <= 0 {
		updatedBan.ID = 0
		updatedBan.CreatedBy = adminUser.Username
	} else {
		updatedBan.ID = args.BanID
		updatedBan.UpdatedBy = adminUser.Username
	}
	if tx := db.Save(&updatedBan); tx.Error != nil {
		l := c.MustGet("logger").(*logrus.Logger)
		l.WithError(tx.Error).WithFields(logrus.Fields{
			"ban_id":   args.BanID,
			"ban_type": args.BanType,
			"user_id":  args.UserID,
			"reason":   args.Reason,
			"admin_id": adminUser.ID,
		}).Error("Error updating ban")
		restapi.Error(c, 500, "Error updating ban")
	} else {
		restapi.Success(c, nil)
	}
}
