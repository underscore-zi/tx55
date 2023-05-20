package migrations

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/restapi/admin"
)

var Logger = logrus.StandardLogger()

type MigrationType string

const (
	GameDBMigrationType  = "GameDB"
	AdminDBMigrationType = "AdminDB"
)

// MigrateModels just wraps a call to gorm's AutoMigrate with the appropiate models for the database type
func MigrateModels(dbType MigrationType, db *gorm.DB) (err error) {
	switch dbType {
	case GameDBMigrationType:
		Logger.WithField("type", dbType).Info("Migrating models")
		if err = db.AutoMigrate(models.All...); err != nil {
			Logger.WithError(err).WithField("type", dbType).Error("Error migrating models")
			return
		} else {
			Logger.WithField("type", dbType).Info("Migration complete")
		}
		if err = Initialize(dbType, db); err != nil {
			Logger.WithError(err).WithField("type", dbType).Error("Error initializing database")
			return
		}
	case AdminDBMigrationType:
		Logger.WithField("type", dbType).Info("Migrating models")
		if err = db.AutoMigrate(admin.AllModels...); err != nil {
			Logger.WithError(err).WithField("type", dbType).Error("Error migrating models")
			return
		} else {
			Logger.WithField("type", dbType).Info("Migration complete")
		}
		if err = Initialize(dbType, db); err != nil {
			Logger.WithError(err).WithField("type", dbType).Error("Error initializing database")
			return
		}
	}
	return
}

// Initialize will create any requires inital data that the database might depend on
func Initialize(dbType MigrationType, db *gorm.DB) (err error) {
	switch dbType {
	case GameDBMigrationType:
		if err = initGameDB(db); err != nil {
			Logger.WithError(err).WithField("type", GameDBMigrationType).Error("Error initializing database")
		}
	case AdminDBMigrationType:
		if err = initAdminDB(db); err != nil {
			Logger.WithError(err).WithField("type", AdminDBMigrationType).Error("Error initializing database")
		}
	}
	return
}
