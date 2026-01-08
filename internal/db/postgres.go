package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string ) (*sql.DB, error) {
	db,err := sql.Open("postgres",dsn)
	if err != nil {
		return nil, err
	}
	
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db,nil
}
