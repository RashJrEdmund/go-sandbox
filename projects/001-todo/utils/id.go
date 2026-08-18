package utils

import (
	"errors"
	"math/rand"
	"slices"
)

const (
	Delimiter  = "---------------------------------------------------------------------------------------------------"
	maxRetries = 10
	idLength   = 5
)

func generateShortCode() string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	bytes := make([]byte, idLength)
	for i := range bytes {
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(bytes)
}

func GenerateUniqueId(existingIds []string) (string, error) {
	var newId string
	foundUnique := false

	for range maxRetries {
		newId = generateShortCode()
		if !slices.Contains(existingIds, newId) {
			foundUnique = true
			break // found an unused unique Id so stopping the loop
		}
	}

	if !foundUnique {
		return "", errors.New("Unique Id could not be generated")
	}

	return newId, nil
}
