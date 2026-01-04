package api

import (
	"goNAS/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterDrives(r *gin.RouterGroup) {
	r.GET("/drives", listDrives)
	r.POST("/:drive_id/adopt", adoptDrive)
}

func adoptDrive(c *gin.Context) {

}

func listDrives(c *gin.Context) {
	drives, err := storage.GetDrives()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, drives)
}

func RegisterPools(r *gin.RouterGroup) {
	r.GET("/pools", listPools)
}

func listPools(c *gin.Context) {
	c.JSON(http.StatusAccepted, NAS.POOLS)
}
