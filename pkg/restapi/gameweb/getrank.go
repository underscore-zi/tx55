package gameweb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"tx55/pkg/metalgearonline1/types"
)

type ArgsGetRanks struct {
	Term    int `form:"term"` // Period (all-time/weekly)
	Rule    int `form:"rule"`
	Skey    int `form:"skey"`    // 0 for Points based, 1 for VS Rating
	From    int `form:"from"`    // Rank/offset to start from
	Records int `form:"records"` // Number of records to fetch at once
	// Fetch rankings around where this player is ranked
	Pid int `form:"pid"` // PlayerID
}

type RankEntry struct {
	Rank        uint
	UserID      uint
	DisplayName string
	Kills       uint
	Deaths      uint
	Points      uint
	VsRating    uint
}

func PostGetRanks(c *gin.Context) {
	args := ArgsGetRanks{}
	_ = c.ShouldBind(&args)

	db := c.MustGet("db").(*gorm.DB)

	var query string
	var queryArgs []interface{}

	// VS Rating request
	if args.Skey == 1 {
		args.Rule = int(types.ModeOverall)
	}

	if types.GameMode(args.Rule) == types.ModeOverall {
		var rankCol string
		if args.Skey == 0 {
			switch types.PlayerStatsPeriod(args.Term) {
			case types.PeriodAllTime:
				rankCol = "overall_rank"
			case types.PeriodWeekly:
				rankCol = "weekly_rank"
			default:
				rankCol = "overall_rank"
			}
		} else {
			rankCol = "vs_rating_rank"
		}

		if args.Pid > 0 {
			var pidRank uint
			tx := db.Raw(`SELECT `+rankCol+` FROM users WHERE id = ?`, args.Pid).Scan(&pidRank)
			if tx.Error != nil {
				log.WithError(tx.Error).Error("Failed to get player rank")
				c.String(500, "Failed to get pid rank")
				return
			}
			args.From = int(pidRank) - args.Records/2
			if args.From < 0 {
				args.From = 0
			}
		}

		query = fmt.Sprintf(
			`SELECT %s as `+"`rank`"+`, id as user_id, display_name, t.kills, t.deaths, t.points, vs_rating FROM users 
			INNER JOIN (
			  SELECT user_id, SUM(kills) as kills, SUM(deaths) as deaths, SUM(points) as points FROM player_stats
			  WHERE period = ?
			  GROUP BY user_id
			) t ON users.id = t.user_id
			ORDER BY %s ASC LIMIT ? OFFSET ?`, rankCol, rankCol)
		queryArgs = append(queryArgs, args.Term, args.Records, args.From)
	} else {
		query = `SELECT t.rank as ` + "`rank`" + `, id as user_id, display_name, t.kills, t.deaths, t.points, vs_rating FROM users 
                 INNER JOIN (
                   SELECT user_id, ` + "`rank`" + `, SUM(kills) as kills, SUM(deaths) as deaths, SUM(points) as points FROM player_stats
                   WHERE period = ? AND mode = ?
				   GROUP BY user_id
				 ) t ON users.id = t.user_id
				 ORDER BY t.rank ASC LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, args.Term, args.Rule, args.Records, args.From)
	}

	var rows []RankEntry
	tx := db.Raw(query, queryArgs...).Scan(&rows)
	if tx.Error != nil {
		log.WithFields(log.Fields{
			"term":    args.Term,
			"rule":    args.Rule,
			"skey":    args.Skey,
			"from":    args.From,
			"records": args.Records,
			"pid":     args.Pid,
		}).WithError(tx.Error).Error("Failed to get rank query")
		c.String(500, "")
		return
	}

	builder := strings.Builder{}

	// This first number is the number of returned elements
	// but the number after the / and on the second row, I don't know, but seem to work with (1, 0)
	builder.WriteString(fmt.Sprintf("%d/%d\n", len(rows), 1))
	builder.WriteString("0\n")

	// Not sure what the second to last number is in each of these rows
	for _, row := range rows {
		builder.WriteString(fmt.Sprintf("%d,%d,%s,%d,%d,%d,0,%d\n", row.Rank, row.UserID, row.DisplayName, row.Kills, row.Deaths, row.Points, row.VsRating))
	}
	builder.WriteString("\n")
	c.String(200, builder.String())
	return
}
