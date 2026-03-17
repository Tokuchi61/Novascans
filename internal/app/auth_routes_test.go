package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func TestAuthRoutesSupportUserAndSessionFlow(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	registerBody := []byte(`{"email":"user@example.com","password":"supersecret"}`)
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	registerRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(registerRec, registerReq)

	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected register status %d, got %d", http.StatusCreated, registerRec.Code)
	}

	var registerResp struct {
		Data struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
			Session struct {
				ID string `json:"id"`
			} `json:"session"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(registerRec.Body).Decode(&registerResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	if registerResp.Data.AccessToken == "" || registerResp.Data.RefreshToken == "" {
		t.Fatal("expected access and refresh token to be returned")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+registerResp.Data.AccessToken)
	meRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("expected me status %d, got %d", http.StatusOK, meRec.Code)
	}

	refreshBody := []byte(`{"refresh_token":"` + registerResp.Data.RefreshToken + `"}`)
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(refreshBody))
	refreshRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(refreshRec, refreshReq)

	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected refresh status %d, got %d", http.StatusOK, refreshRec.Code)
	}

	var refreshResp struct {
		Data struct {
			Session struct {
				ID string `json:"id"`
			} `json:"session"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(refreshRec.Body).Decode(&refreshResp); err != nil {
		t.Fatalf("decode refresh response: %v", err)
	}

	if refreshResp.Data.AccessToken == "" || refreshResp.Data.RefreshToken == "" {
		t.Fatal("expected refreshed access and refresh token to be returned")
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+refreshResp.Data.AccessToken)
	logoutRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusNoContent {
		t.Fatalf("expected logout status %d, got %d", http.StatusNoContent, logoutRec.Code)
	}
}

func TestAuthRoutesReturnValidationErrorEnvelope(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{"email":"invalid","password":"123"}`)))
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected validation status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if response.Error.Code != "validation_error" {
		t.Fatalf("expected validation_error, got %q", response.Error.Code)
	}
}

func TestRouterReturnsStandardNotFoundEnvelope(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/unknown", nil)
	rec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var response struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if response.Error.Code != "not_found" {
		t.Fatalf("expected not_found, got %q", response.Error.Code)
	}
}
