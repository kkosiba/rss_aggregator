package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (db *Database) buildConnectionURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}

func (db *Database) Connect() (*gorm.DB, error) {
	connectionURL := db.buildConnectionURL()
	connection, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return connection, nil
}
