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
	db, err := gorm.Open(mysql.Open(getDSN()), &gorm.Config{})
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

func getDSN() string {
	dbConfig := config.GetConfig().Database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		dbConfig.DBUser,
		dbConfig.DBPassword,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	)

	return dsn
}

func GetDB() *gorm.DB {
	return DB
}
