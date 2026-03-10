package conectdb

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func New() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		return nil, errors.New("DB_DSN belum di set")
	}

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
