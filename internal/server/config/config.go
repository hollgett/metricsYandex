package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

var Config struct {
	Addr            string `env:"ADDRESS"`
	StorageInterval int    `env:"STORE_INTERVAL"`
	PathFileStorage string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func InitConfig() error {
	flag.StringVar(&Config.Addr, "a", "localhost:8080", "server")
	flag.IntVar(&Config.StorageInterval, "i", 300, "storage pull interval")
	flag.StringVar(&Config.PathFileStorage, "f", "", "path to temp file")
	flag.BoolVar(&Config.Restore, "r", true, "flag, load old save data")
	flag.Parse()
	if err := env.Parse(&Config); err != nil {
		return err
	}
	return nil
}
