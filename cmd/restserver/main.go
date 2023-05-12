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

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "tx55/cmd/restserver/docs"
	"tx55/pkg/configurations"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi"
	"tx55/pkg/restapi/admin"
	_ "tx55/pkg/restapi/admin"
	"tx55/pkg/restapi/crons"
	_ "tx55/pkg/restapi/user"
)

var l = logrus.StandardLogger()

func migrate(config configurations.RestAPI) {
	db, err := config.Database.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Error("Unable to open database")
		return
	}
	if err = db.AutoMigrate(models.All...); err != nil {
		l.WithError(err).Error("Unable to migrate database")
		return
	}

	admindb, err := config.AdminDatabase.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Error("Unable to open admin database")
		return
	}
	if err = admindb.AutoMigrate(admin.AllModels...); err != nil {
		l.WithError(err).Error("Unable to migrate admin database")
		return
	}
}

// @title           Metal Gear Online 1 API
// @version         0.1
// @description     API for accessing MGO1 game and user data
// @host      https://tx12.savemgo.com
// @BasePath  /api/v1
func main() {
	configFile := flag.String("config", "", "Path to config file")
	shouldMigrate := flag.Bool("migrate", false, "Run database migrations")
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

	if *shouldMigrate {
		l.Info("Running database migrations")
		migrate(config)
		l.Info("Finished running database migrations")
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

	server.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := server.Engine.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		panic(err)
	}
}
