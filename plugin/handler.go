package plugin

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/glebarez/sqlite"
	"github.com/miekg/dns"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"

	"github.com/cylonchau/hermes/pkg/store"
)

const pluginName = "hermes"

// Hermes represents the main plugin structure
type Hermes struct {
	Next           plugin.Handler
	DatabaseConfig store.DatabaseConfig
	db             *gorm.DB
}

// ServeDNS handles DNS requests
func (h *Hermes) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return h.Next.ServeDNS(ctx, w, r)
}

// Name returns the plugin name
func (h *Hermes) Name() string { return pluginName }

// initAdvancedDBPool initializes database connection pool
func (h *Hermes) initAdvancedDBPool() error {
	newLogger := dblogger.KlogLogger{LogLevel: logger.Info}

	gormConfig := &gorm.Config{
		Logger: newLogger,
	}

	var err error

	switch h.DatabaseConfig.Type {
	case MySQL:
		err = h.initMySQL(gormConfig)
	case SQLite:
		err = h.initSQLite(gormConfig)
	case PostgreSQL:
		err = h.initPostgreSQL(gormConfig)
	default:
		return fmt.Errorf("unsupported database type: %v", h.DatabaseConfig.Type)
	}

	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}

	// Configure connection pool
	if err := h.configureConnectionPool(); err != nil {
		return fmt.Errorf("connection pool configuration failed: %w", err)
	}

	klog.V(2).Infof("Database connection initialized successfully, type: %v", h.DatabaseConfig.Type)
	return nil
}

// initMySQL initializes MySQL connection
func (h *Hermes) initMySQL(config *gorm.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		h.DatabaseConfig.Username,
		h.DatabaseConfig.Password,
		h.DatabaseConfig.Host,
		h.DatabaseConfig.Port,
		h.DatabaseConfig.Database,
	)

	var err error
	h.db, err = gorm.Open(mysql.Open(dsn), config)
	return err
}

// initSQLite initializes SQLite connection
func (h *Hermes) initSQLite(config *gorm.Config) error {
	dbPath := h.DatabaseConfig.File + ".db"

	var err error
	h.db, err = gorm.Open(sqlite.Open(dbPath), config)
	return err
}

// initPostgreSQL initializes PostgreSQL connection
func (h *Hermes) initPostgreSQL(config *gorm.Config) error {
	sslMode := h.DatabaseConfig.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		h.DatabaseConfig.Host,
		h.DatabaseConfig.Username,
		h.DatabaseConfig.Password,
		h.DatabaseConfig.Database,
		h.DatabaseConfig.Port,
		sslMode,
	)

	var err error
	h.db, err = gorm.Open(postgres.Open(dsn), config)
	return err
}

// configureConnectionPool configures database connection pool parameters
func (h *Hermes) configureConnectionPool() error {
	// dbConn is the underlying *sql.DB instance used for configuring connection pool parameters
	// GORM uses the database/sql package internally, and we can access the underlying connection via DB() method
	dbConn, err := h.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	// Set maximum number of open connections
	maxOpen, _ := strconv.Atoi(h.DatabaseConfig.MaxOpenConnection)
	if (maxOpen) <= 0 {
		maxOpen = 25
	}
	dbConn.SetMaxOpenConns(maxOpen)

	// Set maximum number of idle connections
	maxIdle, _ := strconv.Atoi(h.DatabaseConfig.MaxIdleConnection)
	if maxIdle <= 0 {
		maxIdle = 10
	}
	dbConn.SetMaxIdleConns(maxIdle)

	// Set maximum lifetime of connections
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Set maximum idle time for connections
	dbConn.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("database ping test failed: %w", err)
	}

	// Print connection pool statistics
	stats := dbConn.Stats()
	klog.V(4).Infof("Database connection pool stats: MaxOpen=%d, Open=%d, InUse=%d, Idle=%d",
		stats.MaxOpenConnections,
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
	)

	return nil
}

// Close closes the database connection
func (h *Hermes) Close() error {
	if h.db != nil {
		dbConn, err := h.db.DB()
		if err != nil {
			return err
		}
		return dbConn.Close()
	}
	return nil
}

// GetDB returns the GORM database instance for queries
func (h *Hermes) GetDB() *gorm.DB {
	return h.db
}

// HealthCheck performs database health check
func (h *Hermes) HealthCheck() error {
	if h.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	dbConn, err := h.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return dbConn.PingContext(ctx)
}

// MonitorConnectionPool monitors database connection pool status
func (h *Hermes) MonitorConnectionPool() {
	if h.db == nil {
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		dbConn, err := h.db.DB()
		if err != nil {
			klog.Errorf("Failed to get database connection for monitoring: %v", err)
			continue
		}

		stats := dbConn.Stats()
		klog.V(5).Infof("Connection pool status - Open: %d, InUse: %d, Idle: %d, WaitCount: %d",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
		)

		// Log warning if wait count is too high
		if stats.WaitCount > 10 {
			klog.Warningf("High database connection pool wait queue: %d", stats.WaitCount)
		}
	}
}

// ValidateConfig validates database configuration
func (h *Hermes) ValidateConfig() error {
	config := h.DatabaseConfig

	switch config.Type {
	case MySQL, PostgreSQL:
		if config.Host == "" {
			return fmt.Errorf("database host is required for %s", config.Type)
		}
		if config.Port <= 0 {
			return fmt.Errorf("valid database port is required for %s", config.Type)
		}
		if config.Database == "" {
			return fmt.Errorf("database name is required for %s", config.Type)
		}
		if config.Username == "" {
			return fmt.Errorf("database username is required for %s", config.Type)
		}
	case SQLite:
		if config.File == "" {
			return fmt.Errorf("database file path is required for SQLite")
		}
	default:
		return fmt.Errorf("unsupported database type: %v", config.Type)
	}

	return nil
}
