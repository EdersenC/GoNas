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
	if err := db.createPoolTable(ctx); err != nil {
		return err
	}

	if err := db.createDriveTable(ctx); err != nil {
		return err
	}
	return nil
}

func quoteList(vals []string) string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		v = strings.ReplaceAll(v, "'", "''")
		out = append(out, fmt.Sprintf("'%s'", v))
	}
	return strings.Join(out, ", ")
}
