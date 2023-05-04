package restapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

func getNewsList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var news []models.News
	if err := db.Find(&news).Error; err != nil {
		l.WithError(err).Error("Error getting news list")
		Error(c, 500, "Error getting news list")
	} else {
		out := make([]NewsJSON, len(news))
		for i, n := range news {
			out[i] = NewsJSON{
				CreatedAt: n.CreatedAt,
				Title:     n.Topic,
				Content:   n.Body,
			}
		}
		success(c, out)
	}
}
