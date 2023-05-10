package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
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

func whoAmI(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user_id").(uint)

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	db.First(&user, uid)
	success(c, toUserJSON(&user))
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

func getUserGames(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIdParam := c.Param("user_id")
	limit := 50

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		Error(c, 400, "Invalid user id")
		return
	}

	pageParam, found := c.Params.Get("page")
	if !found {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		Error(c, 400, "Invalid page")
		return
	}

	// I didn't want to dive into writing a query here but Gorm doesn't support nesting Joins
	// so using preload this would give us three queries, or I could split it into 2 (player entries, and a game JOIN game_options)
	// or write the fairly simple query myself and get it done in one

	var gamesPlayed []GamePlayedJSON
	query := "SELECT p.game_id, go.name as game_name, go.has_password as game_has_password, go.user_id as game_host_id, p.created_at, p.deleted_at, p.was_kicked, p.score as points, p.kills, p.deaths FROM game_players p JOIN games g ON p.game_id = g.id JOIN game_options go ON g.game_options_id = go.id WHERE p.user_id = ? ORDER BY p.updated_at DESC LIMIT ? OFFSET ?"
	if err := db.Raw(query, userId, limit, (page-1)*limit).Scan(&gamesPlayed).Error; err != nil {
		Error(c, 500, "Database error")
		l.WithError(err).Error("Error getting user's games")
		return
	}

	for i := 0; i < len(gamesPlayed); i++ {
		gamesPlayed[i].GameName = types.BytesToString([]byte(gamesPlayed[i].GameName))
		gamesPlayed[i].WasHost = gamesPlayed[i].GameHostID == uint(userId)
	}
	success(c, gamesPlayed)
}

func getUserOptions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Param("user_id")

	var options models.PlayerSettings
	if err := db.First(&options, "user_id = ?", userID).Error; err != nil {
		Error(c, 404, "User not found")
	} else {
		success(c, toUserSettingsJSON(options))
	}
}
