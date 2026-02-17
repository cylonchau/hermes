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

// Hermes 结构体
type Hermes struct {
	Next           plugin.Handler
	DatabaseConfig store.DatabaseConfig
	Resolver       *resolver.Resolver
}

// ServeDNS 处理 DNS 请求
func (h *Hermes) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	msg, err := h.Resolver.Resolve(ctx, state)
	if err != nil {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	_ = w.WriteMsg(msg)
	return dns.RcodeSuccess, nil
}

// Name 插件名称
func (h *Hermes) Name() string { return pluginName }

// initAdvancedDBPool 初始化数据库连接池
func (h *Hermes) initAdvancedDBPool() error {
	// 使用共享的 RDBStore 初始化数据库
	err := store.GetInstance().Initialize(h.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	logger.Info("数据库连接初始化成功", logger.Any("type", h.DatabaseConfig.Type))
	return nil
}

// Close 关闭数据库连接
func (h *Hermes) Close() error {
	return store.GetInstance().Close()
}

// GetDB 返回 GORM 数据库实例用于查询
func (h *Hermes) GetDB() *gorm.DB {
	return store.GetInstance().GetDB()
}

// HealthCheck 执行数据库健康检查
func (h *Hermes) HealthCheck() error {
	return store.GetInstance().HealthCheck()
}

// MonitorConnectionPool 监控数据库连接池状态
func (h *Hermes) MonitorConnectionPool() {
	// 将监控任务委托给 store 包处理
	ctx := context.Background()
	go store.GetInstance().MonitorConnectionPool(ctx)
}

// ValidateConfig 验证数据库配置
func (h *Hermes) ValidateConfig() error {
	// 目前依赖 Initialize 中的验证逻辑，此处预留
	return nil
}
