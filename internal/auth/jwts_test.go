package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret-key"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned an error: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf(
			"ValidateJWT() returned user ID %v, expected %v",
			validatedUserID,
			userID,
		)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "super-secret-key"

	// A negative duration means the expiration time
	// will already be in the past.
	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected expired JWT to be rejected")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"

	token, err := MakeJWT(userID, correctSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() returned an error: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Errorf("expected JWT signed with wrong secret to be rejected")
	}
}
