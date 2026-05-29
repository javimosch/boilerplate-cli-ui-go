package main

import (
	"flag"
	"fmt"
	"os"
)

const Version = "1.0.0"

func runCommand(args []string) int {
	if len(args) < 1 {
		printHelp()
		return 1
	}

	command := args[0]

	switch command {
	case "start":
		handleStart(args[1:])
	case "stop":
		handleStop()
	case "status":
		handleStatus()
	case "version":
		handleVersion()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printHelp()
		return 1
	}
	return 0
}

func main() {
	os.Exit(runCommand(os.Args[1:]))
}

func handleStart(args []string) {
	startCmd := flag.NewFlagSet("start", flag.ContinueOnError)
	port := startCmd.Int("port", 8080, "Port for HTTP server")
	daemon := startCmd.Bool("daemon", false, "Run as daemon")
	if err := startCmd.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "start: %v\n", err)
		os.Exit(2)
	}

	if *daemon {
		startDaemon(*port)
	} else {
		startServer(*port)
	}
}

func handleStop() {
	stopDaemon()
}

func handleStatus() {
	checkDaemonStatus()
}

func handleVersion() {
	fmt.Printf("boilerplate-cli-ui-go v%s\n", Version)
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