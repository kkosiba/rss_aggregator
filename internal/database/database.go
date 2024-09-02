package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (db *Database) Connect() *pgxpool.Pool {
	connectionURL := db.buildConnectionURL()
	dbpool, err := pgxpool.New(context.Background(), connectionURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()
	return dbpool
}
