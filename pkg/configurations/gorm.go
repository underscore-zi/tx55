package configurations

import (
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type LogLevelOptions string

const (
	Silent LogLevelOptions = "silent"
	Error  LogLevelOptions = "error"
	Warn   LogLevelOptions = "warn"
	Info   LogLevelOptions = "info"
	Debug  LogLevelOptions = "debug"
)

func (l LogLevelOptions) LogrusLevel() logrus.Level {
	switch l {
	case Silent:
		return logrus.PanicLevel
	case Error:
		return logrus.ErrorLevel
	case Warn:
		return logrus.WarnLevel
	case Info:
		return logrus.InfoLevel
	case Debug:
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

type DatabaseType string

const (
	MySQL     DatabaseType = "mysql"
	SQLite    DatabaseType = "sqlite"
	SQLServer DatabaseType = "sqlserver"
	Postgres  DatabaseType = "postgres"
)

type DatabaseConfig struct {
	// Type must be one of mysql, sqlite, sqlserver or postgres as those are the only supported dialects by GORM.
	Type DatabaseType

	// DSN is the Data Source Name string, contains credentials and connection information. It is important the parameter `parseTime=true` is also included.
	// Example: user:password@tcp(localhost:3306)/databaseName?charset=utf8mb4&parseTime=True&loc=Local
	DSN                string
	LogConfig          GormLoggingConfig
	MaxOpenConnections int
	MaxIdleConnections int
	// ConnectionMaxLifetime is the maximum amount of time a connection may be reused, read in nanoseconds (60000000000 == 1 minute)
	ConnectionMaxLifetime time.Duration
}

// GormLoggingConfig basically wraps the logger.Config struct from GORM. (https://gorm.io/docs/logger.html)
type GormLoggingConfig struct {
	Enable                    bool
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	SlowThresholdMS           int
	Level                     LogLevelOptions
}

func (cfg DatabaseConfig) Dialect() (d gorm.Dialector, err error) {
	switch cfg.Type {
	case MySQL:
		d = mysql.Open(cfg.DSN)
	case SQLite:
		d = sqlite.Open(cfg.DSN)
	case SQLServer:
		d = sqlserver.Open(cfg.DSN)
	case Postgres:
		d = postgres.Open(cfg.DSN)
	default:
		err = errors.New("Unknown database type: " + string(cfg.Type))
	}
	return
}

func (cfg DatabaseConfig) Open(config *gorm.Config) (db *gorm.DB, err error) {
	var dialect gorm.Dialector

	dialect, err = cfg.Dialect()
	if err != nil {
		return
	}

	db, err = gorm.Open(dialect, config)
	if err != nil {
		return
	}

	var rawDb *sql.DB
	rawDb, err = db.DB()
	if err != nil {
		return
	}

	rawDb.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	rawDb.SetMaxOpenConns(cfg.MaxOpenConnections)
	rawDb.SetMaxIdleConns(cfg.MaxIdleConnections)

	return
}

func (cfg GormLoggingConfig) LoggerConfig() logger.Config {
	gormLogConfig := logger.Config{
		SlowThreshold:             time.Millisecond * time.Duration(cfg.SlowThresholdMS),
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.ParameterizedQueries,
		Colorful:                  cfg.Colorful,
	}

	switch cfg.Level {
	case Silent:
		gormLogConfig.LogLevel = logger.Silent
	case Error:
		gormLogConfig.LogLevel = logger.Error
	case Warn:
		gormLogConfig.LogLevel = logger.Warn
	case Info:
		gormLogConfig.LogLevel = logger.Info
	case Debug:
		panic("Debug logging is not supporting in GormLoggingConfig")
	}
	return gormLogConfig
}
