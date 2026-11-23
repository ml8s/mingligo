package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("BASE_PATH", "/mingligo/api")
	server, basePath, err := newServer([]byte(`{"openapi":"3.1.0","paths":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	if basePath != "/mingligo/api" {
		t.Fatalf("base path = %q", basePath)
	}
	return httptest.NewServer(server)
}

func TestAPIContractErrors(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		contentType  string
		body         string
		wantStatus   int
		wantType     string
		wantResponse string
	}{
		{name: "openapi", method: "GET", path: "/mingligo/api/openapi.json", wantStatus: 200},
		{name: "health is operational only", method: "GET", path: "/healthz", wantStatus: 200},
		{name: "unknown method", method: "POST", path: "/mingligo/api/not_a_method", contentType: "application/json", body: "{}", wantStatus: 404, wantType: "urn:mingligo:problem:method-not-found"},
		{name: "wrong method", method: "GET", path: "/mingligo/api/time_now", wantStatus: 405, wantType: "urn:mingligo:problem:method-not-allowed"},
		{name: "wrong content type", method: "POST", path: "/mingligo/api/time_now", contentType: "text/plain", body: "{}", wantStatus: 415, wantType: "urn:mingligo:problem:unsupported-media-type"},
		{name: "invalid json", method: "POST", path: "/mingligo/api/time_now", contentType: "application/json", body: "{", wantStatus: 400, wantType: "urn:mingligo:problem:invalid-json"},
		{name: "invalid input", method: "POST", path: "/mingligo/api/time_now", contentType: "application/json", body: "{}", wantStatus: 422, wantType: "urn:mingligo:problem:invalid-input"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, srv.URL+tt.path, bytes.NewReader([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			req.Header.Set("X-Request-Id", "test-request-id")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
			if resp.Header.Get("X-Request-Id") != "test-request-id" {
				t.Fatalf("request id = %q", resp.Header.Get("X-Request-Id"))
			}
			if tt.wantType == "" {
				return
			}
			if got := resp.Header.Get("Content-Type"); got != "application/problem+json; charset=utf-8" {
				t.Fatalf("content type = %q", got)
			}
			var got struct {
				Type      string `json:"type"`
				Status    int    `json:"status"`
				Instance  string `json:"instance"`
				RequestID string `json:"request_id"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if got.Type != tt.wantType || got.Status != tt.wantStatus {
				t.Fatalf("problem = %+v", got)
			}
			if got.Instance != tt.path || got.RequestID != "test-request-id" {
				t.Fatalf("problem metadata = %+v", got)
			}
		})
	}
}

func TestSuccessfulResponseContract(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/mingligo/api/time_now", "application/json", bytes.NewReader([]byte(`{"timezone":"Asia/Shanghai"}`)))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got["local_standard_time"] == "" {
		t.Fatalf("missing domain view field: %+v", got)
	}
}
