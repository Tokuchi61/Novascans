package app

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type accessTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

func buildAccessToken(secret string, claims accessTokenClaims) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal access token claims: %w", err)
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signAccessToken(secret, encodedPayload)

	return encodedPayload + "." + signature, nil
}

func parseAccessToken(secret string, token string) (accessTokenClaims, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 {
		return accessTokenClaims{}, Unauthorized("invalid access token", nil)
	}

	expectedSignature := signAccessToken(secret, parts[0])
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return accessTokenClaims{}, Unauthorized("invalid access token", nil)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return accessTokenClaims{}, Unauthorized("invalid access token", err)
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return accessTokenClaims{}, Unauthorized("invalid access token", err)
	}

	if time.Now().UTC().Unix() >= claims.ExpiresAt {
		return accessTokenClaims{}, Unauthorized("access token expired", nil)
	}

	return claims, nil
}

func signAccessToken(secret string, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func newOpaqueToken() string {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return uuid.NewString()
	}

	return hex.EncodeToString(raw[:])
}

func hashOpaqueToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
