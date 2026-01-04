package api

import (
	"context"
	"errors"
	"goNAS/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Nas        *Nas
	httpServer *http.Server
	Ctx        *context.Context
}

func NewAPIServer(addr string) *Server {
	NAS = &Nas{POOLS: &storage.Pools{}}
	r := gin.New()
	r.Use(gin.Recovery())
	Register(r)
	return &Server{
		Nas: NAS,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
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
	POOLS  *storage.Pools
	Drives *storage.DriveInfo
}

var NAS = &Nas{}

func Register(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	RegisterDrives(v1)
	RegisterPools(v1)
}
