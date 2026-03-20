package plugin

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/dao/memory"
	"github.com/cylonchau/hermes/pkg/resolver"
	"github.com/cylonchau/hermes/pkg/store"
)

// init registers the plugin
func init() {
	caddy.RegisterPlugin(pluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

// setup parses Corefile configuration
func setup(c *caddy.Controller) error {
	h, err := parseHermes(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	// Register startup and shutdown hooks
	c.OnStartup(func() error {
		err := h.initAdvancedDBPool()
		if err != nil {
			return err
		}

		var geoip resolver.GeoIPProvider
		if h.GeoIPPath != "" {
			geoip, err = resolver.NewMaxMindProvider(h.GeoIPPath)
			if err != nil {
				return plugin.Error(pluginName, err)
			}
		}

		// Initialize resolver
		cacheSize := 32 * 1024 * 1024 // Default 32MB L1 Cache
		if h.CacheSizeMB > 0 {
			cacheSize = h.CacheSizeMB * 1024 * 1024
		}
		cache := memory.NewCacheDAO(cacheSize)
		rdbDAO := rdb.NewRecordDAO(h.GetDB())
		cachedDAO := rdb.NewCachedDNSQueryRepository(rdbDAO, cache) // Mount L1 memory cache proxy
		h.Resolver = resolver.NewResolver(cachedDAO, h.GetDB(), geoip)
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

// parseHermes parses hermes configuration block
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

				// Parse database sub-block
				for c.NextBlock() {
					val := c.Val()
					switch val {
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
					case "{", "}":
						// Fault tolerance: explicitly ignore braces
						continue
					default:
						return nil, c.Errf("unknown db property: %s", val)
					}
				}
			case "cache_size":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				size, err := strconv.Atoi(c.Val())
				if err != nil {
					return nil, c.Errf("invalid cache_size value: %s", c.Val())
				}
				h.CacheSizeMB = size
			case "geoip":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				h.GeoIPPath = c.Val()
			default:
				return nil, c.Errf("unknown property: %s", c.Val())
			}
		}
	}

	return h, nil
}
