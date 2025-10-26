package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
}

// var mySigningKey = []byte("AllYourBase")

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("signing failed: %v", err)
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse with claims failed: %v", err)
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		fmt.Println(claims.Issuer)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
	userID, err := uuid.Parse(token.Claims.(*MyCustomClaims).Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("uuid parse failed: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authContent := headers.Get("Authorization")
	if authContent != "" {
		splitted := strings.Split(authContent, " ")
		if splitted[0] != "Bearer" {
			return "", fmt.Errorf("wrong prefix. expected %v, got %v", "Bearer", splitted[0])
		}
		if len(splitted) > 2 {
			return "", fmt.Errorf("extra spaces or invalid token")
		}
		return splitted[1], nil
	}
	return "", fmt.Errorf("failed to get bearer token")
}
