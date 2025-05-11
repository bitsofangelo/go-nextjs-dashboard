package db

import (
	"errors"
	"fmt"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-nextjs-dashboard/internal/config"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		url.QueryEscape(cfg.DBPass),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
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
