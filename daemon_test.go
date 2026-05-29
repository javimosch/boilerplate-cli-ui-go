package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

	output := captureStdout(checkDaemonStatus)
	if !strings.Contains(output, "not running") {
		t.Fatalf("checkDaemonStatus should report 'not running', got: %s", output)
	}
}

func TestStopDaemonNoPIDFile(t *testing.T) {
	os.Remove(pidFile)

	output := captureStdout(stopDaemon)
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

	output := captureStdout(stopDaemon)
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


