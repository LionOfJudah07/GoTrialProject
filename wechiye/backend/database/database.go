package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/001_initial.sql
var initialSchema string

type DB struct {
	*sql.DB
}

func Open(path string, key []byte) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma_key=x'%x'&_pragma_cipher_page_size=4096", path, key)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	_, _ = db.Exec("PRAGMA foreign_keys = ON;")
	_, _ = db.Exec("PRAGMA journal_mode = WAL;")
	_, _ = db.Exec("PRAGMA synchronous = NORMAL;")
	return &DB{db}, nil
}

func Migrate(db *DB) error {
	_, err := db.Exec(initialSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users (full_name, created_at, updated_at) VALUES ('', ?, ?)`, time.Now(), time.Now())
	return err
}

func (db *DB) Close() error {
	return db.DB.Close()
}