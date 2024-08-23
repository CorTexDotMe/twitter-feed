package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"twitter-feed/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection() *gorm.DB {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	connectionProperties := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(connectionProperties), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Message{})

	return db
}
