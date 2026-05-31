package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const subprocessEnv = "TEST_SUBPROCESS_ARGS"

func TestMainCLI(t *testing.T) {
	if args := os.Getenv(subprocessEnv); args != "" {
		os.Args = strings.Fields(args)

		main()

		return
	}

	tests := []struct {
		name        string
		args        string
		wantErr     bool
		wantContain string
	}{
		{name: "no args", args: "gotpl", wantErr: false},
		{name: "help flag", args: "gotpl --help", wantErr: false},
		{name: "version flag", args: "gotpl --version", wantErr: false},
		{
			name:        "unknown command",
			args:        "gotpl nonexistent",
			wantErr:     true,
			wantContain: "unknown command",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=TestMainCLI")

			cmd.Env = append(os.Environ(), subprocessEnv+"="+tc.args)

			out, err := cmd.CombinedOutput()

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got none; output: %s", out)
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v; output: %s", err, out)
			}

			if tc.wantContain != "" && !strings.Contains(string(out), tc.wantContain) {
				t.Errorf("expected output to contain %q, got: %s", tc.wantContain, out)
			}
		})
	}
}
