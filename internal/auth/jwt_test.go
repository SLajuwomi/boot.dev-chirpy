package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		expectedToken string
		wantErr       bool
	}{
		{
			name:          "Correct auth string",
			token:         "Bearer 3jfh56isme019vmsu",
			expectedToken: "3jfh56isme019vmsu",
			wantErr:       false,
		},
		{
			name:          "Missing header",
			token:         "",
			expectedToken: "",
			wantErr:       true,
		},
		{
			name:          "Wrong prefix",
			token:         "Brock 3jfh56isme019vmsu",
			expectedToken: "",
			wantErr:       true,
		},
		{
			name:          "Extra spaces",
			token:         "Bearer    3jfh56isme019vmsu",
			expectedToken: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newHeader := make(http.Header)
			if tt.name != "Missing header" {
				newHeader.Set("Authorization", tt.token)
			}

			bearerToken, err := GetBearerToken(newHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error: %v, wantErr: %v", err, tt.wantErr)
			}
			if !tt.wantErr && bearerToken != tt.expectedToken {
				t.Errorf("GetBearerToken() expects %v, got %v", tt.expectedToken, bearerToken)
			}
		})
	}
}
