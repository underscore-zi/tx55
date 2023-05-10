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
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/crons"
)

var l = logrus.StandardLogger()

func main() {
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configFile == "" {
		l.Fatal("No config file specified")
		return
	}

	var config configurations.RestAPI
	if err := configurations.LoadTOML(*configFile, &config); err != nil {
		l.WithError(err).Fatal("Error loading config file")
		return
	}

	db, err := config.Database.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Error("Unable to open database")
		return
	}

	if config.RunCronJobs {
		scheduler := gocron.NewScheduler(time.UTC)
		if err = crons.Schedule(scheduler, db); err != nil {
			l.WithError(err).Error("Unable to schedule crons")
			return
		} else {
			l.Info("Starting scheduler")
			scheduler.StartAsync()
		}
	}

	server, err := restapi.NewServer(config)
	if err != nil {
		l.WithError(err).Error("Unable to create server")
		return
	}

	if err := server.Engine.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		panic(err)
	}
}
