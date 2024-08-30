package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func New() Database {
	return Database{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		name:     os.Getenv("POSTGRES_DB"),
	}
}

func (db *Database) buildConnectionURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.user, db.password, db.host, db.port, db.name)
}

func (db *Database) Connect() (*gorm.DB, error) {
	connectionURL := db.buildConnectionURL()
	connection, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return connection, nil
}

func (db *Database) Migrate() {
	connection, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = connection.AutoMigrate(&UserModel{})
	if err != nil {
		log.Fatalf("Failed to apply database migrations: %s", err)
	}
}
