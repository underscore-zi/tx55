package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

func GetExternalIP() (ip string, err error) {
	var found bool
	if ip, found = os.LookupEnv("EXTERNAL_IP"); !found {
		err = fmt.Errorf("EXTERNAL_IP environment variable not set")
	}
	return
}

func initGameDB(db *gorm.DB) (err error) {
	Logger.WithField("type", GameDBMigrationType).Info("Initializing database")

	if err = newsTable(db); err != nil {
		return
	}

	if err = lobbiesTable(db); err != nil {

	}

	Logger.WithField("type", GameDBMigrationType).Info("Initialization complete")
	return
}

func newsTable(db *gorm.DB) (err error) {
	Logger.Info("Checking for policy entry in news table")
	var policy models.News
	if err = db.Where("topic = ?", "policy").First(&policy).Error; err == gorm.ErrRecordNotFound {
		Logger.Debug("Creating policy entry in news table")
		policy = models.News{
			Topic: "policy",
			Body:  "Welcome to Metal Gear Online 1!",
		}
		if err = db.Create(&policy).Error; err != nil {
			return
		}
	} else if err != nil {
		return
	} else {
		Logger.Debug("Policy entry already exists in news table")
	}
	return
}

func lobbiesTable(db *gorm.DB) (err error) {
	// Lobbies require atleast three entries to be usable in-game
	// Gate Server, Account Server, and Game Server though only the Game Server is visible to users
	// Defaults to putting the Gate and Account servers on the same ip:port (5731) and game serve ron 5732

	Logger.Info("Checking for gate server in lobbies table")
	var lobby models.Lobby
	if err = db.Where("type = ?", types.LobbyTypeGate).First(&lobby).Error; err == gorm.ErrRecordNotFound {
		Logger.Info("Creating gate server in lobbies table")

		var ip string
		if ip, err = GetExternalIP(); err != nil {
			return
		}

		lobby = models.Lobby{
			Name: "gate-server",
			Type: types.LobbyTypeGate,
			IP:   ip,
			Port: 5731,
		}
		if err = db.Create(&lobby).Error; err != nil {
			return
		}
		Logger.Info("Gate server created")
	} else if err != nil {
		return
	} else {
		Logger.Info("Gate server found")
	}

	Logger.Info("Checking for account server in lobbies table")
	lobby = models.Lobby{}
	if err = db.Where("type = ?", types.LobbyTypeAccount).First(&lobby).Error; err == gorm.ErrRecordNotFound {
		Logger.Info("Creating account server in lobbies table")

		var ip string
		if ip, err = GetExternalIP(); err != nil {
			return
		}

		lobby = models.Lobby{
			Name: "account-server",
			Type: types.LobbyTypeAccount,
			IP:   ip,
			Port: 5731,
		}
		if err = db.Create(&lobby).Error; err != nil {
			return
		}
		Logger.Info("Account server created")
	} else if err != nil {
		return
	} else {
		Logger.Info("Account server found")
	}

	Logger.Info("Checking for atleast one game server in lobbies table")
	lobby = models.Lobby{}
	if err = db.Where("type = ?", types.LobbyTypeGame).First(&lobby).Error; err == gorm.ErrRecordNotFound {
		Logger.Info("Creating game server in lobbies table")

		var ip string
		if ip, err = GetExternalIP(); err != nil {
			return
		}

		lobby = models.Lobby{
			Name: "SNAKE",
			Type: types.LobbyTypeGame,
			IP:   ip,
			Port: 5732,
		}

		if err = db.Create(&lobby).Error; err != nil {
			return
		}
		Logger.Info("Game server created")
	} else if err != nil {
		return
	} else {
		Logger.Info("Game server found")
	}
	return
}
