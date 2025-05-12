package db

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/logger"
)

func Open(cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	wd, _ := os.Getwd()

	logLevel := gormlogger.Warn
	l := &dbLogger{Logger: log, level: logLevel, basePath: wd}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		url.QueryEscape(cfg.DBPass),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: gormlogger.New(l, gormlogger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logLevel,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db connection: %v", err)
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	//	return nil, err
	// }
	// sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

type dbLogger struct {
	logger.Logger
	level    gormlogger.LogLevel
	basePath string
}

func (l dbLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, l.basePath+string(filepath.Separator), "")

	switch l.level {
	case gormlogger.Error:
		l.Error(msg)
	case gormlogger.Warn:
		l.Warn(msg)
	default:
		l.Info(msg)
	}
}

func RecordExists(tx *gorm.DB) (bool, error) {
	var hit int

	if err := tx.Select("1").Limit(1).Scan(&hit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("record exists: %w", err)
	}

	return hit == 1, nil
}
