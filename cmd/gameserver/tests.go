package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/testclient"
	"tx55/pkg/metalgearonline1/types"
)

/* This is not intended to be any sort of test suite, this is really just code I'd write to trigger events
 * to test the handlers. I use hard-coded tests and values. This is just to test the handlers, but I figure the code
 * could be useful to have around.
 */

const (
	IdPS2   = 1
	IdPCSX2 = 4
)

func RunTests(c *testclient.TestClient, db *gorm.DB) {
	defer func() { _ = c.Close() }()
	TestCreateGameAndStats(c, db)
}

func createGame(c *testclient.TestClient, db *gorm.DB) *models.Game {
	if err := c.LoginWithSession("eb82595d-d6b8-4269-80ee-5e9e25703aa9"); err != nil {
		panic(err)
	}
	if err := c.ReportConnectionInfo("192.168.1.111", 5730, 5730); err != nil {
		panic(err)
	}
	if err := c.HostCreateGame("Test Game", "Test Description", ""); err != nil {
		panic(err)
	}
	if err := c.HostPlayerJoinTeam(IdPS2, types.Team(0)); err != nil {
		panic(err)
	}

	var newGame models.Game
	db.Model(&models.Game{}).Where("user_id = ?", IdPS2).Order("id desc").Preload(clause.Associations).First(&newGame)
	if newGame.ID == 0 {
		panic("game not created")
	}
	_ = newGame.Refresh(db)

	if len(newGame.Players) != 1 {
		panic("host was not added as first player")
	}

	return &newGame
}

func playerJoin(c *testclient.TestClient, db *gorm.DB, newGame *models.Game) {
	if err := c.HostPlayerJoin(IdPCSX2); err != nil {
		panic(err)
	}

	_ = newGame.Refresh(db)
	if len(newGame.Players) != 2 {
		panic("New player not added")
	}

	if err := c.HostPlayerJoinTeam(IdPCSX2, types.Team(0)); err != nil {
		panic(err)
	}
	_ = newGame.Refresh(db)
	if newGame.Players[1].Team != types.Team(0) {
		panic("Player not added to team")
	}
}

func TestCreateGameAndStats(c *testclient.TestClient, db *gorm.DB) {
	userId := IdPCSX2
	newGame := createGame(c, db)
	playerJoin(c, db, newGame)

	stats := types.HostReportedStats{
		Kills:      1,
		KillStreak: 1,
	}
	if err := c.HostPlayerStats(types.UserID(userId), stats); err != nil {
		panic(err)
	}
}

//goland:noinspection GoUnusedExportedFunction
func TestCreateGameAndDisconnect(c *testclient.TestClient, db *gorm.DB) {
	newGame := createGame(c, db)
	playerJoin(c, db, newGame)

	if err := c.HostPlayerKicked(IdPCSX2); err != nil {
		panic(err)
	}
	_ = newGame.Refresh(db)
	if len(newGame.Players) != 1 {
		panic("Player not kicked")
	}

	_ = c.Close()

	// Give the server time to see the closed connection as a disconnect
	time.Sleep(time.Second * 2)

	var deletedGame models.Game
	tx := db.Unscoped().Where("id = ?", newGame.ID).Preload(clause.Associations).First(&deletedGame)
	if tx.Error != nil {
		panic(tx.Error)
	}
	if !deletedGame.DeletedAt.Valid {
		panic("Game not deleted on disconnect")
	}
}
