package utils

import "github.com/google/uuid"

func GenerateAccessToken() string {
	return uuid.New().String()
}
