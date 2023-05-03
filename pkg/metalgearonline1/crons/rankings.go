package crons

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

var l = logrus.WithField("pkg", "crons")

func ClearWeeklyStats(db *gorm.DB) {
	// Before we clear weekly, we need to update ranks and then award the champion
	UpdateRankings(db)
	AwardChampion(db)

	updates := map[string]interface{}{
		"kills":                0,
		"deaths":               0,
		"stuns":                0,
		"stuns_received":       0,
		"snake_frags":          0,
		"points":               0,
		"unknown1":             0,
		"unknown2":             0,
		"team_kills":           0,
		"team_stuns":           0,
		"rounds_played":        0,
		"rounds_no_death":      0,
		"kerotans_for_win":     0,
		"kerotans_placed":      0,
		"radio_uses":           0,
		"text_chat_uses":       0,
		"cqc_attacks":          0,
		"cqc_attacks_received": 0,
		"head_shots":           0,
		"head_shots_received":  0,
		"team_wins":            0,
		"kills_with_scorpion":  0,
		"kills_with_knife":     0,
		"times_eaten":          0,
		"rolls":                0,
		"infrared_goggle_uses": 0,
		"play_time":            0,
		"unknown3":             0,
	}
	tx := db.Model(&models.PlayerStats{}).Where("period = ?", types.PeriodWeekly).Updates(updates)
	if tx.Error != nil {
		l.WithError(tx.Error).Error("failed to clear weekly stats")
	}
}

func AwardChampion(db *gorm.DB) {
	emblemText := "Champion"
	// First clear the old oldChampion
	db.Model(&models.User{}).Where("emblem_text = ?", emblemText).Updates(map[string]interface{}{
		"has_emblem":  false,
		"emblem_text": "",
	})

	var oldChampion, newChampion models.User
	if err := db.Model(&models.User{}).Where("emblem_text = ?", emblemText).First(&oldChampion).Error; err != nil {
		l.WithError(err).Error("failed to find oldChampion")
	} else {
		db.Model(oldChampion).Updates(map[string]interface{}{
			"has_emblem":  false,
			"emblem_text": "",
		})
	}

	// Now find the new champion
	if err := db.Model(&models.User{}).Order("weekly_rank desc").Where("id != ? AND NOT has_emblem", oldChampion.ID).First(&newChampion).Error; err != nil {
		l.WithError(err).Error("failed to find a new champion")
		return
	}
	db.Model(&newChampion).Updates(map[string]interface{}{
		"has_emblem":  true,
		"emblem_text": emblemText,
	})
}

func UpdateRankings(db *gorm.DB) {
	// The innerQuery is the key query. It queries for every user and generates their rank for each period/mode
	// The selector takes that query down to a single result that is good for the update query

	// Update the player_stats rank entries. These are stored with each period/mode/map entry. But as ranks in-game are
	// only shown on a per period/mode basis that is all we update to reflect
	innerQuery := "SELECT user_id, mode, period, rank() OVER (PARTITION BY mode, period ORDER BY SUM(points) DESC) AS `rank` FROM player_stats GROUP BY user_id, mode, period"
	selectorQuery := "SELECT `rank` from (" + innerQuery + ") as t WHERE t.user_id=player_stats.user_id AND t.mode=player_stats.mode AND t.period=player_stats.period"
	updateQuery := "UPDATE player_stats SET `rank` = (" + selectorQuery + ")"
	if tx := db.Exec(updateQuery); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to update rankings")
	} else if tx.RowsAffected == 0 {
		l.Error("no rows affected")
	}

	// Update the Overall period(0) ranks in the users table
	innerQuery = "SELECT user_id, period, rank() OVER (PARTITION BY period ORDER BY SUM(points) DESC) AS `rank` FROM player_stats GROUP BY user_id, period"
	selectorQuery = "SELECT `rank` from (" + innerQuery + ") as t WHERE t.user_id=users.id AND t.period=0"
	updateQuery = "UPDATE users SET overall_rank = (" + selectorQuery + ")"
	if tx := db.Exec(updateQuery); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to update overall rankings")
	} else if tx.RowsAffected == 0 {
		l.Error("no rows affected")
	}

	// Update the Weekly period(1) ranks in the users table
	selectorQuery = "SELECT `rank` from (" + innerQuery + ") as t WHERE t.user_id=users.id AND t.period=1"
	updateQuery = "UPDATE users SET weekly_rank = (" + selectorQuery + ")"
	if tx := db.Exec(updateQuery); tx.Error != nil {
		l.WithError(tx.Error).Error("failed to update weekly rankings")
	} else if tx.RowsAffected == 0 {
		l.Error("no rows affected")
	}
}
