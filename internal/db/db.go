package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB() (*sqlx.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "master"
	}
	pass := os.Getenv("POSTGRES_PASSWORD")
	if pass == "" {
		pass = "master"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "master"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
