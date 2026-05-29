package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	handleHealthAPI(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if body["status"] != "healthy" {
		t.Fatalf("expected status 'healthy', got '%s'", body["status"])
	}
}

func TestStatusAPI(t *testing.T) {
	mu.Lock()
	serverStatus = Status{
		Status:    "running",
		Port:      9999,
		StartTime: time.Now(),
	}
	mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handleStatusAPI(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var body Status
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if body.Status != "running" {
		t.Fatalf("expected Status 'running', got '%s'", body.Status)
	}
	if body.Port != 9999 {
		t.Fatalf("expected Port 9999, got %d", body.Port)
	}
	if body.Uptime == "" {
		t.Fatal("expected non-empty Uptime")
	}
}

func TestHealthAPI_PostMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()

	handleHealthAPI(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for POST, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if body["status"] != "healthy" {
		t.Fatalf("expected status 'healthy', got '%s'", body["status"])
	}
}

func TestStatusAPI_JSONFields(t *testing.T) {
	mu.Lock()
	serverStatus = Status{
		Status:    "running",
		Port:      8080,
		StartTime: time.Now(),
	}
	mu.Unlock()

	w := httptest.NewRecorder()
	handleStatusAPI(w, httptest.NewRequest(http.MethodGet, "/api/status", nil))

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Result().Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	expectedFields := []string{"status", "port", "uptime", "start_time"}
	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Fatalf("Status JSON missing field: %s", field)
		}
	}
}

func TestHomePage(t *testing.T) {
	mu.Lock()
	serverStatus = Status{
		Status:    "running",
		Port:      3030,
		StartTime: time.Now(),
	}
	mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handleHome(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" {
		t.Fatalf("expected Content-Type text/html, got %s", ct)
	}

	body := w.Body.String()
	if !strings.Contains(body, "3030") {
		t.Fatal("home page should display the current port (3030)")
	}
	if !strings.Contains(body, "Boilerplate CLI UI") {
		t.Fatal("home page should contain the title")
	}
}

func TestHomePage_ContainsVersion(t *testing.T) {
	mu.Lock()
	serverStatus = Status{
		Status:    "running",
		Port:      8080,
		StartTime: time.Now(),
	}
	mu.Unlock()

	w := httptest.NewRecorder()
	handleHome(w, httptest.NewRequest(http.MethodGet, "/", nil))

	body := w.Body.String()
	if !strings.Contains(body, Version) {
		t.Fatalf("home page should contain Version constant %q", Version)
	}
}
