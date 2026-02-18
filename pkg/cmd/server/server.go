package server

import (
	"flag"
	"fmt"

	"github.com/cylonchau/hermes/pkg/app"
	"github.com/cylonchau/hermes/pkg/config"
	"github.com/cylonchau/hermes/pkg/logger"
	"github.com/cylonchau/hermes/pkg/migration"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/cylonchau/hermes/pkg/store"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Options struct {
	ConfigFile string
	h          bool
	migration  bool
	upgrade    bool
	sqlDriver  string
	errCh      chan error
}

func NewOptions() *Options {
	return &Options{}
}

// NewCommand creates a *cobra.Command object with default parameters
func NewCommand() *cobra.Command {
	opts := NewOptions()

	cmd := &cobra.Command{
		Use:  "server",
		Long: `Hermes server provides DNS governance and management services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			PrintFlags(cmd.Flags())
			if err := opts.Complete(); err != nil {
				return fmt.Errorf("failed complete: %w", err)
			}

			if err := opts.Run(); err != nil {
				logger.Error("Error running "+config.CONFIG.AppName, logger.Err(err))
				return err
			}

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine)

	_ = cmd.MarkFlagFilename("config", "yaml", "yml", "json")

	return cmd
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigFile, "config", "./config.yaml", "The path to the configuration file.")
	fs.BoolVar(&o.migration, "migration", false, "Initial database and tables.")
	fs.BoolVar(&o.upgrade, "upgrade", false, "If true, update the database schema to the latest version.")
	fs.StringVar(&o.sqlDriver, "sql-driver", "sqlite", "enable which sql backend.")
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		logger.Debug(fmt.Sprintf("FLAG: --%s=%q", flag.Name, flag.Value))
	})
}

func (o *Options) Complete() error {
	if len(o.ConfigFile) == 0 {
		logger.Warn("Warning, all flags other than --config are deprecated")
	}
	// Load the config file here in Complete
	if len(o.ConfigFile) > 0 {
		err := config.InitConfiguration(o.ConfigFile)
		if err != nil {
			return err
		}
		// Initialize the logger after config is loaded
		if err := logger.Initialize(logger.Config{Loggers: config.CONFIG.Loggers}); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
	}

	return nil
}

func (o *Options) Run() error {
	// 1. Initialize DB Store first
	if config.CONFIG != nil {
		var dbCfg store.DatabaseConfig
		switch o.sqlDriver {
		case "mysql":
			dbCfg = config.CONFIG.MySQL
			dbCfg.Type = store.MySQL
		case "sqlite":
			dbCfg = config.CONFIG.SQLite
			dbCfg.Type = store.SQLite
		default:
			dbCfg = config.CONFIG.Database // Fallback
		}

		if !dbCfg.IsEmpty() {
			s := store.GetInstance()
			if err := s.Initialize(dbCfg); err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}

			// Sync global model DB
			if err := model.InitDB(o.sqlDriver); err != nil {
				return err
			}
		}
	}

	// 2. Handle migration/upgrade commands
	if o.migration {
		return migration.Migrate(o.sqlDriver)
	}

	if o.upgrade {
		return migration.Upgrade(o.sqlDriver)
	}

	// 3. Start Application
	return app.NewHTTPSever()
}
