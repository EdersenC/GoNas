package DB

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	conn *gorm.DB
}

func NewDB(path string) *DB {
	db, err := Init(path)
	if err != nil {
		panic(err)
	}
	return &DB{conn: db}
}

func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Init(path string) (*gorm.DB, error) {
	// Configure GORM with SQLite driver and pragmas
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", path)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)

	// Explicitly enable foreign keys (in case DSN pragma doesn't work)
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) InitSchema(ctx context.Context) error {
	// Use GORM AutoMigrate to create tables with foreign key constraints
	// The Pool relationship in DriveModel will ensure the foreign key is created
	if err := db.conn.WithContext(ctx).AutoMigrate(&PoolModel{}, &DriveModel{}); err != nil {
		return err
	}

	return nil
}
