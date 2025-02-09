package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

var AgentConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func InitConfig() error {
	flag.StringVar(&AgentConfig.Addr, "a", "localhost:8080", "host:port")
	flag.IntVar(&AgentConfig.ReportInterval, "r", 10, "frequency request")
	flag.IntVar(&AgentConfig.PollInterval, "p", 2, "frequency poll")
	flag.Parse()
	if err := env.Parse(&AgentConfig); err != nil {
		return err
	}
	return nil
}
