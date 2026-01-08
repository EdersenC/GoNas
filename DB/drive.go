package DB

import (
	"context"
	"database/sql"
	"errors"
	"goNAS/storage"
)

func (db *DB) InsertDrive(ctx context.Context, drive *storage.DriveInfo, createdAt string) error {
	const q = `
INSERT INTO drive (kind, value, uuid, createdAt)
VALUES
  (?, ?, ?, ?);
`
	_, err := db.conn.ExecContext(ctx, q,
		drive.DriveKey.Kind, drive.DriveKey.Value, drive.Uuid, createdAt,
	)
	return err
}

func (db *DB) QueryDriveByKey(ctx context.Context, key storage.DriveKey) (storage.AdoptedDrive, bool, error) {
	const q = `
SELECT uuid, createdAt 
FROM drive 
WHERE kind = ? AND value = ?
LIMIT 1;
`

	var (
		adoptedDrive storage.AdoptedDrive
		drive        storage.DriveInfo
		uuid         string
	)

	drive.DriveKey = key
	adoptedDrive.Drive = &drive

	err := db.conn.QueryRowContext(ctx, q, key.Kind, key.Value).
		Scan(&uuid, &adoptedDrive.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return storage.AdoptedDrive{}, false, nil
	}
	if err != nil {
		return storage.AdoptedDrive{}, false, err
	}

	return adoptedDrive, true, nil
}

func (db *DB) QueryAllAdoptedDrives(ctx context.Context) ([]storage.AdoptedDrive, error) {
	const q = `
SELECT kind, value, uuid,createdAt
FROM drive
ORDER BY createdAt DESC;
`

	rows, err := db.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	adoptedDrives := make([]storage.AdoptedDrive, 0, 10)

	for rows.Next() {
		var (
			adoptedDrive storage.AdoptedDrive
			drive        storage.DriveInfo
			uuid         string
		)

		adoptedDrive.Drive = &drive

		if err = rows.Scan(&drive.DriveKey.Kind, &drive.DriveKey.Value, &uuid, &adoptedDrive.CreatedAt); err != nil {
			return nil, err
		}

		adoptedDrive.SetUuid(uuid)
		adoptedDrives = append(adoptedDrives, adoptedDrive)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adoptedDrives, nil
}
