package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id", getUser, nil, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/stats", getUserStats, nil, []restapi.PlayerStatsJSON{})
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/games", getUserGames, nil, []restapi.GameJSON{})
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/games/:page", getUserGames, nil, []restapi.GameJSON{})
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/settings", getUserSettings, nil, restapi.UserSettingsJSON{})

	restapi.Register(restapi.AuthLevelUser, "GET", "/whoami", whoAmI, nil, restapi.UserJSON{})
	restapi.Register(restapi.AuthLevelUser, "POST", "/user/profile", UpdateUserProfile, ArgsUpdateProfile{}, nil)
	restapi.Register(restapi.AuthLevelUser, "POST", "/user/settings", UpdateUserSettings, restapi.UserSettingsJSON{}, nil)
}

func getUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.First(&user, c.Param("user_id")).Error; err != nil {
		restapi.Error(c, 404, "User not found")
	} else {
		restapi.Success(c, restapi.ToUserJSON(&user))
	}
}

func whoAmI(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user_id").(uint)

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	db.First(&user, uid)
	restapi.Success(c, restapi.ToUserJSON(&user))
}

func getUserStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var stats []models.PlayerStats
	if err := db.Find(&stats, "user_id = ?", c.Param("user_id")).Error; err != nil {
		restapi.Success(c, []restapi.PlayerStatsJSON{})
	} else {
		out := make([]restapi.PlayerStatsJSON, len(stats))
		for i, stat := range stats {
			out[i] = restapi.ToPlayerStatsJSON(stat)
		}
		restapi.Success(c, out)
	}
}

func getUserGames(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)
	userIdParam := c.Param("user_id")
	limit := 50

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		restapi.Error(c, 400, "Invalid user id")
		return
	}

	pageParam, found := c.Params.Get("page")
	if !found {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		restapi.Error(c, 400, "Invalid page")
		return
	}

	// I didn't want to dive into writing a query here but Gorm doesn't support nesting Joins
	// so using preload this would give us three queries, or I could split it into 2 (player entries, and a game JOIN game_options)
	// or write the fairly simple query myself and get it done in one

	var gamesPlayed []restapi.GamePlayedJSON
	query := "SELECT p.game_id, go.name as game_name, go.has_password as game_has_password, go.user_id as game_host_id, p.created_at, p.deleted_at, p.was_kicked, p.score as points, p.kills, p.deaths FROM game_players p JOIN games g ON p.game_id = g.id JOIN game_options go ON g.game_options_id = go.id WHERE p.user_id = ? ORDER BY p.updated_at DESC LIMIT ? OFFSET ?"
	if err := db.Raw(query, userId, limit, (page-1)*limit).Scan(&gamesPlayed).Error; err != nil {
		restapi.Error(c, 500, "Database error")
		l.WithError(err).Error("Error getting user's games")
		return
	}

	for i := 0; i < len(gamesPlayed); i++ {
		gamesPlayed[i].GameName = types.BytesToString([]byte(gamesPlayed[i].GameName))
		gamesPlayed[i].WasHost = gamesPlayed[i].GameHostID == uint(userId)
	}
	restapi.Success(c, gamesPlayed)
}

func getUserSettings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Param("user_id")

	var options models.PlayerSettings
	if err := db.First(&options, "user_id = ?", userID).Error; err != nil {
		restapi.Error(c, 404, "User not found")
	} else {
		restapi.Success(c, restapi.ToUserSettingsJSON(options))
	}
}

type ArgsUpdateProfile struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
}

func UpdateUserProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	session := sessions.Default(c)

	var args ArgsUpdateProfile
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	var user models.User
	user.ID = session.Get("user_id").(uint)
	if args.DisplayName != "" {
		if bs, err := iso8859.EncodeAsBytes(args.DisplayName); err != nil {
			restapi.Error(c, 400, "Display name contains characters that can't be typed in-game")
			return
		} else {
			if len(bs) > 16 {
				restapi.Error(c, 400, "Display name too long")
				return
			}
			user.DisplayName = bs
		}
	}

	if args.Password != "" {
		newPassword, err := iso8859.Encode(args.Password)
		if err != nil {
			restapi.Error(c, 400, "New Password contains characters that can't be typed in-game")
			return
		}

		if len(newPassword) < 3 {
			// Game will silently fail on this
			restapi.Error(c, 400, "Password too short")
			return
		}

		if args.OldPassword != "" {
			restapi.Error(c, 400, "Missing old password")
			return
		}

		oldPassword, err := iso8859.EncodeAsBytes(args.OldPassword)
		if err != nil {
			restapi.Error(c, 400, "Old Password contains characters that can't be typed in-game")
			return
		}

		if !user.CheckRawPassword(oldPassword) {
			restapi.Error(c, 400, "Incorrect password")
			return
		}

		user.Password = newPassword
	}

	if tx := db.Updates(&user); tx.RowsAffected != 1 {
		if tx.Error != nil {
			log := c.MustGet("logger").(*logrus.Logger)
			log.WithError(tx.Error).Error("Error updating user")
		}
		restapi.Error(c, 500, "Database error")
		return
	}
	restapi.Success(c, nil)
}

func stringToOrientation(s string) types.SwitchOrientation {
	switch s {
	case types.CameraOrientation.String():
		return types.CameraOrientation
	case types.PlayerOrientation.String():
		return types.PlayerOrientation
	default:
		return types.PlayerOrientation
	}
}

func stringToSwitch(s string) types.GearSwitchMode {
	switch s {
	case types.GearSwitchToggle.String():
		return types.GearSwitchToggle
	case types.GearSwitchFlashback.String():
		return types.GearSwitchFlashback
	case types.GearSwitchCycle.String():
		return types.GearSwitchCycle
	default:
		return types.GearSwitchToggle
	}
}

func UpdateUserSettings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	session := sessions.Default(c)
	uid := session.Get("user_id").(uint)

	var args restapi.UserSettingsJSON
	if err := c.ShouldBindJSON(&args); err != nil {
		restapi.Error(c, 400, err.Error())
		return
	}

	var settings, oldSettings models.PlayerSettings
	if err := db.First(&oldSettings, "user_id = ?", uid).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			restapi.Error(c, 500, "Database error")
		}
		// When there are no settings, we can just create them
	}

	settings.ID = oldSettings.ID
	settings.UserID = uid
	settings.ShowNameTags = args.ShowNameTags
	settings.SwitchSpeed = byte(args.SwitchSpeed) - 1
	settings.FPVVertical = args.FPVVertical
	settings.FPVHorizontal = args.FPVHorizontal
	settings.FPVSwitchOrientation = bool(stringToOrientation(args.FPVSwitchOrientation))
	settings.TPVVertical = args.TPVVertical
	settings.TPVHorizontal = args.TPVHorizontal
	settings.TPVChase = args.TPVChase
	settings.FPVRotationSpeed = byte(args.FPVRotationSpeed) - 1
	settings.EquipmentSwitchStyle = byte(stringToSwitch(args.EquipmentSwitchStyle))
	settings.TPVRotationSpeed = byte(args.TPVRotationSpeed) - 1
	settings.WeaponSwitchStyle = byte(stringToSwitch(args.WeaponSwitchStyle))
	settings.FKey0, _ = iso8859.EncodeAsBytes(args.FKeys[0])
	settings.FKey1, _ = iso8859.EncodeAsBytes(args.FKeys[1])
	settings.FKey2, _ = iso8859.EncodeAsBytes(args.FKeys[2])
	settings.FKey3, _ = iso8859.EncodeAsBytes(args.FKeys[3])
	settings.FKey4, _ = iso8859.EncodeAsBytes(args.FKeys[4])
	settings.FKey5, _ = iso8859.EncodeAsBytes(args.FKeys[5])
	settings.FKey6, _ = iso8859.EncodeAsBytes(args.FKeys[6])
	settings.FKey7, _ = iso8859.EncodeAsBytes(args.FKeys[7])
	settings.FKey8, _ = iso8859.EncodeAsBytes(args.FKeys[8])
	settings.FKey9, _ = iso8859.EncodeAsBytes(args.FKeys[9])
	settings.FKey10, _ = iso8859.EncodeAsBytes(args.FKeys[10])
	settings.FKey11, _ = iso8859.EncodeAsBytes(args.FKeys[11])

	if tx := db.Save(&settings); tx.RowsAffected != 1 {
		if tx.Error != nil {
			log := c.MustGet("logger").(*logrus.Logger)
			log.WithError(tx.Error).Error("Error updating user settings")
		}
		restapi.Error(c, 500, "Database error")
		return
	}
	restapi.Success(c, nil)
}
