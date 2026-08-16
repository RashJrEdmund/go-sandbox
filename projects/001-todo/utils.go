package main

import (
	"errors"
	"fmt"
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

func removeTag(args *[]string, tag string) {
	if slices.Contains(*args, tag) {
		*args = slices.Delete(*args, slices.Index(*args, tag), slices.Index(*args, tag)+1)
	}
}

func ParseInput(args []string) (command string, input string, tags []string) { // input could be a title or an ID
	if len(args) < 1 {
		return command, input, tags // zero values
	}

	fmt.Println(args, len(args))

	if slices.Contains(args, "-p") || slices.Contains(args, "--print") {
		tags = append(tags, "print")
		removeTag(&args, "-p")
		removeTag(&args, "--print")
	}

	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		tags = append(tags, "help")
		removeTag(&args, "-h")
		removeTag(&args, "--help")

		return command, input, tags
	}

	if len(args) >= 1 {
		command = args[0]
	}

	if len(args) >= 2 {
		input = args[1]
	}

	return command, input, tags
}
