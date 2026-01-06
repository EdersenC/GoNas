package api

import (
	"context"
	"errors"
	"goNAS/DB"
	"goNAS/storage"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
)

type Server struct {
	Nas        *Nas
	httpServer *http.Server
	Ctx        *context.Context
	Db         *DB.DB
}

var SERVER = &Server{}

func NewAPIServer(addr string, db *DB.DB) *Server {
	NAS = &Nas{POOLS: &storage.Pools{}}
	r := gin.New()
	r.Use(gin.Recovery())
	Register(r)
	var server = &Server{
		Nas: NAS,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: r,
		},
		Db: db,
	}
	SERVER = server
	return server
}

func (s *Server) Start() error {
	s.Nas.SystemDrives = storage.GetSystemDrives()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

type Nas struct {
	POOLS         *storage.Pools
	SystemDrives  []*storage.DriveInfo
	AdoptedDrives []*AdoptedDrive
}

var NAS = &Nas{}

type AdoptedDrive struct {
	Drive     *storage.DriveInfo `json:"drive"`
	uuid      string
	AdoptedAt string
}

func NewAdoptedDrive(drive *storage.DriveInfo) *AdoptedDrive {
	uuid := uuid2.New().String()
	drive.Uuid = uuid
	return &AdoptedDrive{
		Drive:     drive,
		uuid:      uuid,
		AdoptedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func (a AdoptedDrive) key() string          { return a.Drive.DriveKey.String() }
func (a AdoptedDrive) GetKind() string      { return a.Drive.DriveKey.Kind }
func (a AdoptedDrive) GetKindValue() string { return a.Drive.DriveKey.Value }
func (a AdoptedDrive) Uuid() string         { return a.uuid }

// GetAdoptedDriveByKey retrieves an adopted drive by its key.
func (n *Nas) GetAdoptedDriveByKey(key string) *AdoptedDrive {
	for _, drive := range n.AdoptedDrives {
		if drive.key() == key {
			return drive
		}
	}
	return nil
}

// getDriveByKey retrieves a drive from the system drives by its key.
func (n *Nas) getDriveByKey(key string) *storage.DriveInfo {
	for _, drive := range n.SystemDrives {
		if drive.DriveKey.String() == key {
			return drive
		}
	}
	return nil
}

// AdoptDriveByKey adopts a drive by its key and returns the adopted drive.
// Returns an error if the drive is already adopted or not found.
func (n *Nas) AdoptDriveByKey(key string, c *gin.Context) (*AdoptedDrive, error) {
	adopted := n.GetAdoptedDriveByKey(key)
	if adopted != nil {
		return nil, ErrAlreadyAdopted
	}
	drive := n.getDriveByKey(key)
	if drive == nil {
		return nil, ErrDriveNotFound
	}
	adoptedDrive := NewAdoptedDrive(drive)
	err := SERVER.Db.InsertDrive(c, adoptedDrive.Drive, adoptedDrive.AdoptedAt)
	if err != nil {
		return nil, err
	}
	n.AdoptedDrives = append(n.AdoptedDrives, adoptedDrive)
	return adoptedDrive, nil
}
