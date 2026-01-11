package api

import (
	"goNAS/DB"
	"goNAS/helper"
	"goNAS/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (n *Nas) poolError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch err {
	case storage.ErrPoolNotFound:
		c.JSON(http.StatusNotFound, message)
	case storage.ErrPoolInUse:
		c.JSON(http.StatusConflict, message)
	case storage.ErrDriveNotFoundOrInUse:
		c.JSON(http.StatusConflict, message)
	case storage.ErrInsufficientDrives:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrPoolAlreadyExists:
		c.JSON(http.StatusConflict, message)
	case storage.ErrPoolNotOffline:
		c.JSON(http.StatusConflict, message)
	case storage.ErrPoolFormatRequired:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrInvalidPoolType:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrUnsupportedFormat:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrInvalidStatus:
		c.JSON(http.StatusBadRequest, message)
	case helper.ErrRaid0RequiresDrives, helper.ErrRaid1RequiresDrives, helper.ErrRaid5RequiresDrives, 
		 helper.ErrRaid6RequiresDrives, helper.ErrRaid10RequiresDrives:
		c.JSON(http.StatusBadRequest, message)
	case helper.ErrUnsupportedRaidLevel:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrPoolNotInMemory:
		c.JSON(http.StatusNotFound, message)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
	}
}

func (n *Nas) driveError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch err {
	case storage.ErrDriveNotFound:
		c.JSON(http.StatusNotFound, message)
	case storage.ErrAlreadyAdopted:
		c.JSON(http.StatusConflict, message)
	case storage.ErrDriveNotFoundOrInUse:
		c.JSON(http.StatusConflict, message)
	case storage.ErrNoDrivesToRemove:
		c.JSON(http.StatusBadRequest, message)
	case storage.ErrDuplicateDriveKey:
		c.JSON(http.StatusBadRequest, message)
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
		RaidLevel *int     `json:"raidLevel" binding:"required"`
		Drives    []string `json:"drives" binding:"required"`
		Format    string   `json:"format" binding:"required"`
		Build     bool     `json:"build"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pool, err := NAS.POOLS.NewPool(req.Name, &storage.Raid{Level: *req.RaidLevel}, req.Format)
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
		err = pool.Build()
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
	_ = NAS.RemoveAdoptedDrives(req.Drives, c) // Clean up adopted drives after pool creation
	SuccessResponse(c, pool.Uuid)
}

func deletePool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	err = SERVER.Db.DeletePool(c, pool.Uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	err = NAS.deletePool(pool)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	SuccessResponse(c, gin.H{"deleted": pool.Uuid})
}

func getPool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	SuccessResponse(c, pool)
}

func updatePool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	var req *DB.PoolPatch

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = NAS.ValidatePoolPatch(req); err != nil {
		NAS.poolError(err, c)
		return
	}
	updatedPool, err := SERVER.Db.PatchPool(c, pool, req)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	err = NAS.updatePool(updatedPool)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	SuccessResponse(c, updatedPool)
}

func buildPool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	err = pool.Build()
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	err = SERVER.Db.PatchPoolMount(pool.Uuid, pool.MountPoint)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	SuccessResponse(c, gin.H{"built": pool.Uuid})
}
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
