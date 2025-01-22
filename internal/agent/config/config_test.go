package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		flag        map[string]string
		wantAddr    string
		wantRepInt  int
		wantPollInt int
	}{
		{"without env & flag", nil, nil, "localhost:8080", 10, 2},
		{"without env", nil, map[string]string{"a": "localhost:8085", "r": "5", "p": "1"}, "localhost:8085", 5, 1},
		{"without flag", map[string]string{"ADDRESS": "localhost:8085", "REPORT_INTERVAL": "5", "POLL_INTERVAL": "1"}, nil, "localhost:8085", 5, 1},
		{"flag & env", map[string]string{"ADDRESS": "localhost:8088", "REPORT_INTERVAL": "8", "POLL_INTERVAL": "7"}, map[string]string{"a": "localhost:8085", "r": "5", "p": "1"}, "localhost:8088", 8, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = []string{"cmd"}
			if len(tt.env) > 0 {
				for name, val := range tt.env {
					require.NoError(t, os.Setenv(name, val))
				}
			}
			if len(tt.flag) > 0 {
				for name, val := range tt.flag {
					os.Args = append(os.Args, "-"+name+"="+val)
				}
			}
			require.NoError(t, InitConfig(), "init config")

			assert.Equal(t, tt.wantAddr, AgentConfig.Addr)
			assert.Equal(t, tt.wantRepInt, AgentConfig.ReportInterval)
			assert.Equal(t, tt.wantPollInt, AgentConfig.PollInterval)
		})
	}
}
