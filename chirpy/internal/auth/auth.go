package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expiresIn := time.Duration(60*60) * time.Second // 1 hour

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(token.Claims.(*jwt.RegisteredClaims).Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// Reading the Authorization header: Bearer TOKEN_STRING
func GetBearerToken(headers http.Header) (string, error) {
	parts := strings.Fields(headers.Get("Authorization"))
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}

	if parts[1] == "" {
		return "", fmt.Errorf("no token provided")
	}

	return parts[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	parts := strings.Fields(headers.Get("Authorization"))
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid authorization header")
	}

	if parts[1] == "" {
		return "", fmt.Errorf("no api key provided")
	}

	return parts[1], nil
}

func GetUserIdFromBearerTokenHeader(h http.Header, jwtSecret string) (uuid.UUID, error) {
	bearerToken, err := GetBearerToken(h)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userIDFromToken, err := ValidateJWT(bearerToken, jwtSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return userIDFromToken, nil
}

// REFRESH TOKEN

/*
For the creation of this function and links to other resources,
see the https://www.boot.dev/lessons/f7285cef-5185-4b15-b5fc-9533ccaafe8a course
*/
func MakeRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
