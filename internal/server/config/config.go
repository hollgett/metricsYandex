package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr            string `env:"ADDRESS"`
	StorageInterval int    `env:"STORE_INTERVAL"`
	PathFileStorage string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
}

func New() (*Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server")
	flag.IntVar(&cfg.StorageInterval, "i", 300, "storage pull interval")
	flag.StringVar(&cfg.PathFileStorage, "f", "", "path to temp file")
	flag.BoolVar(&cfg.Restore, "r", true, "flag, load old save data")
	flag.StringVar(&cfg.DataBaseDSN, "d", "", "database DSN")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
