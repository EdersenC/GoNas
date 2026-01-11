package DB

import (
	"context"
	"database/sql"
	"fmt"
	"goNAS/storage"
	"strings"
)

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
  format    Text Null,
  createdAt  TEXT NOT NULL
);
`, quoteList(statuses), quoteList(poolTypes))

	_, err := db.conn.ExecContext(ctx, ddl)
	return err
}

func (db *DB) InsertPool(ctx context.Context, pool *storage.Pool, createdAt string) error {
	const q = `
INSERT INTO Pool (uuid, name, mdDevice, status, poolType, format,createdAt)
VALUES
  (?, ?, ?, ?, ?, ?,?);
`
	_, err := db.conn.ExecContext(ctx, q,
		pool.Uuid,
		pool.Name,
		pool.MdDevice,
		pool.Status.ToLower(),
		pool.Type.Value(),
		pool.Format,
		createdAt,
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
	Name   string         `json:"name"`
	Status storage.Status `json:"status"`
	Format string         `json:"format"`
}

func buildPoolPatch(pool *storage.Pool, patch *PoolPatch) (set []string, args []any, updatedPool *storage.Pool) {
	set = make([]string, 0, 4)
	args = make([]any, 0, 5)
	updatedPool = pool.Clone()

	if patch == nil {
		return
	}

	if patch.Name != "" {
		set = append(set, "Name = ?")
		args = append(args, patch.Name)
		updatedPool.SetName(patch.Name)
	}

	if patch.Status != "" {
		set = append(set, "status = ?")
		args = append(args, patch.Status.ToLower())
		updatedPool.SetStatus(patch.Status.ToLower())
	}

	if patch.Format != "" {
		set = append(set, "format = ?")
		args = append(args, patch.Format)
		updatedPool.SetFormat(patch.Format)
	}

	return
}

func (db *DB) PatchPoolMount(uuid, mount string) error {
	const q = `
UPDATE Pool SET mountPoint = ? WHERE uuid = ?;
`
	_, err := db.conn.ExecContext(context.Background(), q, mount, uuid)
	return err
}

func (db *DB) PatchPool(ctx context.Context, pool *storage.Pool, patch *PoolPatch) (*storage.Pool, error) {
	set, args, updatedPool := buildPoolPatch(pool, patch)

	if len(set) == 0 {
		return pool, nil
	}

	args = append(args, pool.Uuid)

	q := "UPDATE Pool SET " + strings.Join(set, ", ") + " WHERE uuid = ?"
	res, err := db.conn.ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}

	return updatedPool, nil
}

func (db *DB) QueryAllPools(ctx context.Context) (map[string]storage.Pool, error) {
	const q = `
SELECT uuid, name, mountPoint, mdDevice, status, poolType, format, createdAt
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
			format     sql.NullString
			createdAt  string
		)

		if err = rows.Scan(&uuid,
			&name,
			&mountPoint,
			&mdDevice,
			&status,
			&poolType,
			&format,
			&createdAt,
		); err != nil {
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
			Format:        format.String,
			CreatedAt:     createdAt,
		}

		pools[uuid] = p
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pools, nil
}
