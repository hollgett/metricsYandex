package main

import (
	"flag"
	"strconv"
	"strings"
)

type AgentArgs struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func InitAgentCommand() *AgentArgs {
	addr := flag.String("a", "localhost:8080", "host:port")
	reportInt := flag.Int("r", 10, "frequency request")
	pollInt := flag.Int("p", 2, "frequency poll")

	flag.Parse()
	hp := strings.Split(*addr, ":")
	if _, err := strconv.Atoi(hp[1]); len(hp) != 2 || err != nil {
		panic("command error address")
	}

	return &AgentArgs{
		Addr:           *addr,
		ReportInterval: *reportInt,
		PollInterval:   *pollInt,
	}
}
