package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"psam/database/models"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		if err := os.MkdirAll("data", os.ModePerm); err != nil {
			fmt.Printf("failed to create directory %s", err)
			return
		}
		db, err = gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database:", err)
		}

		err = db.AutoMigrate(&models.APIKey{})
		if err != nil {
			log.Fatal("failed to migrate database:", err)
		}
	})
	return db
}
