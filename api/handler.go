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
	server := &Server{
		Nas: NAS,
		httpServer: &http.Server{
			Addr:         addr,
			WriteTimeout: 0,
			ReadTimeout:  5 * time.Minute,
			IdleTimeout:  10 * time.Second,
			Handler:      r,
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
	err := SERVER.LoadData(context.Background())
	if err != nil {
		return nil
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) LoadData(c context.Context) error {
	err := s.Nas.LoadAdoptedDrives(c)
	if err != nil {
		return err
	}
	return nil
}

type Nas struct {
	POOLS         *storage.Pools
	SystemDrives  []*storage.DriveInfo
	AdoptedDrives []*storage.AdoptedDrive
}

var NAS = &Nas{}

func (n *Nas) LoadAdoptedDrives(c context.Context) error {
	adoptedDrives, err := SERVER.Db.QueryAllAdoptedDrives(c)
	if err != nil {
		return err
	}
	n.AdoptedDrives = make([]*storage.AdoptedDrive, 0, len(adoptedDrives))
	for i := range adoptedDrives {
		adoptedDrive := adoptedDrives[i]
		drive := n.getDriveByKey(adoptedDrives[i].Key())
		if drive != nil {
			drive.Uuid = adoptedDrive.GetUuid()
			adoptedDrive.Drive = drive
		}
		n.AdoptedDrives = append(n.AdoptedDrives, &adoptedDrive)
	}
	return nil
}

// GetAdoptedDriveByKey retrieves an adopted drive by its key.
func (n *Nas) GetAdoptedDriveByKey(key string) *storage.AdoptedDrive {
	for _, drive := range n.AdoptedDrives {
		if drive.Key() == key {
			return drive
		}
	}
	return nil
}

func (n *Nas) GetDriveByUuid(id string) *storage.DriveInfo {
	for _, drive := range n.AdoptedDrives {
		if drive.GetUuid() == id {
			return drive.Drive
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

func ensureUniqueKeys(keys ...string) error {
	keySet := make(map[string]bool)
	for _, key := range keys {
		if keySet[key] {
			return errors.New("duplicate drive key found: " + key)
		}
		keySet[key] = true
	}
	return nil
}

// PopulatePool populates a storage pool with the specified drives ids.
func (n *Nas) PopulatePool(pool *storage.Pool, drives []string, c *gin.Context) error {
	if err := ensureUniqueKeys(drives...); err != nil {
		return err
	}

	poolDrives := make([]*storage.DriveInfo, 0, len(drives))
	for _, driveID := range drives {
		drive := n.GetDriveByUuid(driveID)
		if drive == nil {
			return ErrDriveNotFound
		}
		poolDrives = append(poolDrives, drive)
	}

	pool.AddDrives(poolDrives...)
	err := SERVER.Db.InsertPool(c, pool, pool.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// AdoptDriveByKey adopts a drive by its key and returns the adopted drive.
// Returns an error if the drive is already adopted or not found.
func (n *Nas) AdoptDriveByKey(key string, c *gin.Context) (*storage.AdoptedDrive, error) {
	adopted := n.GetAdoptedDriveByKey(key)
	if adopted != nil {
		return nil, ErrAlreadyAdopted
	}
	drive := n.getDriveByKey(key)
	if drive == nil {
		return nil, ErrDriveNotFound
	}
	adoptedDrive := storage.NewAdoptedDrive(drive)
	err := SERVER.Db.InsertDrive(c, adoptedDrive.Drive, adoptedDrive.CreatedAt)
	if err != nil {
		return nil, err
	}
	n.AdoptedDrives = append(n.AdoptedDrives, adoptedDrive)
	return adoptedDrive, nil
}
