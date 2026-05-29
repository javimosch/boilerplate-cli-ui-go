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
