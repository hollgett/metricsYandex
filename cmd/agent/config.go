package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type AgentArgs struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

var Cfg AgentArgs

func validatePort(agentCfg *AgentArgs) error {
	parsedURL, err := url.Parse(agentCfg.Addr)
	if err != nil {
		return fmt.Errorf("addr error: %w", err)
	}
	if parsedURL.Scheme != `http` {
		agentCfg.Addr = `http://` + parsedURL.String()
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

func InitConfig() {
	agentArgs := &AgentArgs{}
	if err := env.Parse(agentArgs); err != nil {
		panic(err)
	}
	if agentArgs.Addr == "" || agentArgs.ReportInterval == 0 || agentArgs.PollInterval == 0 {
		initAgentCommand(agentArgs)
	}
	if err := validatePort(agentArgs); err != nil {
		panic(err)
	}
	Cfg = *agentArgs
}
