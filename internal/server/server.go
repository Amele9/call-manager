package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Amele9/call-manager/internal/config"
	"github.com/Amele9/call-manager/internal/database"
)

type Server interface {
	// Run starts the server
	Run() error

	// Shutdown shuts down the server
	Shutdown(ctx context.Context) error
}

// New returns an instance of the GinServer
func New(config *config.Configuration, database database.Database) (*GinServer, error) {
	server := &GinServer{Database: database}

	router := server.registerRoutes(gin.Default())

	server.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return server, nil
}

// GinServer is the structure for working with the Gin server
type GinServer struct {
	server *http.Server

	Database database.Database
}

// Run starts the Gin server
func (s *GinServer) Run() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	return nil
}

// Shutdown shuts down the Gin server
func (s *GinServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
