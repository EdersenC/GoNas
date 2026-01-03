package api

import (
	"goNAS/storage"

	"github.com/gin-gonic/gin"
)

type Nas struct {
	POOLS *storage.Pools
}

var NAS = &Nas{}

func Register(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	RegisterDrives(v1)
	RegisterPools(v1)
}

func Run(nas *Nas, addr string) error {
	NAS = nas
	r := gin.New()
	Register(r)
	err := r.Run(addr)
	if err != nil {
		return err
	}
	return nil
}
