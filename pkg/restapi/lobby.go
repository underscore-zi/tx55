package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

func init() {
	Register(AuthLevelNone, "GET", "/lobby/list", getLobbyList, nil, []LobbyJSON{})
}

func getLobbyList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)

	var lobbies []models.Lobby
	if err := db.Where("type = ?", types.LobbyTypeGame).Find(&lobbies).Error; err != nil {
		l.WithError(err).Error("Error getting lobbies list")
		Error(c, 500, "Error getting lobbies list")
	} else {
		out := make([]LobbyJSON, len(lobbies))
		for i, lobby := range lobbies {
			out[i] = toLobbyJSON(lobby)
		}
		success(c, out)
	}
}
