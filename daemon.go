package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	pidFile = "/tmp/boilerplate-cli-ui-go.pid"
	logFile = "/tmp/boilerplate-cli-ui-go.log"
)

func startDaemon(port int) int {
	// Check if already running
	if isDaemonRunning() {
		fmt.Println("Daemon is already running")
		return ExitSuccess
	}

	// Get the current executable path (resolves symlinks)
	execPath, err := getExecutablePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		return ExitSoftwareError
	}

	// Create command to run server in foreground
	cmd := exec.Command(execPath, "start", fmt.Sprintf("-port=%d", port))
	
	// Set up logging
	logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		return ExitResourceError
	}
	defer logFileHandle.Close()

	cmd.Stdout = logFileHandle
	cmd.Stderr = logFileHandle

	// Start the process
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting daemon: %v\n", err)
		return ExitIntegrationError
	}

	// Write PID file
	pid := cmd.Process.Pid
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing PID file: %v\n", err)
		cmd.Process.Kill()
		return ExitResourceError
	}

	fmt.Printf("Daemon started with PID %d\n", pid)
	fmt.Printf("Logs: %s\n", logFile)
	return ExitSuccess
}

func readPIDFile(path string) (int, error) {
	pidData, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("reading PID file: %w", err)
	}

	pidStr := strings.TrimSpace(string(pidData))
	if pidStr == "" {
		return 0, fmt.Errorf("invalid PID in file %s", path)
	}

	var pid int
	n, err := fmt.Sscanf(pidStr, "%d", &pid)
	if err != nil || n != 1 || pid <= 0 {
		return 0, fmt.Errorf("invalid PID in file %s", path)
	}

	// Ensure the entire string was consumed (no extra characters)
	if fmt.Sprintf("%d", pid) != pidStr {
		return 0, fmt.Errorf("invalid PID in file %s", path)
	}

	return pid, nil
}

func stopDaemon() int {
	pid, err := readPIDFile(pidFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Daemon is not running")
			return ExitSuccess
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Remove(pidFile)
		return ExitResourceError
	}

	// Send SIGTERM to the process
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding process: %v\n", err)
		os.Remove(pidFile)
		return ExitSoftwareError
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		if errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
			os.Remove(pidFile)
			fmt.Printf("Daemon stopped (PID %d)\n", pid)
			return ExitSuccess
		}
		fmt.Fprintf(os.Stderr, "Error stopping process: %v\n", err)
		os.Remove(pidFile)
		return ExitIntegrationError
	}

	// Wait for process to terminate (max 5 seconds)
	timeout := time.After(5 * time.Second)
	forceTimeout := time.After(6 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	forceKilled := false

	for {
		select {
		case <-timeout:
			if !forceKilled {
				process.Signal(syscall.SIGKILL)
				forceKilled = true
			}
		case <-forceTimeout:
			os.Remove(pidFile)
			fmt.Fprintf(os.Stderr, "Warning: daemon PID %d did not stop, removing PID file\n", pid)
			return ExitIntegrationError
		case <-ticker.C:
			if err := process.Signal(syscall.Signal(0)); err != nil {
				// Process has terminated
				os.Remove(pidFile)
				fmt.Printf("Daemon stopped (PID %d)\n", pid)
				return ExitSuccess
			}
		}
	}
}

func checkDaemonStatus() int {
	var status DaemonStatus

	if isDaemonRunning() {
		pid, _ := readPIDFile(pidFile)
		status = DaemonStatus{
			Running: true,
			PID:     pid,
			LogFile: logFile,
		}
	} else {
		status = DaemonStatus{
			Running: false,
		}
	}

	// Check for --human flag (agent-first design: JSON by default)
	if len(os.Args) > 2 && os.Args[2] == "--human" {
		if status.Running {
			fmt.Printf("Daemon is running (PID %d)\n", status.PID)
			fmt.Printf("Logs: %s\n", status.LogFile)
		} else {
			fmt.Println("Daemon is not running")
		}
	} else {
		json.NewEncoder(os.Stdout).Encode(status)
	}

	return ExitSuccess
}

func isDaemonRunning() bool {
	pid, err := readPIDFile(pidFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			os.Remove(pidFile)
		}
		return false
	}

	// Check if process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	if err := process.Signal(syscall.Signal(0)); err != nil {
		// Process not running, clean up PID file
		os.Remove(pidFile)
		return false
	}

	return true
}

func getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Resolve symlinks
	resolvedPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return execPath, nil
	}

	return resolvedPath, nil
}