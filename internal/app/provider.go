package app

import (
	"context"

	"github.com/Amele9/call-manager/internal/config"
	"github.com/Amele9/call-manager/internal/database"
)

// Provider contains application services
type Provider struct {
	config *config.Configuration

	database database.Database
}

// NewProvider returns an instance of the Provider
func NewProvider(config *config.Configuration) (*Provider, error) {
	return &Provider{config: config}, nil
}

// Database returns the database service
func (p *Provider) Database() (database.Database, error) {
	if p.database == nil {
		database_, err := database.New(p.config)
		if err != nil {
			return nil, err
		}

		p.database = database_
	}

	return p.database, nil
}

// Shutdown shuts down all application services
func (p *Provider) Shutdown(ctx context.Context) error {
	return p.database.Shutdown(ctx)
}
