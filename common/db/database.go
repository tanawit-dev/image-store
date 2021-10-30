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

var DB *gorm.DB

func Init() {
	c := config.GetConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		c.Database.DBUser,
		c.Database.DBPassword,
		c.Database.DBHost,
		c.Database.DBPort,
		c.Database.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	log.Println("connect database successfully")

	migrateDB()
}

func migrateDB() {
	if err := DB.AutoMigrate(&models.Image{}); err != nil {
		log.Fatalln(err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
