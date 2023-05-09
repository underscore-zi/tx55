package restapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"tx55/pkg/metalgearonline1/types"
)

func getRankings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	limit := 50
	AllModes := types.GameMode(255)

	var period types.PlayerStatsPeriod
	var page int
	gameMode := AllModes

	period, valid := PeriodParam(c.Param("period")).PlayerStatsPeriod()
	if !valid {
		Error(c, 400, "Invalid period")
		return
	}

	// page is optional so just default to 1
	pageParam, found := c.Params.Get("page")
	if !found {
		pageParam = "1"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		Error(c, 400, "Invalid page")
		return
	}

	if modeParam := c.Query("mode"); modeParam != "" {
		gameMode, valid = GameModeParam(modeParam).GameMode()
		if !valid {
			Error(c, 400, "Invalid mode")
			return
		}
	}

	rankings := make([]RankingEntryJSON, 0, limit)
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
			Error(c, 400, "Cannot get overall ranking from archive")
			return
		}
		if err := db.Raw(query, period, limit, (page-1)*limit).Scan(&rankings).Error; err != nil {
			Error(c, 500, "Database error")
			l.WithError(err).Error("Error getting overall rankings")
			return
		}
	} else {
		query := "SELECT `rank`, user_id, u.display_name, SUM(points) as points FROM player_stats INNER JOIN users u ON u.id = player_stats.user_id WHERE period = ? AND mode = ? GROUP BY user_id ORDER BY `rank` LIMIT ? OFFSET ?"
		if err = db.Raw(query, period, gameMode, limit, (page-1)*limit).Scan(&rankings).Error; err != nil {
			Error(c, 500, "Database error")
			l.WithError(err).Error("Error getting mode rankings")
			return
		}
	}
	success(c, rankings)
}
