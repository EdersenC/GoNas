package api

import (
	"errors"
	"fmt"
	"goNAS/DB"
	"goNAS/helper"
	"goNAS/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// poolError writes a pool-related error response with the appropriate status.
func (n *Nas) poolError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch {
	case errors.Is(err, storage.ErrPoolNotFound):
		c.JSON(http.StatusNotFound, message)
	case errors.Is(err, storage.ErrPoolInUse):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrDriveNotFoundOrInUse):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrInsufficientDrives):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrPoolAlreadyExists):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrPoolNotOffline):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrPoolFormatRequired):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrInvalidPoolType):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrUnsupportedFormat):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrInvalidRequestBody):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, helper.ErrRaid0RequiresDrives),
		errors.Is(err, helper.ErrRaid1RequiresDrives),
		errors.Is(err, helper.ErrRaid5RequiresDrives),
		errors.Is(err, helper.ErrRaid6RequiresDrives),
		errors.Is(err, helper.ErrRaid10RequiresDrives):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, helper.ErrUnsupportedRaidLevel):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, helper.ErrInvalidSizeInput),
		errors.Is(err, helper.ErrInvalidAmountInput),
		errors.Is(err, helper.ErrRootPrivilegesNeeded),
		errors.Is(err, helper.ErrWorkdirResolve),
		errors.Is(err, helper.ErrPackageManagerMissing),
		errors.Is(err, helper.ErrMdadmInstall),
		errors.Is(err, helper.ErrMdadmInstallVerify),
		errors.Is(err, helper.ErrMdadmArgsEmpty),
		errors.Is(err, helper.ErrMdadmBuild),
		errors.Is(err, helper.ErrMountPointCreate),
		errors.Is(err, helper.ErrMountRaidDevice),
		errors.Is(err, helper.ErrFormatRaidDevice):
		c.JSON(http.StatusInternalServerError, message)
	case errors.Is(err, storage.ErrPoolNotInMemory):
		c.JSON(http.StatusNotFound, message)
	case errors.Is(err, storage.ErrPoolDeleteUnmount),
		errors.Is(err, storage.ErrPoolDeleteRmdir),
		errors.Is(err, storage.ErrPoolDeleteRemove),
		errors.Is(err, storage.ErrPoolDeleteStop),
		errors.Is(err, storage.ErrPoolDeleteZeroSB),
		errors.Is(err, storage.ErrPoolCapacityRead),
		errors.Is(err, storage.ErrPoolCapacityParse):
		c.JSON(http.StatusInternalServerError, message)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
	}
}

// driveError writes a drive-related error response with the appropriate status.
func (n *Nas) driveError(err error, c *gin.Context) {
	message := gin.H{"error": err.Error()}
	switch {
	case errors.Is(err, storage.ErrDriveNotFound):
		c.JSON(http.StatusNotFound, message)
	case errors.Is(err, storage.ErrAlreadyAdopted):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrDriveNotFoundOrInUse):
		c.JSON(http.StatusConflict, message)
	case errors.Is(err, storage.ErrNoDrivesToRemove):
		c.JSON(http.StatusBadRequest, message)
	case errors.Is(err, storage.ErrDuplicateDriveKey):
		c.JSON(http.StatusBadRequest, message)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: " + err.Error()})
	}
}

// listAdoptedDrives returns all adopted drives.
func listAdoptedDrives(c *gin.Context) {
	SuccessResponse(c, NAS.AdoptedDrives)
}

// Todo Make UUID System for drives
// adoptDrive adopts a system drive by its key.
func adoptDrive(c *gin.Context) {
	key := c.Param("key")
	driveToAdopt, err := NAS.AdoptDriveByKey(key, c)
	if err != nil {
		NAS.driveError(err, c)
		return
	}
	SuccessResponse(c, driveToAdopt)
}

// listDrives returns known drives, optionally rescanning system devices.
func listDrives(c *gin.Context, rescan bool) {
	if len(NAS.SystemDrives) == 0 || rescan {
		NAS.SystemDrives = storage.GetSystemDriveMap()
	}
	SuccessResponse(c, NAS.SystemDrives)
}

// listPools returns all pools from memory.
func listPools(c *gin.Context) {
	SuccessResponse(c, NAS.POOLS)
}

// createPool validates input, persists, and optionally builds a pool.
func createPool(c *gin.Context) {
	var req struct {
		Name      string   `json:"name" binding:"required"`
		RaidLevel *int     `json:"raidLevel" binding:"required"`
		Drives    []string `json:"drives" binding:"required"`
		Format    string   `json:"format" binding:"required"`
		Build     bool     `json:"build"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		NAS.poolError(fmt.Errorf("%w: %v", storage.ErrInvalidRequestBody, err), c)
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

	err = NAS.AddPool(pool, c)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	_ = NAS.RemoveAdoptedDrives(req.Drives, c) // Clean up adopted drives after pool creation

	if req.Build {
		err = pool.Build()
		if err != nil {
			NAS.poolError(err, c)
			return
		}
	}

	SuccessResponse(c, pool)
}

// deletePool removes the pool from the database and memory.
func deletePool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	err = NAS.deletePool(pool)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	err = SERVER.Db.DeletePool(c, pool.Uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	SuccessResponse(c, gin.H{"Deleted": pool.Uuid})
}

// getPool returns a pool by UUID.
func getPool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}
	SuccessResponse(c, pool)
}

// updatePool applies a patch to a pool in storage and memory.
func updatePool(c *gin.Context) {
	uuid := c.Param("uuid")
	pool, err := NAS.POOLS.GetPool(uuid)
	if err != nil {
		NAS.poolError(err, c)
		return
	}

	var req *DB.PoolPatch

	if err = c.ShouldBindJSON(&req); err != nil {
		NAS.poolError(fmt.Errorf("%w: %v", storage.ErrInvalidRequestBody, err), c)
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

// buildPool builds an existing pool and persists its mount point.
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

// SuccessResponse writes a standard success response envelope.
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
