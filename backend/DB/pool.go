package DB

import (
	"context"
	"goNAS/storage"
)

// InsertPool persists a new pool record.
func (db *DB) InsertPool(ctx context.Context, pool *storage.Pool, createdAt string) error {
	model := &PoolModel{}
	model.FromStoragePool(pool)
	model.CreatedAt = createdAt

	return db.conn.WithContext(ctx).Create(model).Error
}

// DeletePool removes a pool record by UUID.
func (db *DB) DeletePool(ctx context.Context, uuid string) error {
	return db.conn.WithContext(ctx).Delete(&PoolModel{}, "uuid = ?", uuid).Error
}

type PoolPatch struct {
	Name   string         `json:"name"`
	Status storage.Status `json:"status"`
	Format string         `json:"format"`
}

// applyPoolPatch returns a modified copy of the pool based on the patch.
func applyPoolPatch(pool *storage.Pool, patch *PoolPatch) (updatedPool *storage.Pool) {
	updatedPool = pool.Clone()

	if patch == nil {
		return
	}

	if patch.Name != "" {
		updatedPool.SetName(patch.Name)
	}

	if patch.Status != "" {
		updatedPool.SetStatus(patch.Status.ToLower())
	}

	if patch.Format != "" {
		updatedPool.SetFormat(patch.Format)
	}

	return
}

// PatchPoolMount updates the mount point for a pool record.
func (db *DB) PatchPoolMount(uuid, mount string) error {
	return db.conn.Model(&PoolModel{}).
		Where("uuid = ?", uuid).
		Update("mountPoint", mount).Error
}

// PatchPool applies a patch to a pool and persists changes.
func (db *DB) PatchPool(ctx context.Context, pool *storage.Pool, patch *PoolPatch) (*storage.Pool, error) {
	updatedPool := applyPoolPatch(pool, patch)

	// Build update map
	updates := make(map[string]interface{})
	if patch.Name != "" {
		updates["name"] = patch.Name
	}
	if patch.Status != "" {
		updates["status"] = patch.Status.ToLower()
	}
	if patch.Format != "" {
		updates["format"] = patch.Format
	}

	if len(updates) == 0 {
		return pool, nil
	}

	result := db.conn.WithContext(ctx).Model(&PoolModel{}).
		Where("uuid = ?", pool.Uuid).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, storage.ErrPoolNotFound
	}

	return updatedPool, nil
}

// QueryAllPools returns all pools indexed by UUID.
func (db *DB) QueryAllPools(ctx context.Context) (map[string]storage.Pool, error) {
	var models []PoolModel
	if err := db.conn.WithContext(ctx).Order("createdAt DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	pools := make(map[string]storage.Pool)
	for _, model := range models {
		pool, err := model.ToStoragePool()
		if err != nil {
			return nil, err
		}
		pools[model.UUID] = pool
	}

	return pools, nil
}
