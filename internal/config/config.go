package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

type CommandAddr struct {
	Addr string `env:"ADDRESS"`
}

func validatePort(addr string) error {
	hp := strings.Split(addr, ":")
	if _, err := strconv.Atoi(hp[1]); len(hp) != 2 || err != nil {
		return fmt.Errorf("address must be in the format `:port`, got: %s", addr)
	}
	return nil
}

func InitConfig() *CommandAddr {
	addr := flag.String("a", ":8080", "setup server address host:port")
	commandAddr := &CommandAddr{}
	if err := env.Parse(commandAddr); err != nil {
		panic(err)
	}

	if commandAddr.Addr == "" {
		flag.Parse()
		commandAddr.Addr = *addr
	}

	if err := validatePort(commandAddr.Addr); err != nil {
		panic(err)
	}

	return commandAddr
}
