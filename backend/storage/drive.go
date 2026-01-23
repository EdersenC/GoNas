package storage

import (
	"goNAS/helper"
	"time"

	uuid2 "github.com/google/uuid"
)

// CreationTime returns the current UTC time in RFC3339Nano format.
func CreationTime() string { return time.Now().UTC().Format(time.RFC3339Nano) }

type AdoptedDrive struct {
	Drive     *DriveInfo `json:"drive"`
	Uuid      string     `json:"uuid"`
	PoolID    string     `json:"poolID"`
	CreatedAt string     `json:"createdAt"`
}

// NewAdoptedDrive creates an adopted drive with a stable UUID.
func NewAdoptedDrive(drive *DriveInfo) *AdoptedDrive {
	uuid := ""
	if drive.GetUuid() != "" {
		uuid = drive.GetUuid()
	} else {
		uuid = uuid2.New().String()
		drive.Uuid = uuid
	}
	return &AdoptedDrive{
		Drive:     drive,
		Uuid:      uuid,
		CreatedAt: CreationTime(),
	}
}

// Key returns the DriveKey string for the adopted drive.
func (a *AdoptedDrive) Key() string          { return a.Drive.DriveKey.String() }
// GetKind returns the drive key kind.
func (a *AdoptedDrive) GetKind() string      { return a.Drive.DriveKey.Kind }
// GetKindValue returns the drive key value.
func (a *AdoptedDrive) GetKindValue() string { return a.Drive.DriveKey.Value }

// GetUuid returns the adopted drive UUID.
func (a *AdoptedDrive) GetUuid() string   { return a.Uuid }
// SetUuid sets the adopted drive UUID and updates the underlying DriveInfo.
func (a *AdoptedDrive) SetUuid(id string) { a.Uuid = id; a.Drive.SetUuid(id) }

// GetPoolID returns the owning pool UUID for the adopted drive.
func (a *AdoptedDrive) GetPoolID() string   { return a.PoolID }
// SetPoolID associates the adopted drive with a pool UUID.
func (a *AdoptedDrive) SetPoolID(id string) { a.PoolID = id }

// GetSystemDrives returns system drives filtered by name and minimum size.
func GetSystemDrives(names ...string) []*DriveInfo {
	drives, _ := GetDrives()
	drives = FilterFor(DriveFilter{
		Names:   names,
		MinSize: 1 * helper.Gigabyte,
	}, drives...)
	return drives
}

// GetSystemDriveMap returns a map of system drives keyed by drive key string.
func GetSystemDriveMap(names ...string) map[string]*DriveInfo {
	drives := GetSystemDrives(names...)
	driveMap := make(map[string]*DriveInfo)
	for _, d := range drives {
		driveMap[d.DriveKey.String()] = d
	}
	return driveMap
}
