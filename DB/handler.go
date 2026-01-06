package DB

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func NewDB(path string) *DB {
	db, err := Init(path)
	if err != nil {
		panic(err)
	}
	return &DB{conn: db}
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func Init(path string) (*sql.DB, error) {
	base := fmt.Sprintf("file:%s", path)
	params := []string{
		"_pragma=busy_timeout(5000)",
		"_pragma=foreign_keys(1)",
		"_pragma=journal_mode(WAL)",
	}
	dsn := base + "?" + strings.Join(params, "&")

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil

}

func (db *DB) InitSchema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS drives (
  kind       TEXT NOT NULL,
  value      TEXT NOT NULL,
  uuid       TEXT NOT NULL,
  adopted_at TEXT NOT NULL,
  PRIMARY KEY (kind, value),
  UNIQUE (uuid)
); 
`
	_, err := db.conn.ExecContext(ctx, ddl)
	return err
}
