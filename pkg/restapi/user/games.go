package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/games/list", getGamesList)
	restapi.Register(restapi.AuthLevelNone, "GET", "/games/:game_id", getGame)
}

// getGame godoc
// @Summary      Retrieve a Game by ID
// @Description  Retrieves all game settings and complete player list along with their stats overview and timestamps.
// @Tags         Games
// @Produce      json
// @Param        game_id  path  string  true   "Game ID"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.GameJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /games/{game_id} [get]
func getGame(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var game models.Game

	gameId := restapi.ParamAsUint(c, "game_id", 0)
	if gameId == 0 {
		restapi.Error(c, 400, "Invalid game ID")
		return
	}

	// This one is unscoped to allow for querying of games that have already ended
	q := db.Unscoped().Joins("GameOptions")
	q = q.Preload("Players").Preload("Players.User")
	if err := q.First(&game, gameId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			restapi.Error(c, 404, "Game not found")
		default:
			l.WithError(err).Error("Error getting game")
			restapi.Error(c, 500, "Error getting game")
		}
	} else {
		restapi.Success(c, restapi.ToGameJSON(game))
	}
}

// getGame godoc
// @Summary      Retrieve all currently active games
// @Description  Retrieves all active games and current player list. The `user` field of the GamePlayer object will be null.
// @Tags         Games
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.GameJSON{}}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /games/list [get]
func getGamesList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var games []models.Game
	q := db.Joins("GameOptions").Preload("Players")
	if err := q.Find(&games).Error; err != nil {
		l.WithError(err).Error("Error getting games list")
		restapi.Error(c, 500, "Error getting games list")
	} else {
		out := make([]restapi.GameJSON, len(games))
		for i, game := range games {
			out[i] = restapi.ToGameJSON(game)
		}
		restapi.Success(c, out)
	}

}
