/*
X POST
https://x.com/orashus/status/2078129654901198971
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

var (
	store = make(map[string]string) // shortCode -> longUrl
	mu    sync.RWMutex
)

const consoleSpace string = "-----------------------------------"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(size int) string {
	bytes := make([]byte, size)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

const maxRetries = 10

func generateUniqueShortCode(size int, longUrl string) (string, error) {
	var code string

	mu.RLock()
	for range maxRetries {
		code = generateShortCode(size)
		if _, exists := store[code]; !exists {
			break // found an unused unique code so stopping the loop
		}
	}
	mu.RUnlock()

	if code == "" {
		return "", errors.New("Unique Code could not be generated")
	}

	mu.Lock()
	store[code] = longUrl // reserves shortcode while we still hold the lock
	mu.Unlock()

	return code, nil
}

// POST /shorten?url=https://www.example.com/path123?query=123#fragment123
func shortenHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(consoleSpace)
	log.Println("Handling Shorten...")

	longUrl := req.URL.Query().Get("url")

	if longUrl == "" {
		http.Error(w, "Missing URL Parameter", http.StatusBadRequest)
		return
	}

	code, err := generateUniqueShortCode(6, longUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "http://localhost:8080/%s", code)
}

// GET /<shortcode>
func redirectHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(consoleSpace)
	log.Println("Handling Redirect...")

	code := req.URL.Path[1:] // strip leading "/"

	mu.RLock()
	longUrl, ok := store[code]
	mu.RUnlock()

	if !ok {
		http.NotFound(w, req)
		return
	}

	fmt.Println("Reading url:", req.URL)
	fmt.Println("Redirecting to", longUrl)

	http.Redirect(w, req, longUrl, http.StatusFound)
}

func main() {
	// could still use http.HandleFunc("/shorten", shortenHandler) this this would enforce only "POST" verb
	http.HandleFunc("POST /shorten", shortenHandler)
	http.HandleFunc("GET /", redirectHandler)

	log.Println("\nListening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
