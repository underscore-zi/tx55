package metalgearonline1

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"tx55/pkg/konamiserver"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
)

type GameServer struct {
	Db           *gorm.DB
	Sessions     map[string]*session.Session
	sessionLock  sync.Mutex
	LobbyID      types.LobbyID
	KonamiServer *konamiserver.Server
	Log          logrus.FieldLogger
}

type Config struct {
	Address string
	Db      *gorm.DB
	LobbyID types.LobbyID
	Log     logrus.FieldLogger
}

func (gs *GameServer) DeleteSession(id string) {
	gs.sessionLock.Lock()
	defer gs.sessionLock.Unlock()

	delete(gs.Sessions, id)

	// Doing the DB update here while we hold the lock to work double duty and prevent races
	gs.Db.Model(&models.Lobby{ID: uint32(gs.LobbyID)}).Update("players", gorm.Expr("players - 1"))
}

func (gs *GameServer) NewSession() *session.Session {
	gs.sessionLock.Lock()
	defer gs.sessionLock.Unlock()

	newSession := &session.Session{
		ID:      uuid.New().String(),
		DB:      gs.Db,
		LobbyID: gs.LobbyID,
		Log:     gs.Log,
	}
	gs.Sessions[newSession.ID] = newSession

	// Doing the DB update here while we hold the lock to work double duty and prevent races
	gs.Db.Model(&models.Lobby{ID: uint32(gs.LobbyID)}).Update("players", gorm.Expr("players + 1"))
	return newSession
}

func (gs *GameServer) Start() error {
	// Clear out any old games that would be disconnected at this point
	var games []models.Game
	gs.Db.Where("lobby_id = ?", gs.LobbyID).Find(&games)
	var ids []uint
	for _, game := range games {
		ids = append(ids, game.ID)
	}

	gs.Db.Where("game_id IN ?", ids).Delete(&models.GamePlayers{})
	gs.Db.Where("lobby_id = ?", gs.LobbyID).Delete(&models.Game{})
	gs.Db.Model(&models.Lobby{ID: uint32(gs.LobbyID)}).Update("players", 0)

	return gs.KonamiServer.Start()
}

func (gs *GameServer) ClientFactory(id string) konamiserver.GameClient {
	return &GameClient{
		Server:  gs,
		Session: gs.NewSession(),
	}
}

func NewGameServer(cfg Config) *GameServer {
	gs := &GameServer{
		Sessions: make(map[string]*session.Session),
		Db:       cfg.Db,
		LobbyID:  cfg.LobbyID,
		Log:      cfg.Log,
	}

	gs.KonamiServer = konamiserver.NewServer(konamiserver.Config{
		Address:       cfg.Address,
		Key:           types.XORKEY,
		ClientFactory: gs.ClientFactory,
		Log:           cfg.Log,
	})

	return gs
}
