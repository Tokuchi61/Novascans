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

	createUserBody := []byte(`{"email":"user@example.com","password":"supersecret"}`)
	createUserReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/users", bytes.NewReader(createUserBody))
	createUserRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(createUserRec, createUserReq)

	if createUserRec.Code != http.StatusCreated {
		t.Fatalf("expected create user status %d, got %d", http.StatusCreated, createUserRec.Code)
	}

	var createUserResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(createUserRec.Body).Decode(&createUserResp); err != nil {
		t.Fatalf("decode create user response: %v", err)
	}

	getUserReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/users/"+createUserResp.Data.ID, nil)
	getUserRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(getUserRec, getUserReq)

	if getUserRec.Code != http.StatusOK {
		t.Fatalf("expected get user status %d, got %d", http.StatusOK, getUserRec.Code)
	}

	createSessionBody := []byte(`{"email":"user@example.com","password":"supersecret"}`)
	createSessionReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sessions", bytes.NewReader(createSessionBody))
	createSessionRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(createSessionRec, createSessionReq)

	if createSessionRec.Code != http.StatusCreated {
		t.Fatalf("expected create session status %d, got %d", http.StatusCreated, createSessionRec.Code)
	}

	var createSessionResp struct {
		Data struct {
			ID    string `json:"id"`
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(createSessionRec.Body).Decode(&createSessionResp); err != nil {
		t.Fatalf("decode create session response: %v", err)
	}

	if createSessionResp.Data.Token == "" {
		t.Fatal("expected session token to be returned")
	}

	revokeReq := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/sessions/"+createSessionResp.Data.ID, nil)
	revokeRec := httptest.NewRecorder()

	runtime.Router.ServeHTTP(revokeRec, revokeReq)

	if revokeRec.Code != http.StatusNoContent {
		t.Fatalf("expected revoke session status %d, got %d", http.StatusNoContent, revokeRec.Code)
	}
}

func TestAuthRoutesReturnValidationErrorEnvelope(t *testing.T) {
	runtime, err := Bootstrap(t.Context(), config.Default(), WithSkipDB())
	if err != nil {
		t.Fatalf("bootstrap returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/users", bytes.NewReader([]byte(`{"email":"invalid","password":"123"}`)))
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
