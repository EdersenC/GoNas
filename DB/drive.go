package DB

import (
	"context"
	"database/sql"
	"errors"
	"goNAS/storage"
	"strings"
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

func (db *DB) createDriveTable(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS Drive (
  kind       TEXT NOT NULL,
  value      TEXT NOT NULL,
  uuid       TEXT NOT NULL,
  poolID     TEXT NULL,               
  createdAt  TEXT NOT NULL,

  PRIMARY KEY (kind, value),
  UNIQUE (uuid),

  FOREIGN KEY (poolID)
    REFERENCES Pool(uuid)
    ON DELETE SET NULL
);
`
	_, err := db.conn.ExecContext(ctx, ddl)
	return err
}

func (db *DB) QueryAllAdoptedDrives(ctx context.Context) ([]storage.AdoptedDrive, error) {
	const q = `
SELECT kind, value, uuid, poolID ,createdAt
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
			poolID       sql.NullString
		)

		adoptedDrive.Drive = &drive

		if err = rows.Scan(&drive.DriveKey.Kind, &drive.DriveKey.Value, &uuid, &poolID, &adoptedDrive.CreatedAt); err != nil {
			return nil, err
		}

		if poolID.Valid {
			adoptedDrive.SetPoolID(poolID.String)
		}

		adoptedDrive.SetUuid(uuid)
		adoptedDrives = append(adoptedDrives, adoptedDrive)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return adoptedDrives, nil
}

type DrivePatch struct {
	PoolID *string
}

func (db *DB) PatchDrive(ctx context.Context, driveUuid string, p DrivePatch) error {
	set := make([]string, 0, 1)
	args := make([]any, 0, 2)

	if p.PoolID != nil {
		set = append(set, "poolID = ?")
		args = append(args, *p.PoolID)
	}
	if len(set) == 0 {
		return nil
	}
	args = append(args, driveUuid)

	q := "UPDATE Drive SET " + strings.Join(set, ", ") + " WHERE uuid = ?"
	res, err := db.conn.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
