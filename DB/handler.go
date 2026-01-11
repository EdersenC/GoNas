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
	// Manually create the Pool table first
	if err := db.conn.WithContext(ctx).AutoMigrate(&PoolModel{}); err != nil {
		return err
	}

	// Then create the Drive table with foreign key
	// We need to create it manually to ensure proper foreign key with ON DELETE SET NULL
	createDriveTable := `
		CREATE TABLE IF NOT EXISTS Drive (
			kind TEXT NOT NULL,
			value TEXT NOT NULL,
			uuid TEXT NOT NULL UNIQUE,
			poolID TEXT NULL,
			createdAt TEXT NOT NULL,
			PRIMARY KEY (kind, value),
			FOREIGN KEY (poolID) REFERENCES Pool(uuid) ON DELETE SET NULL
		)
	`
	if err := db.conn.WithContext(ctx).Exec(createDriveTable).Error; err != nil {
		// Table might already exist, try AutoMigrate instead
		if err := db.conn.WithContext(ctx).AutoMigrate(&DriveModel{}); err != nil {
			return err
		}
	}

	return nil
}
