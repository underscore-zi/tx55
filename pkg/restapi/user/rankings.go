package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/restapi"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/rankings/:period", getRankings, nil, []restapi.RankingEntryJSON{})
	restapi.Register(restapi.AuthLevelNone, "GET", "/rankings/:period/:mode", getRankings, nil, []restapi.RankingEntryJSON{})
}

func getRankings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)
	limit := 50
	AllModes := types.GameMode(255)

	var period types.PlayerStatsPeriod
	gameMode := AllModes

	period, valid := restapi.PeriodParam(c.Param("period")).PlayerStatsPeriod()
	if !valid {
		restapi.Error(c, 400, "Invalid period")
		return
	}

	page := restapi.ParamAsInt(c, "page", 1)

	if modeParam := c.Query("mode"); modeParam != "" {
		gameMode, valid = restapi.GameModeParam(modeParam).GameMode()
		if !valid {
			restapi.Error(c, 400, "Invalid mode")
			return
		}
	}

	rankings := make([]restapi.RankingEntryJSON, 0, limit)
	if gameMode == AllModes {
		var query string
		switch period {
		case types.PeriodAllTime:
			query = `SELECT overall_rank as ` + "`rank`" + `, id as user_id, t.points, display_name FROM users
					INNER JOIN (SELECT user_id, SUM(points) as points FROM player_stats WHERE period = ? GROUP BY user_id) t ON users.id = t.user_id
					WHERE users.overall_rank > 0 ORDER BY users.overall_rank ASC LIMIT ? OFFSET ?`
		case types.PeriodWeekly:
			query = `SELECT weekly_rank as ` + "`rank`" + `, id as user_id, t.points, display_name FROM users
					INNER JOIN (SELECT user_id, SUM(points) as points FROM player_stats WHERE period = ? GROUP BY user_id) t ON users.id = t.user_id
					WHERE users.weekly_rank > 0 ORDER BY users.weekly_rank ASC LIMIT ? OFFSET ?`
		case types.PeriodArchive:
			restapi.Error(c, 400, "Cannot get overall ranking from archive")
			return
		}
		if err := db.Raw(query, period, limit, (page-1)*limit).Scan(&rankings).Error; err != nil {
			restapi.Error(c, 500, "Database error")
			l.WithError(err).Error("Error getting overall rankings")
			return
		}
	} else {
		query := "SELECT `rank`, user_id, u.display_name, SUM(points) as points FROM player_stats INNER JOIN users u ON u.id = player_stats.user_id WHERE period = ? AND mode = ? GROUP BY user_id ORDER BY `rank` LIMIT ? OFFSET ?"
		if err := db.Raw(query, period, gameMode, limit, (page-1)*limit).Scan(&rankings).Error; err != nil {
			restapi.Error(c, 500, "Database error")
			l.WithError(err).Error("Error getting mode rankings")
			return
		}
	}
	restapi.Success(c, rankings)
}
