package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/news/create", createNews)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/news/:id/update", updateNews)
	restapi.Register(restapi.AuthLevelAdmin, "POST", "/admin/policy/update", updatePolicy)
}

type ArgsUpdateNews struct {
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

func createNews(c *gin.Context) {
	if !CheckPrivilege(c, PrivManageNews) {
		restapi.Error(c, 403, "Insufficient privileges")
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var args ArgsUpdateNews
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, "Invalid arguments")
		return
	}

	entry := models.News{
		Topic: args.Topic,
		Body:  args.Body,
	}
	if err := db.Create(&entry).Error; err != nil {
		l.WithError(err).Error("Error creating news")
		restapi.Error(c, 500, "Error creating news")
	} else {
		restapi.Success(c, restapi.NewsJSON{
			ID:        entry.ID,
			CreatedAt: entry.CreatedAt,
			Title:     entry.Topic,
		})
	}
}

func updateNews(c *gin.Context) {
	if !CheckPrivilege(c, PrivManageNews) {
		restapi.Error(c, 403, "Insufficient privileges")
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)
	id := restapi.ParamAsUint(c, "id", 0)

	var args ArgsUpdateNews
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, "Invalid arguments")
		return
	}

	if args.Topic == "" && args.Body == "" {
		//Deleting
		if err := db.Delete(&models.News{}, id).Error; err != nil {
			l.WithError(err).Error("Error deleting news")
			restapi.Error(c, 500, "Error deleting news")
		} else {
			restapi.Success(c, nil)
		}
	} else {
		updates := map[string]interface{}{
			"topic": args.Topic,
			"body":  args.Body,
		}
		if err := db.Model(&models.News{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			l.WithError(err).Error("Error updating news")
			restapi.Error(c, 500, "Error updating news")
		} else {
			restapi.Success(c, nil)
		}
	}
}

type ArgsUpdatePolicy struct {
	Body string `json:"body"`
}

func updatePolicy(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var args ArgsUpdatePolicy
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, "Invalid arguments")
		return
	}

	var entry models.News
	if err := db.Find(&entry, "topic = 'policy'").Error; err == gorm.ErrRecordNotFound {
		entry.Topic = "policy"
		entry.Body = args.Body
	} else if err != nil {
		l.WithError(err).Error("Error getting policy")
		restapi.Error(c, 500, "Error getting policy")
		return
	} else {
		entry.Body = args.Body
	}
	if err := db.Save(&entry).Error; err != nil {
		l.WithError(err).Error("Error saving policy")
		restapi.Error(c, 500, "Error saving policy")
		return
	}
	restapi.Success(c, nil)
}
