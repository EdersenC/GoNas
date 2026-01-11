package api

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	RegisterDrives(v1)
	RegisterPools(v1)
}

func RegisterPools(r *gin.RouterGroup) {
	r.GET("/pools", listPools)
	r.GET("/pool/:uuid", getPool)
	r.POST("/pool/:uuid/build", buildPool)
	r.POST("/pool", createPool)
	r.PATCH("/pool/:uuid", updatePool)
	r.DELETE("/pool/:uuid", deletePool)
}

func RegisterDrives(r *gin.RouterGroup) {
	r.GET("/drives", func(c *gin.Context) {
		listDrives(c, false)
	})
	r.GET("/drives/scan", func(c *gin.Context) {
		listDrives(c, true)
	})
	r.GET("/drives/adopted", listAdoptedDrives)

	r.POST("/drives/adopt/:key", adoptDrive)
}
