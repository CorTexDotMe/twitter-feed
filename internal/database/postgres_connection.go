package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"twitter-feed/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection() *gorm.DB {
	connectionProperties := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"))

	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(postgres.Open(connectionProperties), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Connection failed(%s); will retry...", err.Error())
		time.Sleep(5 * time.Second)
	}

	db.AutoMigrate(&model.Message{})

	return db
}
