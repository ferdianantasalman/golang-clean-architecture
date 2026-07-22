package util

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	JTI    string    `json:"jti"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, name string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}

	expiryMin := 15
	if v, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MIN")); err == nil && v > 0 {
		expiryMin = v
	}

	claims := Claims{
		UserID: userID,
		Name:   name,
		JTI:    uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func Sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}

func RefreshExpiryDuration() time.Duration {
	hours := 168
	if v, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_HOUR")); err == nil && v > 0 {
		hours = v
	}
	return time.Duration(hours) * time.Hour
}
