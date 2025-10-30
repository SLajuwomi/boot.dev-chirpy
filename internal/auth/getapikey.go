package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authContent := headers.Get("Authorization")
	if authContent != "" {
		splitted := strings.Split(authContent, " ")
		if splitted[0] != "ApiKey" {
			return "", fmt.Errorf("wrong prefix. expected %v, got %v", "ApiKey", splitted[0])
		}
		if len(splitted) > 2 {
			return "", fmt.Errorf("extra spaces or invalid token")
		}
		return splitted[1], nil
	}
	return "", fmt.Errorf("failed to get api key")
}
