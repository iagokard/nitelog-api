package util

import (
	"crypto/rand"
	"encoding/base64"
	"nitelog/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @model ErrorResponse
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// @model MessageResponse
type MessageResponse struct {
	Message string `json:"message" example:"Sample status message"`
}

func NormalizeDate(date time.Time) (*time.Time, error) {
	cfg := config.Load()
	location, err := time.LoadLocation(cfg.Timezone)

	if err != nil {
		return nil, err
	}

	normalizedDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		location,
	).UTC()

	return &normalizedDate, nil
}

func GenerateMeetingCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func GenerateJWT(userID, secret string) (string, error) {
	expirationTime := 24 * time.Hour
	claims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(expirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
