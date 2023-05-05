package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"tx55/pkg/configurations"
	"tx55/pkg/konamiserver"
	"tx55/pkg/metalgearonline1"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/testclient"
	"tx55/pkg/metalgearonline1/types"
	// Handlers need to be imported to be registered
	// annoying, but it allows for easy swapping out and testing

	_ "tx55/pkg/metalgearonline1/handlers/auth"
	_ "tx55/pkg/metalgearonline1/handlers/general"
	_ "tx55/pkg/metalgearonline1/handlers/hostgame"
	_ "tx55/pkg/metalgearonline1/handlers/lobby"
)

var l = logrus.StandardLogger()

func main() {
	var err error
	l.SetLevel(logrus.InfoLevel)

	configFile := flag.String("config", "", "Path to config file")
	shouldMigrate := flag.Bool("migrate", false, "Run database migrations")
	withTests := flag.Bool("test", false, "Run developer tests")
	doTrace := flag.Bool("trace", false, "Display full packet traces")
	flag.Parse()

	// Load the configuration

	if *configFile == "" {
		l.Fatal("No config file specified")
		return
	}

	var serverConfig configurations.MetalGearOnline1
	if err = configurations.LoadTOML(*configFile, &serverConfig); err != nil {
		l.WithError(err).Fatal("Error loading config file")
		return
	}

	db, err := serverConfig.Database.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), serverConfig.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Error("Unable to open database")
		return
	}

	if *shouldMigrate {
		if err := db.AutoMigrate(models.All...); err != nil {
			l.WithError(err).Error("Error migrating database")
			return
		}
	}

	cfg := metalgearonline1.Config{
		Address: fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Db:      db,
		LobbyID: types.LobbyID(serverConfig.LobbyID),
		Log:     l,
	}

	server := metalgearonline1.NewGameServer(cfg)
	server.KonamiServer.Debug = *doTrace

	// Temporary Hook to transfer logins from old database to new
	if OriginalDb != nil {
		l.Info("Enabling sessions from old database")
		GormDb = db
		server.KonamiServer.AddHook(uint16(types.ClientLogin), konamiserver.HookBefore, hookLogin)
	}

	/* Example Hook: Hooking the game host info packet containing IP/port info to redirect any game to connect to a
	   controlled IP address. When testing games with a server and clients all in the same network the remote addresses
	   were not correct making this necessary */
	/*
		matchedIp := []byte("55.66.77.88")
		replacedIp := "11.22.33.44\x00"
		server.KonamiServer.AddHook(uint16(types.ServerHostInfo), konamiserver.HookOutputPacket, func(p, req *packet.Packet, out chan packet.Packet) konamiserver.HookResult {
			data := (*p).Data()
			idx := bytes.Index(data, matchedIp)
			if idx >= 0 {
				copy(data[idx:], replacedIp)
				(*p).SetData(data)
			}
			return konamiserver.HookResultContinue
		})
		//*/
	l.WithField("address", cfg.Address).Info("Starting server")
	if endpoint, found := os.LookupEnv("EVENTS_ENDPOINT"); !found {
		l.Info("EVENTS_ENDPOINT not set, events will not be broadcast to external service")
	} else {
		l.WithField("events_endpoint", endpoint).Info("Sending events to external service")
	}

	if !*withTests {
		if err := server.Start(); err != nil {
			panic(err)
		}
	} else {
		go func() {
			if err := server.Start(); err != nil {
				panic(err)
			}
		}()
		l.Info("Starting Tests")
		c := testclient.TestClient{Key: server.KonamiServer.Config.Key}
		if err = c.Connect(cfg.Address); err != nil {
			panic(err)
		}
		RunTests(&c, db)
	}
}
