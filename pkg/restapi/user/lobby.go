package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/lobby/list", getLobbyList, nil, []restapi.LobbyJSON{})
}

// getLobbyList godoc
// @Summary      Retrieve lobby list
// @Description  Retrieves all game lobbies and player counts
// @Tags         Games
// @Produce      json
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.LobbyJSON{}}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /lobby/list [get]
func getLobbyList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var lobbies []models.Lobby
	if err := db.Where("type = ?", types.LobbyTypeGame).Find(&lobbies).Error; err != nil {
		l.WithError(err).Error("Error getting lobbies list")
		restapi.Error(c, 500, "Error getting lobbies list")
	} else {
		out := make([]restapi.LobbyJSON, len(lobbies))
		for i, lobby := range lobbies {
			out[i] = restapi.ToLobbyJSON(lobby)
		}
		restapi.Success(c, out)
	}
}
