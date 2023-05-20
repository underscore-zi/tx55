package gameweb

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"tx55/pkg/metalgearonline1/models"
)

func GetTextFile(c *gin.Context) {
	filename := c.Param("filename")
	if strings.HasPrefix(filename, "policy") {
		db := c.MustGet("db").(*gorm.DB)

		var policy models.News
		if tx := db.Where("topic = ?", "policy").First(&policy); tx.Error == nil {
			c.String(200, policy.Body)
			return
		} else {
			c.String(200, "No policy set")
			return
		}
	}
	c.String(404, "Not found")
}
