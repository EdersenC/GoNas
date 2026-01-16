package storage

import (
	"goNAS/helper"
	"time"

	uuid2 "github.com/google/uuid"
)

func CreationTime() string { return time.Now().UTC().Format(time.RFC3339Nano) }

type AdoptedDrive struct {
	Drive     *DriveInfo `json:"drive"`
	Uuid      string     `json:"uuid"`
	PoolID    string     `json:"poolID"`
	CreatedAt string     `json:"createdAt"`
}

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

func (a *AdoptedDrive) Key() string          { return a.Drive.DriveKey.String() }
func (a *AdoptedDrive) GetKind() string      { return a.Drive.DriveKey.Kind }
func (a *AdoptedDrive) GetKindValue() string { return a.Drive.DriveKey.Value }

func (a *AdoptedDrive) GetUuid() string   { return a.Uuid }
func (a *AdoptedDrive) SetUuid(id string) { a.Uuid = id; a.Drive.SetUuid(id) }

func (a *AdoptedDrive) GetPoolID() string   { return a.PoolID }
func (a *AdoptedDrive) SetPoolID(id string) { a.PoolID = id }

func GetSystemDrives(names ...string) []*DriveInfo {
	drives, _ := GetDrives()
	drives = FilterFor(DriveFilter{
		Names:   names,
		MinSize: 1 * helper.Gigabyte,
	}, drives...)
	return drives
}

func GetSystemDriveMap(names ...string) map[string]*DriveInfo {
	drives := GetSystemDrives(names...)
	driveMap := make(map[string]*DriveInfo)
	for _, d := range drives {
		driveMap[d.DriveKey.String()] = d
	}
	return driveMap
}
