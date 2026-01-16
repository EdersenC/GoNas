package DB

import (
	"goNAS/storage"
	"time"

	"gorm.io/gorm"
)

// PoolModel represents the Pool table in GORM
type PoolModel struct {
	UUID       string `gorm:"primaryKey;column:uuid"`
	Name       string `gorm:"unique;not null;column:name"`
	MountPoint string `gorm:"column:mountPoint"`
	MdDevice   string `gorm:"unique;column:mdDevice"`
	Status     string `gorm:"not null;column:status"`
	PoolType   string `gorm:"not null;column:poolType"`
	Format     string `gorm:"column:format"`
	CreatedAt  string `gorm:"not null;column:createdAt"`
}

// TableName sets the table name for GORM
func (PoolModel) TableName() string {
	return "Pool"
}

// ToStoragePool converts GORM model to storage.Pool
func (p *PoolModel) ToStoragePool() (storage.Pool, error) {
	convertedType, err := storage.ParsePoolType(p.PoolType)
	if err != nil {
		return storage.Pool{}, err
	}

	return storage.Pool{
		Uuid:          p.UUID,
		Name:          p.Name,
		MdDevice:      p.MdDevice,
		AdoptedDrives: make(map[string]*storage.AdoptedDrive),
		Status:        storage.Status(p.Status),
		Type:          convertedType,
		MountPoint:    p.MountPoint,
		Format:        p.Format,
		CreatedAt:     p.CreatedAt,
	}, nil
}

// FromStoragePool converts storage.Pool to GORM model
func (p *PoolModel) FromStoragePool(pool *storage.Pool) {
	p.UUID = pool.Uuid
	p.Name = pool.Name
	p.MdDevice = pool.MdDevice
	p.Status = string(pool.Status.ToLower())
	p.PoolType = pool.Type.Value()
	p.Format = pool.Format
	p.CreatedAt = pool.CreatedAt
	p.MountPoint = pool.MountPoint
}

// DriveModel represents the Drive table in GORM
type DriveModel struct {
	Kind      string     `gorm:"primaryKey;not null;column:kind"`
	Value     string     `gorm:"primaryKey;not null;column:value"`
	UUID      string     `gorm:"unique;not null;column:uuid"`
	PoolID    *string    `gorm:"column:poolID"` // Pointer handles NULL (nil = NULL in DB)
	CreatedAt string     `gorm:"not null;column:createdAt"`
	Pool      *PoolModel `gorm:"foreignKey:PoolID;references:UUID;constraint:OnDelete:SET NULL;"`
}

// TableName sets the table name for GORM
func (DriveModel) TableName() string {
	return "Drive"
}

// ToAdoptedDrive converts GORM model to storage.AdoptedDrive
func (d *DriveModel) ToAdoptedDrive() storage.AdoptedDrive {
	drive := &storage.DriveInfo{
		DriveKey: storage.DriveKey{
			Kind:  d.Kind,
			Value: d.Value,
		},
		Uuid: d.UUID,
	}

	adoptedDrive := storage.AdoptedDrive{
		Drive:     drive,
		CreatedAt: d.CreatedAt,
	}

	adoptedDrive.SetUuid(d.UUID)
	if d.PoolID != nil {
		adoptedDrive.SetPoolID(*d.PoolID)
	}

	return adoptedDrive
}

// FromDriveInfo converts storage.DriveInfo to GORM model
func (d *DriveModel) FromDriveInfo(drive *storage.DriveInfo, createdAt string) {
	d.Kind = drive.DriveKey.Kind
	d.Value = drive.DriveKey.Value
	d.UUID = drive.Uuid
	d.CreatedAt = createdAt
}

// BeforeCreate hook to set default timestamp if not provided
func (p *PoolModel) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt == "" {
		p.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	return nil
}

// BeforeCreate hook to set default timestamp if not provided
func (d *DriveModel) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedAt == "" {
		d.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	return nil
}
