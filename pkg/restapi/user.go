package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	Register(AuthLevelNone, "GET", "/user/:user_id", getUser, nil, UserJSON{})
	Register(AuthLevelNone, "GET", "/user/:user_id/stats", getUserStats, nil, []PlayerStatsJSON{})
	Register(AuthLevelNone, "GET", "/user/:user_id/games", getUserGames, nil, []GameJSON{})
	Register(AuthLevelNone, "GET", "/user/:user_id/games/:page", getUserGames, nil, []GameJSON{})
	Register(AuthLevelNone, "GET", "/user/:user_id/settings", getUserSettings, nil, UserSettingsJSON{})

	Register(AuthLevelUser, "GET", "/whoami", whoAmI, nil, UserJSON{})
	Register(AuthLevelUser, "POST", "/user/profile", UpdateUserProfile, ArgsUpdateProfile{}, UserSettingsJSON{})
}

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
	l := c.MustGet("logger").(*logrus.Logger)
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

func getUserSettings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Param("user_id")

	var options models.PlayerSettings
	if err := db.First(&options, "user_id = ?", userID).Error; err != nil {
		Error(c, 404, "User not found")
	} else {
		success(c, toUserSettingsJSON(options))
	}
}

func UpdateUserProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	session := sessions.Default(c)

	var args ArgsUpdateProfile
	if err := c.ShouldBindJSON(&args); err != nil {
		Error(c, 400, err.Error())
		return
	}

	var user models.User
	user.ID = session.Get("user_id").(uint)
	if args.DisplayName != "" {
		if bs, err := iso8859.EncodeAsBytes(args.DisplayName); err != nil {
			Error(c, 400, "Display name contains characters that can't be typed in-game")
			return
		} else {
			if len(bs) > 16 {
				Error(c, 400, "Display name too long")
				return
			}
			user.DisplayName = bs
		}
	}

	if args.Password != "" {
		newPassword, err := iso8859.Encode(args.Password)
		if err != nil {
			Error(c, 400, "New Password contains characters that can't be typed in-game")
			return
		}

		if len(newPassword) < 3 {
			// Game will silently fail on this
			Error(c, 400, "Password too short")
			return
		}

		if args.OldPassword != "" {
			Error(c, 400, "Missing old password")
			return
		}

		oldPassword, err := iso8859.EncodeAsBytes(args.OldPassword)
		if err != nil {
			Error(c, 400, "Old Password contains characters that can't be typed in-game")
			return
		}

		if !user.CheckRawPassword(oldPassword) {
			Error(c, 400, "Incorrect password")
			return
		}

		user.Password = newPassword
	}

	if tx := db.Debug().Updates(&user); tx.RowsAffected != 1 {
		if tx.Error != nil {
			log := c.MustGet("logger").(*logrus.Logger)
			log.WithError(tx.Error).Error("Error updating user")
		}
		Error(c, 500, "Database error")
		return
	}
	success(c, nil)
}
