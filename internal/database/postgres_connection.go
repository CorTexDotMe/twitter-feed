package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"twitter-feed/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxRetries = 10
	delay      = 10 * time.Second
)

type PostgresDB struct {
	host   string
	port   string
	user   string
	dbName string
	db     *gorm.DB
	once   sync.Once
}

func (p *PostgresDB) NewDB() (*gorm.DB, error) {
	var err error
	p.once.Do(func() {
		if p.db, err = p.connect(); err != nil {
			return
		}
		if err = p.createDatabase(); err != nil {
			return
		}
		if err := p.migrate(); err != nil {
			return
		}
	})
	return p.db, err
}

func (p *PostgresDB) connect() (*gorm.DB, error) {
	for i := 0; i < maxRetries; i++ {
		connectionProperties := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			p.host, p.port, p.user, p.dbName)
		db, err := gorm.Open(postgres.Open(connectionProperties), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("postgress connection failed: will retry...")
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("failed to connect to postgres after %d retries", maxRetries)
}

func (p *PostgresDB) createDatabase() error {
	return p.db.Exec("CREATE DATABASE IF NOT EXISTS " + p.dbName + ";").Error
}

func (p *PostgresDB) migrate() error {
	return p.db.AutoMigrate(&model.Message{})
}

func NewPostgresDB() PostgresDB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	return PostgresDB{
		host:   host,
		port:   port,
		user:   user,
		dbName: dbName,
	}
}
