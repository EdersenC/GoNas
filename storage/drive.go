package storage

import (
	"time"

	uuid2 "github.com/google/uuid"
)

type AdoptedDrive struct {
	Drive     *DriveInfo `json:"drive"`
	uuid      string
	CreatedAt string
}

func CreationTime() string { return time.Now().UTC().Format(time.RFC3339Nano) }

func NewAdoptedDrive(drive *DriveInfo) *AdoptedDrive {
	uuid := uuid2.New().String()
	drive.Uuid = uuid
	return &AdoptedDrive{
		Drive:     drive,
		uuid:      uuid,
		CreatedAt: CreationTime(),
	}
}

func (a *AdoptedDrive) Key() string          { return a.Drive.DriveKey.String() }
func (a *AdoptedDrive) GetKind() string      { return a.Drive.DriveKey.Kind }
func (a *AdoptedDrive) GetKindValue() string { return a.Drive.DriveKey.Value }
func (a *AdoptedDrive) GetUuid() string      { return a.uuid }
func (a *AdoptedDrive) SetUuid(id string)    { a.uuid = id; a.Drive.SetUuid(id) }
