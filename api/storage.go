package api

import (
	"errors"
	"goNAS/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrDriveNotFound            = errors.New("drive not found")
	DriveNotFoundOrAlreadyInUse = errors.New("drive not found or already in use")
	ErrAlreadyAdopted           = errors.New("drive already adopted")

	ErrPoolNotFound       = errors.New("pool not found")
	ErrPoolInUse          = errors.New("pool is currently in use")
	ErrInsufficientDrives = errors.New("insufficient drives for the requested pool type")
)

func (n *Nas) poolError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch {
	case errors.Is(err, ErrPoolNotFound):
		c.JSON(http.StatusNotFound, message)
	case errors.Is(err, ErrPoolInUse):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, DriveNotFoundOrAlreadyInUse):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, ErrInsufficientDrives):
		c.JSON(http.StatusBadRequest, message)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
	}
}

func (n *Nas) driveError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch {
	case errors.Is(err, ErrDriveNotFound):
		c.JSON(http.StatusNotFound, message)
	case errors.Is(err, ErrAlreadyAdopted):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, DriveNotFoundOrAlreadyInUse):
		c.JSON(http.StatusConflict, message)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
	}
}

func listAdoptedDrives(c *gin.Context) {
	SuccessResponse(c, NAS.AdoptedDrives)
}

// Todo Make UUID System for drives
func adoptDrive(c *gin.Context) {
	key := c.Param("key")
	driveToAdopt, err := NAS.AdoptDriveByKey(key, c)
	if err != nil {
		NAS.driveError(err, c)
		return
	}
	SuccessResponse(c, driveToAdopt)
}

func listDrives(c *gin.Context, rescan bool) {
	if len(NAS.SystemDrives) == 0 || rescan {
		NAS.SystemDrives = storage.GetSystemDriveMap()
	}
	SuccessResponse(c, NAS.SystemDrives)
}

func listPools(c *gin.Context) {
	SuccessResponse(c, NAS.POOLS)
}

func createPool(c *gin.Context) {
	var req struct {
		Name      string   `json:"name" binding:"required"`
		RaidLevel int      `json:"raidLevel" binding:"required"`
		Drives    []string `json:"drives" binding:"required"`
		Format    string   `json:"format"`
		Build     bool     `json:"build"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pool, err := NAS.POOLS.NewPool(req.Name, &storage.Raid{Level: req.RaidLevel}, nil)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	err = NAS.PopulatePool(pool, req.Drives, c)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	if req.Build {
		err = pool.Build(req.Format)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err = NAS.AddPool(pool, c)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	_ = NAS.RemoveAdoptedDrives(req.Drives, c)
	SuccessResponse(c, pool.Uuid)
}

func deletePool(c *gin.Context) {}

func getPool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	SuccessResponse(c, pool)
}

func updatePool(c *gin.Context) {}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
