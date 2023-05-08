package session

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

type HostSession struct {
	GameID       types.GameID
	Rules        [15]types.GameRules
	CurrentRound byte
	Players      map[types.UserID]bool
	Lock         sync.Mutex
}

type Session struct {
	ID        string
	User      *models.User
	GameState *HostSession
	DB        *gorm.DB
	IP        string
	//ActiveConnection gets filled in with user controlled data when they first connect to a game lobby
	ActiveConnection models.Connection

	isHost  bool
	LobbyID types.LobbyID
	Log     logrus.FieldLogger
}

func (s *Session) IsLoggedIn() bool {
	return s.User != nil && s.User.ID > 0
}

func (s *Session) IsHost() bool {
	return s.GameState != nil && s.GameState.GameID > 0
}

func (s *Session) LogFields() logrus.Fields {
	f := logrus.Fields{
		"id": s.ID,
		"ip": s.IP,
	}
	if s.IsLoggedIn() && s.LobbyID > 0 {
		f["state"] = "in-lobby"
		f["lobby"] = s.LobbyID
		f["user"] = string(s.User.Username)
		f["user_id"] = s.User.ID
	} else {
		f["state"] = "unconnected"
	}

	if s.IsHost() {
		f["state"] = "hosting"
		f["game"] = s.GameState.GameID
	}
	return f
}

// --- State Changes ---

// Login is also where any first-time setup should be done
func (s *Session) Login(user *models.User) {
	s.User = user

	s.DB.Model(user).Updates(map[string]interface{}{
		"updated_at": gorm.Expr("NOW()"),
	})

	if s.User.PlayerSettings.UserID == 0 {
		var settings models.PlayerSettings
		tx := s.DB.Model(&models.PlayerSettings{}).Where("user_id = ?", s.User.ID).First(&settings)
		if tx.RowsAffected == 1 {
			s.User.PlayerSettings = settings
		}
		// We could insert custom defaults here, or just let the game do it
		// in the future it might be fun to set some custom F-keys for the user
	}
}
