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
	"tx55/pkg/migrations"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"tx55/cmd/restserver/docs"
	"tx55/pkg/configurations"
	"tx55/pkg/restapi"
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
		l.WithError(err).Fatal("Unable to open database")
	}

	migrations.Logger = l
	if err = migrations.MigrateModels(migrations.GameDBMigrationType, db); err != nil {
		l.Fatal("Failed to migrate GameDB models")
	}

	admindb, err := config.AdminDatabase.Open(&gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", 0), config.Database.LogConfig.LoggerConfig()),
	})
	if err != nil {
		l.WithError(err).Fatal("Unable to open admin database")
	}

	migrations.Logger = l
	if err = migrations.MigrateModels(migrations.AdminDBMigrationType, admindb); err != nil {
		l.Fatal("Failed to migrate AdminDB models")
	}
}

// @title           Metal Gear Online 1 API
// @version         0.1
// @description     API for accessing MGO1 game and user data
// @host      https://tx12.savemgo.com
// @BasePath  /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-TOKEN
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

	if v, found := os.LookupEnv("SWAGGER_HOST"); found {
		docs.SwaggerInfo.Host = v
	}
	server.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := server.Engine.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		panic(err)
	}
}
