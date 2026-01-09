package DB

import (
	"context"
	"database/sql"
	"fmt"
	"goNAS/storage"
	"strings"
)

// Todo implement network table creation
func (db *DB) createPoolTable(ctx context.Context) error {
	statuses := []string{
		string(storage.Healthy),
		string(storage.Offline),
		string(storage.Degraded),
	}
	poolTypes := []string{"raid0", "raid1", "raid5", "raid6", "raid10"}

	ddl := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS Pool (
  uuid        TEXT PRIMARY KEY,
  name        TEXT UNIQUE NOT NULL,
  mountPoint TEXT  NULL,
  mdDevice   TEXT  UNIQUE,
  status      TEXT NOT NULL CHECK (status IN (%s)),
  poolType   TEXT NOT NULL CHECK (poolType IN (%s)),
  createdAt  TEXT NOT NULL
);
`, quoteList(statuses), quoteList(poolTypes))

	_, err := db.conn.ExecContext(ctx, ddl)
	return err
}

func (db *DB) InsertPool(ctx context.Context, pool *storage.Pool, createdAt string) error {
	const q = `
INSERT INTO Pool (uuid, name, mdDevice, status, poolType, createdAt)
VALUES
  (?, ?, ?, ?, ?, ?);
`
	_, err := db.conn.ExecContext(ctx, q,
		pool.Uuid, pool.Name, pool.MdDevice, pool.Status, pool.Type.Value(), createdAt,
	)
	return err
}

func (db *DB) DeletePool(ctx context.Context, uuid string) error {
	const q = `
DELETE FROM Pool
WHERE uuid = ?;
`
	_, err := db.conn.ExecContext(ctx, q, uuid)
	return err
}

type PoolPatch struct {
	Status     *string
	MountPoint *string
}

func (db *DB) PatchPool(ctx context.Context, poolUUID string, p PoolPatch) error {
	set := make([]string, 0, 2)
	args := make([]any, 0, 3)

	if p.Status != nil {
		set = append(set, "status = ?")
		args = append(args, *p.Status)
	}
	if p.MountPoint != nil {
		set = append(set, "mount_point = ?")
		args = append(args, *p.MountPoint)
	}
	if len(set) == 0 {
		return nil
	}

	args = append(args, poolUUID)

	q := "UPDATE pools SET " + strings.Join(set, ", ") + " WHERE uuid = ?"
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

func (db *DB) QueryAllPools(ctx context.Context) (map[string]storage.Pool, error) {
	const q = `
SELECT uuid, name, mountPoint, mdDevice, status, poolType, createdAt
FROM Pool
ORDER BY createdAt DESC;
`

	rows, err := db.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	pools := make(map[string]storage.Pool)

	for rows.Next() {
		var (
			uuid       string
			name       string
			mountPoint sql.NullString
			mdDevice   sql.NullString
			status     string
			poolType   string
			createdAt  string
		)

		if err = rows.Scan(&uuid, &name, &mountPoint, &mdDevice, &status, &poolType, &createdAt); err != nil {
			return nil, err
		}

		convertedType, err := storage.ParsePoolType(poolType)
		if err != nil {
			return nil, err
		}

		p := storage.Pool{
			Uuid:          uuid,
			Name:          name,
			MdDevice:      mdDevice.String,
			AdoptedDrives: make(map[string]*storage.AdoptedDrive),
			Status:        storage.Status(status),
			Type:          convertedType,
			MountPoint:    mountPoint.String,
			CreatedAt:     createdAt,
		}

		pools[uuid] = p
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pools, nil
}
