package auth

import "testing"

func TestPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "normal password",
			password: "mypassword123",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "password with spaces",
			password: "my password 123",
		},
		{
			name:     "password with special characters",
			password: "!@#$%^&*()_+-=",
		},
		{
			name:     "very long password",
			password: "this-is-a-very-long-password-that-keeps-going-and-going-and-going-1234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() returned an error: %v", err)
			}

			match, err := CheckPasswordHash(tt.password, hash)
			if err != nil {
				t.Fatalf("CheckPasswordHash() returned an error: %v", err)
			}

			if !match {
				t.Errorf("expected password to match hash")
			}
		})
	}
}

func TestWrongPassword(t *testing.T) {
	password := "correctpassword"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an error: %v", err)
	}

	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an error: %v", err)
	}

	if match {
		t.Errorf("expected wrong password not to match hash")
	}
}

func TestInvalidHash(t *testing.T) {
	password := "mypassword123"
	invalidHash := "this-is-not-a-valid-hash"

	_, err := CheckPasswordHash(password, invalidHash)

	if err == nil {
		t.Errorf("expected an error for invalid hash")
	}
}
