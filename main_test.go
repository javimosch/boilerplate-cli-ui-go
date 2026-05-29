package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old
	return <-out
}

func captureStderr(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	f()

	w.Close()
	os.Stderr = old
	return <-out
}

func TestHandleVersion(t *testing.T) {
	output := captureStdout(handleVersion)
	if !strings.Contains(output, "boilerplate-cli-ui-go v1.0.0") {
		t.Fatalf("handleVersion output should contain version, got: %s", output)
	}
}

func TestPrintHelp(t *testing.T) {
	output := captureStdout(printHelp)
	if !strings.Contains(output, "Usage:") {
		t.Fatalf("printHelp should contain Usage:, got: %s", output)
	}
	if !strings.Contains(output, "start") {
		t.Fatalf("printHelp should mention start command, got: %s", output)
	}
}

func TestRunCommand_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantExit  int
		wantOut   string
		wantErr   string
	}{
		{
			name:     "version command",
			args:     []string{"version"},
			wantExit: 0,
			wantOut:  "boilerplate-cli-ui-go v1.0.0",
		},
		{
			name:     "help command",
			args:     []string{"help"},
			wantExit: 0,
			wantOut:  "Usage:",
		},
		{
			name:     "--help flag",
			args:     []string{"--help"},
			wantExit: 0,
			wantOut:  "Usage:",
		},
		{
			name:     "-h flag",
			args:     []string{"-h"},
			wantExit: 0,
			wantOut:  "Usage:",
		},
		{
			name:     "no args",
			args:     []string{},
			wantExit: 1,
			wantOut:  "Usage:",
		},
		{
			name:     "unknown command",
			args:     []string{"foobar"},
			wantExit: 1,
			wantOut:  "Usage:",
			wantErr:  "Unknown command: foobar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := captureStdout(func() {
				stderr := captureStderr(func() {
					got := runCommand(tt.args)
					if got != tt.wantExit {
						t.Fatalf("runCommand() exit code = %d, want %d", got, tt.wantExit)
					}
				})
				if tt.wantErr != "" && !strings.Contains(stderr, tt.wantErr) {
					t.Fatalf("runCommand() stderr should contain %q, got: %s", tt.wantErr, stderr)
				}
			})
			if !strings.Contains(stdout, tt.wantOut) {
				t.Fatalf("runCommand() stdout should contain %q, got: %s", tt.wantOut, stdout)
			}
		})
	}
}
