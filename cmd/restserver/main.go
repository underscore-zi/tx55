package main

import (
	"flag"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"tx55/pkg/configurations"
	"tx55/pkg/metalgearonline1/crons"
	"tx55/pkg/metalgearonline1/restapi"
)

var l = logrus.StandardLogger()

func main() {
	port := flag.Int("port", 8888, "Port to listen on")
	noCrons := flag.Bool("no-crons", false, "Disable scheduled tasks")
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		l.Fatal("No config file specified")
		return
	}

	var serverConfig configurations.MetalGearOnline1
	if err := configurations.LoadTOML(*configFile, &serverConfig); err != nil {
		l.WithError(err).Fatal("Error loading config file")
		return
	}

	//serverConfig.Database.LogConfig.Level = "info"

	db, err := serverConfig.Database.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), serverConfig.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Error("Unable to open database")
		return
	}

	if !*noCrons {
		scheduler := gocron.NewScheduler(time.UTC)
		if err = crons.Schedule(scheduler, db); err != nil {
			l.WithError(err).Error("Unable to schedule crons")
			return
		} else {
			l.Info("Starting scheduler")
			scheduler.StartAsync()
		}
	}

	server := restapi.NewServer(db)
	if err := server.Engine.Run(fmt.Sprintf(":%d", *port)); err != nil {
		panic(err)
	}
}
