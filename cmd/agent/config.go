package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

type AgentArgs struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func validatePort(addr string) error {
	hp := strings.Split(addr, ":")
	if _, err := strconv.Atoi(hp[1]); len(hp) != 2 || err != nil {
		return fmt.Errorf("address must be in the format `:port`, got: %s", addr)
	}
	return nil
}

func initAgentCommand(agentConfig *AgentArgs) {
	addr := flag.String("a", "localhost:8080", "host:port")
	reportInt := flag.Int("r", 10, "frequency request")
	pollInt := flag.Int("p", 2, "frequency poll")

	flag.Parse()
	if agentConfig.Addr == "" {
		agentConfig.Addr = *addr
	}
	if agentConfig.ReportInterval == 0 {
		agentConfig.ReportInterval = *reportInt
	}
	if agentConfig.PollInterval == 0 {
		agentConfig.PollInterval = *pollInt
	}

}

func InitConfig() *AgentArgs {
	agentArgs := &AgentArgs{}
	if err := env.Parse(agentArgs); err != nil {
		panic(err)
	}
	if agentArgs.Addr == "" || agentArgs.ReportInterval == 0 || agentArgs.PollInterval == 0 {
		initAgentCommand(agentArgs)
	}
	if err := validatePort(agentArgs.Addr); err != nil {
		panic(err)
	}

	return agentArgs
}
