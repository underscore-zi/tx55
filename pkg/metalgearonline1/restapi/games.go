package restapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
)

func getGame(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var game models.Game

	// This one is unscoped to allow for querying of games that have already ended
	q := db.Unscoped().Joins("GameOptions")
	q = q.Preload("Players").Preload("Players.User")
	if err := q.First(&game, c.Param("game_id")).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			Error(c, 404, "Game not found")
		default:
			l.WithError(err).Error("Error getting game")
			Error(c, 500, "Error getting game")
		}
	} else {
		success(c, toGameJSON(game))
	}
}

func getGamesList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var games []models.Game
	q := db.Joins("GameOptions").Preload("Players")
	if err := q.Find(&games).Error; err != nil {
		l.WithError(err).Error("Error getting games list")
		Error(c, 500, "Error getting games list")
	} else {
		out := make([]GameJSON, len(games))
		for i, game := range games {
			out[i] = toGameJSON(game)
		}
		success(c, out)
	}

}
