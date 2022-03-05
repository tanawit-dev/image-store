package db

import (
	"fmt"
	"github.com/tanawit-dev/image-store/config"
	"github.com/tanawit-dev/image-store/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func ProvideDatabase(config config.Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		config.Database.DBUser,
		config.Database.DBPassword,
		config.Database.DBHost,
		config.Database.DBPort,
		config.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10) // TODO read from config
	sqlDB.SetMaxOpenConns(100) // TODO read from config
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("connect database successfully")

	defer migrateDB(db)

	return db, nil
}

func migrateDB(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Image{}); err != nil {
		log.Fatalln(err)
	}
}
