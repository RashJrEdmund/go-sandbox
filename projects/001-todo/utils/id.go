package utils

import (
	"errors"
	"math/rand"
	"slices"
	"sync"
)

var (
	idStore = []string{}
	mu      sync.RWMutex
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

func GenerateUniqueId() (string, error) {
	var newId string

	mu.RLock()
	for range maxRetries {
		newId = generateShortCode()
		if !slices.Contains(idStore, newId) {
			break // found an unused unique Id so stopping the loop
		}
	}
	mu.RUnlock()

	if newId == "" {
		return "", errors.New("Unique Id could not be generated")
	}

	mu.Lock()
	idStore = append(idStore, newId) // reserves Id while we still hold the lock
	mu.Unlock()

	return newId, nil
}
