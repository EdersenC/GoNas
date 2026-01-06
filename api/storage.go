package api

import (
	"errors"
	"goNAS/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrDriveNotFound  = errors.New("drive not found")
	ErrAlreadyAdopted = errors.New("drive already adopted")
)

func (n *Nas) driveError(err error, c *gin.Context) {
	switch {
	case errors.Is(err, ErrDriveNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrAlreadyAdopted):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		NAS.SystemDrives = storage.GetSystemDrives()
	}
	SuccessResponse(c, NAS.SystemDrives)
}

func listPools(c *gin.Context) {
	SuccessResponse(c, NAS.POOLS)
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}
