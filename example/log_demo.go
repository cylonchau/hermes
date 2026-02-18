package main

import (
	"fmt"

	"github.com/cylonchau/hermes/pkg/logger"
)

func main() {
	// 1. Initialize with custom configuration
	// This shows how to setup multiple drivers and formats
	config := logger.Config{
		Loggers: map[string]logger.LoggerConfig{
			"business": {
				Enabled: true,
				Level:   logger.LevelInfo,
				Format:  logger.FormatJSON, // Use JSON for business logs
				Outputs: []logger.OutputConfig{
					{Type: "stdout"},
					{
						Type: "file",
						File: logger.FileConfig{
							Filename:   "logs/app.log",
							MaxSize:    "10M",
							MaxBackups: 3,
							MaxAge:     "7d",
							Compress:   true,
						},
					},
				},
			},
			"sql": {
				Enabled: true,
				Level:   logger.LevelDebug,
				Format:  logger.FormatText, // Use Text for SQL logs
				Outputs: []logger.OutputConfig{
					{Type: "stdout"},
					{
						Type: "udp",
						UDP:  logger.UDPConfig{Addr: "localhost:514"},
					},
				},
			},
		},
	}

	if err := logger.Initialize(config); err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return
	}

	biz := logger.GetLogger("business")
	sql := logger.GetLogger("sql")

	// 2. Business logs will now be in JSON format (Stdout + File)
	biz.Info("Starting service in JSON format", logger.String("mode", "production"))

	// 3. SQL logs will be in Text format (Stdout + UDP Placeholder)
	sql.Debug("executing database query",
		logger.String("query", "SELECT * FROM zones"),
		logger.Int("rows", 5),
	)

	// 4. Convenience global functions (uses "business" logger by default)
	logger.Warn("This is a global warning log")
}
