package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func TestAccessMeResolvesPrincipalFromAuthToken(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{"email":"user@example.com","password":"supersecret"}`)))
	registerRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(registerRec, registerReq)

	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected register status %d, got %d", http.StatusCreated, registerRec.Code)
	}

	var registerResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(registerRec.Body).Decode(&registerResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/access/me", nil)
	req.Header.Set("Authorization", "Bearer "+registerResp.Data.AccessToken)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected access me status %d, got %d", http.StatusOK, rec.Code)
	}

	var response struct {
		Data struct {
			IsGuest  bool   `json:"is_guest"`
			BaseRole string `json:"base_role"`
			Email    string `json:"email"`
		} `json:"data"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode access me response: %v", err)
	}

	if response.Data.IsGuest {
		t.Fatal("expected authenticated principal, got guest")
	}

	if response.Data.BaseRole != "user" {
		t.Fatalf("expected base role user, got %q", response.Data.BaseRole)
	}

	if response.Data.Email != "user@example.com" {
		t.Fatalf("expected user email, got %q", response.Data.Email)
	}
}
