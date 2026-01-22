package api

import (
	"context"
	"errors"
	"goNAS/DB"
	"goNAS/storage"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
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
	s.Nas.SystemDrives = storage.GetSystemDriveMap()
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
	log.Println("Server started on", s.httpServer.Addr)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) LoadData(c context.Context) error {
	err := s.Nas.LoadPools(c)
	if err != nil {
		return err
	}
	err = s.Nas.LoadAdoptedDrives(c)
	if err != nil {
		return err
	}
	return nil
}

type Nas struct {
	POOLS         *storage.Pools
	SystemDrives  map[string]*storage.DriveInfo
	AdoptedDrives map[string]*storage.AdoptedDrive
}

var NAS = &Nas{}

func (n *Nas) LoadAdoptedDrives(c context.Context) error {
	adoptedDrives, err := SERVER.Db.QueryAllAdoptedDrives(c)
	if err != nil {
		return err
	}
	n.AdoptedDrives = make(map[string]*storage.AdoptedDrive)
	for i := range adoptedDrives {
		adoptedDrive := adoptedDrives[i]
		drive := n.getDriveByKey(adoptedDrives[i].Key())
		if err = n.ClaimDrive(drive, adoptedDrive); err != nil {
			log.Println("Error claiming drive:", err)
		}
	}
	return nil
}

func (n *Nas) ClaimDrive(drive *storage.DriveInfo, adoptedDrive storage.AdoptedDrive) error {
	if drive == nil {
		return storage.ErrDriveNotFound
	}
	drive.Uuid = adoptedDrive.GetUuid()
	adoptedDrive.Drive = drive

	if adoptedDrive.GetPoolID() != "" {
		pool, err := n.POOLS.GetPool(adoptedDrive.GetPoolID())
		if err != nil {
			return err
		}
		pool.AddDrives(drive)
		return nil
	}

	//if drive is not part of a pool, add to adopted drives
	n.AdoptedDrives[adoptedDrive.GetUuid()] = &adoptedDrive
	return nil
}

func (n *Nas) LoadPools(c context.Context) error {
	pools, err := SERVER.Db.QueryAllPools(c)
	if err != nil {
		return err
	}
	n.POOLS = &storage.Pools{}
	for _, pool := range pools {
		err = n.POOLS.AddPool(&pool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Nas) updatePool(pool *storage.Pool) error {
	if _, exists := (*n.POOLS)[pool.Uuid]; !exists {
		return storage.ErrPoolNotInMemory
	}
	(*n.POOLS)[pool.Uuid] = pool
	return nil
}

func (n *Nas) deletePool(p *storage.Pool) error {
	adoptedDrives := p.AdoptedDrives
	for _, adopt := range adoptedDrives {
		adopt.SetPoolID("")
		n.AdoptedDrives[adopt.GetUuid()] = adopt
	}
	err := n.setOffline(p)
	if err != nil {
		return err
	}

	err = NAS.POOLS.DeletePool(p.Uuid)
	if err != nil {
		return err
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
	if drive, exists := n.AdoptedDrives[id]; exists {
		return drive.Drive
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
			return storage.ErrDuplicateDriveKey
		}
		keySet[key] = true
	}
	return nil
}

func (n *Nas) AddPool(p *storage.Pool, c *gin.Context) error {
	err := n.POOLS.AddPool(p)
	if err != nil {
		return err
	}
	drivePatch := DB.DrivePatch{
		PoolID: &p.Uuid,
	}

	for _, drive := range p.AdoptedDrives {
		err = SERVER.Db.PatchDrive(c, drive.GetUuid(), drivePatch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Nas) AreDrivesAlreadyInPool(d []string) (string, bool) {
	for _, uuid := range d {
		adoptedDrive, ok := n.AdoptedDrives[uuid]
		if !ok {
			continue
		}
		if adoptedDrive.GetPoolID() != "" {
			return uuid + ":" + adoptedDrive.GetPoolID(), true
		}
	}
	return "", false
}

// RemoveAdoptedDrives removes an in mem adopted drive by its UUID.
// Returns an error if the drive is not found or if it is currently in use by a pool.
func (n *Nas) RemoveAdoptedDrives(uuid []string, c *gin.Context) int {
	removedCount := 0
	for _, id := range uuid {
		if _, exists := n.AdoptedDrives[id]; !exists {
			continue
		}
		removedCount++
		delete(n.AdoptedDrives, id)
	}
	return removedCount
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
			return storage.ErrDriveNotFoundOrInUse
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
		return nil, storage.ErrAlreadyAdopted
	}
	drive := n.getDriveByKey(key)
	if drive == nil {
		return nil, storage.ErrDriveNotFound
	}
	adoptedDrive := storage.NewAdoptedDrive(drive)
	err := SERVER.Db.InsertDrive(c, adoptedDrive.Drive, adoptedDrive.CreatedAt)
	if err != nil {
		return nil, err
	}
	n.AdoptedDrives[adoptedDrive.GetUuid()] = adoptedDrive
	return adoptedDrive, nil
}

func (n *Nas) ValidatePoolPatch(patch *DB.PoolPatch) error {
	if patch.Status != "" {
		err := storage.ValidateStatus(patch.Status)
		if err != nil {
			return err
		}
	}

	if patch.Format != "" {
		err := storage.ValidatePoolFormat(patch.Format)
		if err != nil {
			return err
		}
	}

	return nil
}

// setOffline sets the status of the pool to offline.
// Todo: implement actual offline logic
func (n *Nas) setOffline(pool *storage.Pool) error {
	pool.SetStatus(storage.Offline)

	return nil
}
