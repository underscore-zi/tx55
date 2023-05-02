package restapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

func getUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.First(&user, c.Param("user_id")).Error; err != nil {
		Error(c, 404, "User not found")
	} else {
		success(c, toUserJSON(&user))
	}
}

func getUserStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var stats []models.PlayerStats
	if err := db.Find(&stats, "user_id = ?", c.Param("user_id")).Error; err != nil {
		success(c, []PlayerStatsJSON{})
	} else {
		out := make([]PlayerStatsJSON, len(stats))
		for i, stat := range stats {
			out[i] = toPlayerStatsJSON(stat)
		}
		success(c, out)
	}
}
