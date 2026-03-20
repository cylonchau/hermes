package plugin

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/logger"
	"github.com/cylonchau/hermes/pkg/resolver"
	"github.com/cylonchau/hermes/pkg/store"
)

const pluginName = "hermes"

// Hermes struct
type Hermes struct {
	Next           plugin.Handler
	DatabaseConfig store.DatabaseConfig
	Resolver       *resolver.Resolver
	GeoIPPath      string
	CacheSizeMB    int // Cache Size limit, Unit: MB
}


// ServeDNS handles DNS requests
func (h *Hermes) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	msg, err := h.Resolver.Resolve(ctx, state)
	if err != nil {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	_ = w.WriteMsg(msg)
	return dns.RcodeSuccess, nil
}

// Name returns the plugin name
func (h *Hermes) Name() string { return pluginName }

// initAdvancedDBPool initializes database connection pool
func (h *Hermes) initAdvancedDBPool() error {
	// Initialize DB using shared RDBStore
	err := store.GetInstance().Initialize(h.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	logger.Info("Database connection initialized successfully", logger.Any("type", h.DatabaseConfig.Type))
	return nil
}

// Close closes database connections
func (h *Hermes) Close() error {
	return store.GetInstance().Close()
}

// GetDB returns GORM database instance for queries
func (h *Hermes) GetDB() *gorm.DB {
	return store.GetInstance().GetDB()
}

// HealthCheck executes database healthcheck
func (h *Hermes) HealthCheck() error {
	return store.GetInstance().HealthCheck()
}

// MonitorConnectionPool monitors database connection pool status
func (h *Hermes) MonitorConnectionPool() {
	// Delegate monitoring tasks to store package
	ctx := context.Background()
	go store.GetInstance().MonitorConnectionPool(ctx)
}

// ValidateConfig validates database configuration
func (h *Hermes) ValidateConfig() error {
	// Currently relies on Initialize validation, reserved here
	return nil
}
