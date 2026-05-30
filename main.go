package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const Version = "1.0.0"

// JSON output structures for agent-first design
type VersionInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DaemonStatus struct {
	Running bool   `json:"running"`
	PID     int    `json:"pid,omitempty"`
	LogFile string `json:"log_file,omitempty"`
}

// Agent-first design semantic exit codes
const (
	ExitSuccess         = 0  // Success
	ExitUserError       = 80 // User errors (invalid arguments, missing required input)
	ExitUserInvalidFlag = 81 // Invalid command line flags
	ExitResourceError   = 90 // Resource errors (file access, network, etc.)
	ExitIntegrationError = 100 // Integration errors (external service failures)
	ExitSoftwareError   = 110 // Software errors (internal bugs, unexpected conditions)
)

func runCommand(args []string) int {
	if len(args) < 1 {
		printHelp()
		return ExitUserError
	}

	command := args[0]

	switch command {
	case "start":
		return handleStart(args[1:])
	case "stop":
		return handleStop()
	case "status":
		return handleStatus()
	case "version":
		handleVersion()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printHelp()
		return ExitUserError
	}
	return ExitSuccess
}

func main() {
	os.Exit(runCommand(os.Args[1:]))
}

func handleStart(args []string) int {
	startCmd := flag.NewFlagSet("start", flag.ContinueOnError)
	port := startCmd.Int("port", 8080, "Port for HTTP server")
	daemon := startCmd.Bool("daemon", false, "Run as daemon")
	if err := startCmd.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "start: %v\n", err)
		return ExitUserInvalidFlag
	}

	if *daemon {
		return startDaemon(*port)
	} else {
		return startServer(*port)
	}
}

func handleStop() int {
	return stopDaemon()
}

func handleStatus() int {
	return checkDaemonStatus()
}

func handleVersion() {
	versionInfo := VersionInfo{
		Name:    "boilerplate-cli-ui-go",
		Version: Version,
	}

	// Check for --human flag (agent-first design: JSON by default)
	if len(os.Args) > 2 && os.Args[2] == "--human" {
		fmt.Printf("boilerplate-cli-ui-go v%s\n", Version)
	} else {
		json.NewEncoder(os.Stdout).Encode(versionInfo)
	}
}

func printHelp() {
	fmt.Println("boilerplate-cli-ui-go - Go CLI with HTTP UI and daemon management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  boilerplate-cli-ui-go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  start       Start HTTP server (UI)")
	fmt.Println("  stop        Stop daemon server")
	fmt.Println("  status      Check daemon status")
	fmt.Println("  version     Show version information")
	fmt.Println("  help        Show this help message")
	fmt.Println()
	fmt.Println("Start Options:")
	fmt.Println("  -port int   Port for HTTP server (default 8080)")
	fmt.Println("  -daemon     Run as daemon (background)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  boilerplate-cli-ui-go start")
	fmt.Println("  boilerplate-cli-ui-go start -port 3000")
	fmt.Println("  boilerplate-cli-ui-go start -daemon")
	fmt.Println("  boilerplate-cli-ui-go start -port 3000 -daemon")
	fmt.Println("  boilerplate-cli-ui-go stop")
	fmt.Println("  boilerplate-cli-ui-go status")
	fmt.Println("  boilerplate-cli-ui-go version")
}