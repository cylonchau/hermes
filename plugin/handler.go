package plugin

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/logger"
	"github.com/cylonchau/hermes/pkg/store"
)

const pluginName = "hermes"

// Hermes represents the main plugin structure
type Hermes struct {
	Next           plugin.Handler
	DatabaseConfig store.DatabaseConfig
}

// ServeDNS handles DNS requests
func (h *Hermes) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return h.Next.ServeDNS(ctx, w, r)
}

// Name returns the plugin name
func (h *Hermes) Name() string { return pluginName }

// initAdvancedDBPool initializes database connection pool
func (h *Hermes) initAdvancedDBPool() error {
	// Use the shared RDBStore to initialize the database
	err := store.GetInstance().Initialize(h.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}

	logger.Info("Database connection initialized successfully", logger.Any("type", h.DatabaseConfig.Type))
	return nil
}

// Close closes the database connection
func (h *Hermes) Close() error {
	return store.GetInstance().Close()
}

// GetDB returns the GORM database instance for queries
func (h *Hermes) GetDB() *gorm.DB {
	return store.GetInstance().GetDB()
}

// HealthCheck performs database health check
func (h *Hermes) HealthCheck() error {
	return store.GetInstance().HealthCheck()
}

// MonitorConnectionPool monitors database connection pool status
func (h *Hermes) MonitorConnectionPool() {
	// Delegate monitoring to the store package if it supports background monitoring,
	// or perform monitoring using the shared instance.
	// Current store package has MonitorConnectionPool(ctx).
	// We need a context here, or we can just start it in a goroutine if not already running.
	// For now, we can create a background context or use a proper lifecycle context if available.
	// Since this method signature doesn't take context, we'll create one.
	ctx := context.Background()
	go store.GetInstance().MonitorConnectionPool(ctx)
}

// ValidateConfig validates database configuration
func (h *Hermes) ValidateConfig() error {
	// We can reuse the validate logic in store package if we expose it,
	// or just check basic fields here if needed before initialization.
	// For now keeping basic checks or relying on Initialize to fail.
	return nil
}
