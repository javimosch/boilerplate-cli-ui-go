package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempPIDFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "boilerplate-pid-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp PID file: %v", err)
	}
	f.Close()
	if content != "" {
		if err := os.WriteFile(f.Name(), []byte(content), 0644); err != nil {
			os.Remove(f.Name())
			t.Fatalf("failed to write temp PID file: %v", err)
		}
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestGetExecutablePath(t *testing.T) {
	path, err := getExecutablePath()
	if err != nil {
		t.Fatalf("getExecutablePath() returned error: %v", err)
	}
	if path == "" {
		t.Fatal("getExecutablePath() returned empty path")
	}
	if !filepath.IsAbs(path) {
		t.Fatalf("getExecutablePath() returned relative path: %s", path)
	}
}

func TestIsDaemonRunningNoPIDFile(t *testing.T) {
	os.Remove(pidFile)
	if isDaemonRunning() {
		t.Fatal("isDaemonRunning() should return false when no PID file exists")
	}
}

func TestIsDaemonRunningInvalidPID(t *testing.T) {
	os.Remove(pidFile)
	if err := os.WriteFile(pidFile, []byte("invalid"), 0644); err != nil {
		t.Fatalf("failed to write test PID file: %v", err)
	}
	defer os.Remove(pidFile)

	if isDaemonRunning() {
		t.Fatal("isDaemonRunning() should return false for invalid PID file content")
	}
}

func TestIsDaemonRunningZeroPID(t *testing.T) {
	os.Remove(pidFile)
	if err := os.WriteFile(pidFile, []byte("0"), 0644); err != nil {
		t.Fatalf("failed to write test PID file: %v", err)
	}
	defer os.Remove(pidFile)

	if isDaemonRunning() {
		t.Fatal("isDaemonRunning() should return false for PID 0")
	}
}

func TestIsDaemonRunningNegativePID(t *testing.T) {
	os.Remove(pidFile)
	if err := os.WriteFile(pidFile, []byte("-1"), 0644); err != nil {
		t.Fatalf("failed to write test PID file: %v", err)
	}
	defer os.Remove(pidFile)

	if isDaemonRunning() {
		t.Fatal("isDaemonRunning() should return false for negative PID")
	}
}

func TestCheckDaemonStatusNotRunning(t *testing.T) {
	os.Remove(pidFile)

	// Test --human flag output
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"program", "status", "--human"}

	output := captureStdout(func() { checkDaemonStatus() })
	if !strings.Contains(output, "not running") {
		t.Fatalf("checkDaemonStatus --human should report 'not running', got: %s", output)
	}
}

func TestCheckDaemonStatusJSON(t *testing.T) {
	// Test JSON output by default (agent-first design)
	os.Remove(pidFile)

	output := captureStdout(func() { checkDaemonStatus() })

	var status DaemonStatus
	if err := json.Unmarshal([]byte(output), &status); err != nil {
		t.Fatalf("checkDaemonStatus should output valid JSON by default, got: %s", output)
	}

	if status.Running != false {
		t.Fatal("expected Running to be false when daemon not running")
	}
	if status.PID != 0 {
		t.Fatalf("expected PID to be 0 when daemon not running, got %d", status.PID)
	}
}

func TestCheckDaemonStatusJSONRunning(t *testing.T) {
	os.Remove(pidFile)
	myPID := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", myPID)), 0644); err != nil {
		t.Fatalf("failed to write PID file: %v", err)
	}
	defer os.Remove(pidFile)

	output := captureStdout(func() { checkDaemonStatus() })

	var status DaemonStatus
	if err := json.Unmarshal([]byte(output), &status); err != nil {
		t.Fatalf("checkDaemonStatus should output valid JSON, got: %s", output)
	}

	if status.Running != true {
		t.Fatal("expected Running to be true when daemon running")
	}
	if status.PID != myPID {
		t.Fatalf("expected PID to be %d, got %d", myPID, status.PID)
	}
	if status.LogFile != logFile {
		t.Fatalf("expected LogFile to be %s, got %s", logFile, status.LogFile)
	}
}

func TestStopDaemonNoPIDFile(t *testing.T) {
	os.Remove(pidFile)

	output := captureStdout(func() { stopDaemon() })
	if !strings.Contains(output, "not running") {
		t.Fatalf("stopDaemon should report 'not running', got: %s", output)
	}
}

func TestStopDaemonStalePID(t *testing.T) {
	os.Remove(pidFile)
	// Use a PID that is very unlikely to exist
	stalePID := 999999999
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", stalePID)), 0644); err != nil {
		t.Fatalf("failed to write test PID file: %v", err)
	}
	defer os.Remove(pidFile)

	output := captureStdout(func() { stopDaemon() })
	if !strings.Contains(output, "Daemon stopped") {
		t.Fatalf("stopDaemon should report 'Daemon stopped', got: %s", output)
	}
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Fatal("stopDaemon should remove stale PID file")
	}
}

func TestIsDaemonRunningCleansStalePID(t *testing.T) {
	os.Remove(pidFile)
	stalePID := 999999999
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", stalePID)), 0644); err != nil {
		t.Fatalf("failed to write test PID file: %v", err)
	}
	defer os.Remove(pidFile)

	if isDaemonRunning() {
		t.Fatal("isDaemonRunning() should return false for stale PID")
	}
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Fatal("isDaemonRunning() should remove stale PID file")
	}
}

func TestStartDaemonAlreadyRunning(t *testing.T) {
	os.Remove(pidFile)
	myPID := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", myPID)), 0644); err != nil {
		t.Fatalf("failed to write PID file: %v", err)
	}
	defer os.Remove(pidFile)

	output := captureStdout(func() {
		startDaemon(8080)
	})
	if !strings.Contains(output, "already running") {
		t.Fatalf("startDaemon should report 'already running', got: %s", output)
	}
}

func TestReadPIDFile_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantPID int
		wantErr bool
	}{
		{"non-existent file", "", 0, true},
		{"empty file", "", 0, true},
		{"invalid string", "not-a-number", 0, true},
		{"negative PID", "-1", 0, true},
		{"zero PID", "0", 0, true},
		{"valid PID", "12345", 12345, false},
		{"PID with whitespace", "  67890  ", 67890, false},
		{"PID with newline", "9999\n", 9999, false},
		{"multiple numbers", "123 456", 0, true},
		{"non-numeric suffix", "123abc", 0, true},
		{"leading zeros", "00123", 0, true},
		{"large PID", "2147483647", 2147483647, false},
		{"too large PID", "9223372036854775808", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tempPIDFile(t, tt.content)

			got, err := readPIDFile(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("readPIDFile() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.wantPID {
				t.Fatalf("readPIDFile() = %d, want %d", got, tt.wantPID)
			}
		})
	}
}


