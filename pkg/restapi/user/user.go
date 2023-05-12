package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/iso8859"
)

func init() {
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id", getUser)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/stats", getUserStats)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/games", getUserGames)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/games/:page", getUserGames)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/:user_id/settings", getUserSettings)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/search/:name", SearchByName)
	restapi.Register(restapi.AuthLevelNone, "GET", "/user/search/:name/:page", SearchByName)

	restapi.Register(restapi.AuthLevelUser, "GET", "/whoami", whoAmI)
	restapi.Register(restapi.AuthLevelUser, "POST", "/user/profile", UpdateUserProfile)
	restapi.Register(restapi.AuthLevelUser, "POST", "/user/settings", UpdateUserSettings)
}

// getUser godoc
// @Summary      Game user's profile
// @Description  Provides the public profile of a game user (display name, last seen, rankings, etc.)
// @Tags         GameUser
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserJSON}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/{user_id} [get]
func getUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	uid := restapi.ParamAsUint(c, "user_id", 0)
	if uid == 0 {
		restapi.Error(c, 400, "Invalid user id")
		return
	}

	var user models.User
	if err := db.First(&user, uid).Error; err != nil {
		restapi.Error(c, 404, "User not found")
	} else {
		restapi.Success(c, restapi.ToUserJSON(&user))
	}
}

// whoAmI godoc
// @Summary      Public profile of the currently logged in game user
// @Description  Provides the public profile of a the currently logged in user
// @Tags         GameUserLogin
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserJSON}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Router       /whoami [get]
func whoAmI(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("user_id").(uint)

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	db.First(&user, uid)
	restapi.Success(c, restapi.ToUserJSON(&user))
}

// getUserStats godoc
// @Summary      Retrieve stats for a user
// @Description  Retrieves all stats generated by a particular user. Stats are split-up by game mode, map, and period. So a user may not have stats generated for certain combinations yet. But if there is a All-Time stat entry, there will be a matching weekly one, even if it's empty.
// @Tags         GameUser
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.PlayerStatsJSON}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/{user_id}/stats [get]
func getUserStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	uid := restapi.ParamAsUint(c, "user_id", 0)
	if uid == 0 {
		restapi.Error(c, 400, "Invalid user id")
		return
	}
	var stats []models.PlayerStats
	if err := db.Find(&stats, "user_id = ?", uid).Error; err != nil {
		restapi.Success(c, []restapi.PlayerStatsJSON{})
	} else {
		out := make([]restapi.PlayerStatsJSON, len(stats))
		for i, stat := range stats {
			out[i] = restapi.ToPlayerStatsJSON(stat)
		}
		restapi.Success(c, out)
	}
}

// getUserGames godoc
// @Summary      Retrieve summary of games a user has played in
// @Description  Retrieves a high-level summary of the games a user has played in. Summary includes points, kills, deaths along with time-stamps.
// @Tags         GameUser
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Param        page  path  int  false  "Page"
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.GamePlayedJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/{user_id}/games/{page} [get]
func getUserGames(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	l := c.MustGet("logger").(*logrus.Logger)
	limit := 50

	userId := restapi.ParamAsUint(c, "user_id", 0)
	if userId == 0 {
		restapi.Error(c, 400, "Invalid user id")
		return
	}

	page := restapi.ParamAsInt(c, "page", 1)

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

// getUserSettings godoc
// @Summary      Retrieve user's in-game settings
// @Description  Retrieves a the in-game settings for a particular user
// @Tags         GameUser
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserSettingsJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      404  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/{user_id}/settings [get]
func getUserSettings(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := restapi.ParamAsUint(c, "user_id", 0)
	if userID == 0 {
		restapi.Error(c, 400, "Invalid user id")
		return
	}

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

// UpdateUserProfile godoc
// @Summary      Update User Profile
// @Description  Updates the profile of the currently logged in user. You can either provide just the new display name or just old and new password.
// @Tags         GameUserLogin
// @Accept       json
// @Produce      json
// @Param 	     body  body  ArgsUpdateProfile  true  "Body"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserSettingsJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/profile [post]
// @Security ApiKeyAuth
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

// UpdateUserSettings godoc
// @Summary      Update User Settings
// @Description  Updates the in-game settings of the currently logged in user. ALl settings must be provided.
// @Tags         GameUserLogin
// @Accept       json
// @Produce      json
// @Param 	     body  body  restapi.UserSettingsJSON  true  "Body"
// @Success      200  {object}  restapi.ResponseJSON{data=restapi.UserSettingsJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Failure      500  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/settings [post]
// @Security ApiKeyAuth
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

// SearchByName godoc
// @Summary      Search User by Display Name
// @Description  Find users by portions of their display name
// @Tags         GameUser
// @Produce      json
// @Param        name  path  string  true  "Name"
// @Param        page  path  string  false  "Page"
// @Success      200  {object}  restapi.ResponseJSON{data=[]restapi.UserJSON{}}
// @Failure      400  {object}  restapi.ResponseJSON{data=string}
// @Router       /user/search/{name}/{page} [get]
func SearchByName(c *gin.Context) {
	var limit = 50
	l := c.MustGet("logger").(*logrus.Logger)
	db := c.MustGet("db").(*gorm.DB)

	name := c.Param("name")
	if name == "" {
		restapi.Error(c, 400, "Missing name")
		return
	}

	name, err := iso8859.Encode(name)
	if err != nil {
		restapi.Error(c, 400, "Name contains invalid characters")
		return
	}

	page := restapi.ParamAsInt(c, "page", 1)
	var users []models.User
	if err := db.Debug().Where("CAST(display_name as CHAR(20) CHARACTER SET latin1) LIKE ?", "%"+name+"%").Limit(limit).Offset((page - 1) * limit).Find(&users).Error; err != nil {
		l.WithError(err).WithFields(logrus.Fields{
			"page":  page,
			"limit": limit,
			"name":  name,
		}).Error("Error searching for users")
		restapi.Error(c, 500, "Database error")
		return
	}

	var out []restapi.UserJSON
	for _, user := range users {
		out = append(out, *restapi.ToUserJSON(&user))
	}

	restapi.Success(c, out)
}
