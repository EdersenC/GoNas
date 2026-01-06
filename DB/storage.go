package DB

import (
	"context"
	"database/sql"
	"errors"
	"goNAS/storage"
)

func (db *DB) InsertDrive(ctx context.Context, drive *storage.DriveInfo, adoptedAt string) error {
	const q = `
INSERT INTO drives (kind, value, uuid, adopted_at)
VALUES
  (?, ?, ?, ?);
`
	_, err := db.conn.ExecContext(ctx, q,
		drive.DriveKey.Kind, drive.DriveKey.Value, drive.Uuid, adoptedAt,
	)
	return err
}

type DriveRow struct {
	Kind      string
	Value     string
	UUID      string
	AdoptedAt string // or time.Time if you scan/parse
}

func (db *DB) queryDriveByKey(ctx context.Context, key storage.DriveKey) (DriveRow, bool, error) {
	const q = `
SELECT uuid, adopted_at
FROM drives
WHERE kind = ? AND value = ?
`
	var row DriveRow
	err := db.conn.QueryRowContext(ctx, q, key.Kind, key.Value).
		Scan(&row.Kind, &row.Value, &row.UUID, &row.AdoptedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return DriveRow{}, false, nil
	}
	if err != nil {
		return DriveRow{}, false, err
	}
	return row, true, nil
}
