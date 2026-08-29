package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":    "chirpy-acces",
		"IssuedAt":  time.Now(),
		"ExpiresAt": time.Now() + expiresIn,
		"Subject":   phrases(userID),
	})
}
