package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Amele9/call-manager/internal/config"
	"github.com/Amele9/call-manager/internal/server"
)

const ShutdownTimeout = 30 * time.Second

// App is the structure of the application
type App struct {
	provider *Provider

	server server.Server
}

// New returns an instance of the App
func New() (*App, error) {
	configuration, err := config.Get()
	if err != nil {
		return nil, err
	}

	provider, err := NewProvider(configuration)
	if err != nil {
		return nil, err
	}

	database, err := provider.Database()
	if err != nil {
		return nil, err
	}

	server_, err := server.New(configuration, database)
	if err != nil {
		return nil, err
	}

	return &App{
		provider: provider,

		server: server_,
	}, nil
}

// Run starts the application and waits for a signal
func (a *App) Run() error {
	go func() {
		err := a.server.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	return a.ShutdownOnSignal()
}

// ShutdownOnSignal shuts down the application on a signal
func (a *App) ShutdownOnSignal() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = a.provider.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
