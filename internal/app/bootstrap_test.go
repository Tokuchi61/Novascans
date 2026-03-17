package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func TestBootstrapRegistersIdentityAuthModule(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/ping", nil)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body struct {
		Data struct {
			Module string `json:"module"`
			Status string `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Data.Module != "identity.auth" {
		t.Fatalf("expected module identity.auth, got %q", body.Data.Module)
	}

	if body.Data.Status != "ok" {
		t.Fatalf("expected status ok, got %q", body.Data.Status)
	}
}

func TestBootstrapRegistersHealthzRoute(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestBootstrapRegistersReadyzRoute(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestBootstrapRegistersMetricsRoute(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
