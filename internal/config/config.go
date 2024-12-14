package config

import (
	"flag"
	"strconv"
	"strings"
)

type CommandAddr struct {
	Addr string
}

func InitConfig() *CommandAddr {
	addr := flag.String("a", "localhost:8080", "setup server address host:port")

	hp := strings.Split(*addr, ":")
	if _, err := strconv.Atoi(hp[1]); len(hp) != 2 || err != nil {
		panic("error value setup server")
	}

	return &CommandAddr{
		Addr: *addr,
	}
}
