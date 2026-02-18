package app

import (
	"fmt"

	"github.com/cylonchau/hermes/pkg/app/router"
	"github.com/cylonchau/hermes/pkg/config"
	"github.com/cylonchau/hermes/pkg/logger"

	"github.com/gin-gonic/gin"
)

// NewHTTPSever 启动 HTTP 管理服务
func NewHTTPSever() error {
	cfg := config.Get()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Hermes HTTP server listening", logger.String("addr", addr))

	// Initialize Gin
	engine := gin.Default()

	// Register Routers
	router.RegisteredRouter(engine)

	// Start Server
	return engine.Run(addr)
}
