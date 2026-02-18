package app

import (
	"fmt"

	"github.com/cylonchau/hermes/pkg/app/router"
	"github.com/cylonchau/hermes/pkg/config"
	"github.com/gin-gonic/gin"
)

// NewHTTPSever 启动 HTTP 管理服务
func NewHTTPSever() error {
	cfg := config.Get()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Hermes HTTP server listening on %s\n", addr)

	// Initialize Gin
	engine := gin.Default()

	// Register Routers
	router.RegisteredRouter(engine)

	// Start Server
	return engine.Run(addr)
}
