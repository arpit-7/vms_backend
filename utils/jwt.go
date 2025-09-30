package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-auth/config"
	"go-auth/models"
)

// JWT structures
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTPayload struct {
	Sub      string `json:"sub"`
	ID       string `json:"id"`
	Username string `json:"username"`
	GroupId  int    `json:"groupId"`
	AreaName string `json:"areaName"`
	Role     string `json:"role"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

// CreateSessionJWT creates a NextAuth compatible JWT for sessions
func CreateSessionJWT(user models.User) (string, error) {
	now := time.Now().Unix()
	expiresAt := now + (60 * 60 * 24 * 365 * config.TOKEN_EXPIRY_YEARS)

	payload := map[string]interface{}{
		"sub":      fmt.Sprintf("%d", user.ID),
		"id":       fmt.Sprintf("%d", user.ID),
		"username": user.Username,
		"groupId":  user.GroupId,
		"areaName": user.AreaName,
		"role":     user.Role,
		"iat":      now,
		"exp":      expiresAt,
	}

	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	message := headerB64 + "." + payloadB64
	signature := CreateHMACSignature(message, config.SESSION_SECRET)

	token := message + "." + signature

	return token, nil
}

// CreateTokenJWT creates a JWT for token-based authentication
func CreateTokenJWT(user models.User, expiresAt time.Time) (string, error) {
	payload := JWTPayload{
		Sub:      fmt.Sprintf("%d", user.ID),
		ID:       fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		GroupId:  user.GroupId,
		AreaName: user.AreaName,
		Role:     user.Role,
		Iat:      time.Now().Unix(),
		Exp:      expiresAt.Unix(),
	}

	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	message := headerB64 + "." + payloadB64
	signature := CreateHMACSignature(message, config.JWT_SECRET)

	token := message + "." + signature

	return token, nil
}

// CreateNextAuthJWT creates a NextAuth compatible JWT from token record
func CreateNextAuthJWT(tokenRecord models.Token) (string, error) {
	now := time.Now().Unix()
	expiresAt := now + (60 * 60 * 24 * 365 * config.TOKEN_EXPIRY_YEARS)

	payload := map[string]interface{}{
		"sub":      fmt.Sprintf("%d", tokenRecord.UserID),
		"id":       fmt.Sprintf("%d", tokenRecord.UserID),
		"username": tokenRecord.Username,
		"groupId":  tokenRecord.GroupId,
		"areaName": tokenRecord.AreaName,
		"role":     tokenRecord.Role,
		"iat":      now,
		"exp":      expiresAt,
	}

	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	message := headerB64 + "." + payloadB64
	signature := CreateHMACSignature(message, config.JWT_SECRET)

	token := message + "." + signature

	return token, nil
}

// VerifySessionJWT verifies and parses a session JWT
func VerifySessionJWT(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	signature := parts[2]

	expectedSignature := CreateHMACSignature(message, config.SESSION_SECRET)
	if signature != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	// Check expiry
	if exp, ok := payload["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return payload, nil
}

// VerifyTokenJWT verifies a token JWT signature
func VerifyTokenJWT(token string) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	signature := parts[2]

	expectedSignature := CreateHMACSignature(message, config.JWT_SECRET)

	return signature == expectedSignature, nil
}

// CreateHMACSignature creates an HMAC-SHA256 signature
func CreateHMACSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}