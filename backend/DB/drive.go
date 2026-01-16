package DB

import (
	"context"
	"goNAS/storage"

	"gorm.io/gorm"
)

func (db *DB) InsertDrive(ctx context.Context, drive *storage.DriveInfo, createdAt string) error {
	model := &DriveModel{}
	model.FromDriveInfo(drive, createdAt)

	return db.conn.WithContext(ctx).Create(model).Error
}

func (db *DB) QueryDriveByKey(ctx context.Context, key storage.DriveKey) (storage.AdoptedDrive, bool, error) {
	var model DriveModel
	err := db.conn.WithContext(ctx).
		Where("kind = ? AND value = ?", key.Kind, key.Value).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return storage.AdoptedDrive{}, false, nil
	}
	if err != nil {
		return storage.AdoptedDrive{}, false, err
	}

	adoptedDrive := model.ToAdoptedDrive()
	return adoptedDrive, true, nil
}

func (db *DB) QueryAllAdoptedDrives(ctx context.Context) ([]storage.AdoptedDrive, error) {
	var models []DriveModel
	if err := db.conn.WithContext(ctx).Order("createdAt DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	adoptedDrives := make([]storage.AdoptedDrive, 0, len(models))
	for _, model := range models {
		adoptedDrive := model.ToAdoptedDrive()
		adoptedDrives = append(adoptedDrives, adoptedDrive)
	}

	return adoptedDrives, nil
}

type DrivePatch struct {
	PoolID *string
}

func (db *DB) PatchDrive(ctx context.Context, driveUuid string, p DrivePatch) error {
	updates := make(map[string]interface{})

	if p.PoolID != nil {
		updates["poolID"] = *p.PoolID
	}

	if len(updates) == 0 {
		return nil
	}

	result := db.conn.WithContext(ctx).Model(&DriveModel{}).
		Where("uuid = ?", driveUuid).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
