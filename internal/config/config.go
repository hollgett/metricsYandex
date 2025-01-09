package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type CommandAddr struct {
	Addr string `env:"ADDRESS"`
}

var Cfg = CommandAddr{}

func validatePort(serverCfg *CommandAddr) error {
	parsedURL, err := url.Parse(serverCfg.Addr)
	if err != nil {
		return fmt.Errorf("addr error: %w", err)
	}
	if parsedURL.Scheme == `http` {
		parsedURL.Scheme = ""
		serverCfg.Addr = parsedURL.Host
	}
	return nil
}

func InitConfig() error {
	addr := flag.String("a", "localhost:8080", "setup server address host:port")
	commandAddr := &CommandAddr{}
	if err := env.Parse(commandAddr); err != nil {
		return err
	}

	if commandAddr.Addr == "" {
		flag.Parse()
		commandAddr.Addr = *addr
	}

	if err := validatePort(commandAddr); err != nil {
		return err
	}

	Cfg = *commandAddr
	return nil
}
