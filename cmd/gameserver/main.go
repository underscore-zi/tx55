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
		l.Info("Enabling Hook: Transfer sessions from old database to new")
		GormDb = db
		server.KonamiServer.AddHook(uint16(types.ClientLogin), konamiserver.HookBefore, hookLogin)
	}

	// This is a development hook, rewrites the remote_addr in the host's connection info
	// since when working inside a LAN the address can be incorrect/inaccessible
	if v, found := os.LookupEnv("FORCED_HOST_REMOTE_ADDR"); found {
		l.Info("Enabling Hook: Rewrite all host's remote_addr with: " + v)
		server.KonamiServer.AddHook(uint16(types.ServerHostInfo), konamiserver.HookOutputPacket, hookConnectionInfo)
	}

	l.WithField("address", cfg.Address).Info("Starting server")
	if endpoint, found := os.LookupEnv("EVENTS_ENDPOINT"); !found {
		l.Info("EVENTS_ENDPOINT not set, events will not be broadcast to external service")
	} else {
		l.WithField("events_endpoint", endpoint).Info("Sending events to external service")
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}
