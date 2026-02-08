package plugin

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/resolver"
	"github.com/cylonchau/hermes/pkg/store"
)

// init 注册插件
func init() {
	caddy.RegisterPlugin(pluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

// setup 解析 Corefile 配置
func setup(c *caddy.Controller) error {
	h, err := parseHermes(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	// 注册启动和关闭钩子
	c.OnStartup(func() error {
		err := h.initAdvancedDBPool()
		if err != nil {
			return err
		}
		// 初始化解析器
		h.Resolver = resolver.NewResolver(rdb.NewRecordDAO(h.GetDB()))
		return nil
	})

	c.OnShutdown(func() error {
		return h.Close()
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		h.Next = next
		return h
	})

	return nil
}

// parseHermes 解析 hermes 配置块
func parseHermes(c *caddy.Controller) (*Hermes, error) {
	h := &Hermes{}

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "db":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				dbTypeStr := c.Val()
				switch dbTypeStr {
				case "mysql":
					h.DatabaseConfig.Type = store.MySQL
				case "postgres":
					h.DatabaseConfig.Type = store.PostgreSQL
				case "sqlite":
					h.DatabaseConfig.Type = store.SQLite
				default:
					return nil, c.Errf("unsupported database type: %s", dbTypeStr)
				}

				// 解析数据库子块
				for c.NextBlock() {
					switch c.Val() {
					case "host":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.Host = c.Val()
					case "port":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						port, _ := strconv.Atoi(c.Val())
						h.DatabaseConfig.Port = port
					case "username":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.Username = c.Val()
					case "password":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.Password = c.Val()
					case "database":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.Database = c.Val()
					case "file":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.File = c.Val()
					case "ssl_mode":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.SSLMode = c.Val()
					case "max_open":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.MaxOpenConnection = c.Val()
					case "max_idle":
						if !c.NextArg() {
							return nil, c.ArgErr()
						}
						h.DatabaseConfig.MaxIdleConnection = c.Val()
					default:
						return nil, c.Errf("unknown db property: %s", c.Val())
					}
				}
			default:
				return nil, c.Errf("unknown property: %s", c.Val())
			}
		}
	}

	return h, nil
}
