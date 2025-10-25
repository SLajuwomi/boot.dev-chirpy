package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	user1 := uuid.New()
	token1Secret := "AllYourBase"
	token1Expiry := 2 * time.Minute
	token1, _ := MakeJWT(user1, token1Secret, token1Expiry)

	user2 := uuid.New()
	token2Secret := "BrokenDogs"
	token2Expiry := 1 * time.Second
	token2, _ := MakeJWT(user2, token2Secret, token2Expiry)

	tests := []struct {
		name        string
		token       string
		tokenSecret string
		wantErr     bool
		user        uuid.UUID
	}{
		{
			name:        "Correct JWT",
			token:       token1,
			tokenSecret: token1Secret,
			wantErr:     false,
			user:        user1,
		},
		{
			name:        "Expired token",
			token:       token2,
			tokenSecret: token2Secret,
			wantErr:     true,
			user:        user2,
		},
		{
			name:        "Wrong secret",
			token:       token1,
			tokenSecret: token2Secret,
			wantErr:     true,
			user:        user1,
		},
	}

	time.Sleep(5 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userFromValidateJWT, err := ValidateJWT(tt.token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userFromValidateJWT != tt.user {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.user, userFromValidateJWT)
			}
		})
	}
}
