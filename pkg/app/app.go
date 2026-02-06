package app

import (
	"fmt"
	"net/http"

	"github.com/cylonchau/hermes/pkg/config"
)

// NewHTTPSever 启动 HTTP 管理服务
func NewHTTPSever() error {
	cfg := config.Get()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Hermes HTTP server listening on %s\n", addr)

	// TODO: 注册路由和处理器
	return http.ListenAndServe(addr, nil)
}
