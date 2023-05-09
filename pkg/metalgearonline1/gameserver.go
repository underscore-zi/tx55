package metalgearonline1

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
	"tx55/pkg/konamiserver"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
)

type GameServer struct {
	Db            *gorm.DB
	Sessions      map[string]*session.Session
	sessionLock   sync.Mutex
	LobbyID       types.LobbyID
	KonamiServer  *konamiserver.Server
	Log           logrus.FieldLogger
	LastBanUpdate time.Time
	BannedIPs     map[string]bool
}

type Config struct {
	Address string
	Db      *gorm.DB
	LobbyID types.LobbyID
	Log     logrus.FieldLogger
}

func (gs *GameServer) IsBannedIP(ip string) bool {
	if gs.LastBanUpdate.Before(time.Now().Add(-5 * time.Minute)) {
		query := "SELECT remote_addr FROM connections WHERE user_id IN (SELECT user_id FROM bans WHERE expires_at > NOW() AND type=?) GROUP BY remote_addr"

		var ips []string
		gs.Db.Raw(query, models.IPBan).Find(&ips)
		gs.LastBanUpdate = time.Now()

		gs.BannedIPs = make(map[string]bool)
		for _, i := range ips {
			gs.BannedIPs[i] = true
		}
	}

	if _, found := gs.BannedIPs[ip]; found {
		return true
	}
	return false
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

func (gs *GameServer) ClientFactory(_ string) konamiserver.GameClient {
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
