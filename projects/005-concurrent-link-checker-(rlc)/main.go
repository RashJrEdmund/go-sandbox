package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const requestTimeout = 5 * time.Second

type status struct {
	ALIVE   string
	DEAD    string
	INVALID string
}

var allowedStatuses = status{
	// ALIVE:   "ALIVE",
	DEAD:    "DEAD",
	INVALID: "INVALID",
}

type result struct {
	url     string
	status  string
	elapsed time.Duration
}

func checkLink(inputUrl string, wg *sync.WaitGroup, results chan<- result) {
	defer wg.Done()

	client := http.Client{Timeout: requestTimeout}
	start := time.Now()

	if u, err := url.Parse(inputUrl); err != nil || u.Scheme == "" {
		// for some reason, parsing a random string does not cause an error. so checking for a required field like the protocol at u.Scheme
		results <- result{inputUrl, allowedStatuses.INVALID, time.Since(start)}
		return
	}

	resp, err := client.Get(inputUrl)
	elapsed := time.Since(start)

	if err != nil {
		results <- result{inputUrl, allowedStatuses.DEAD, elapsed}
		return
	}
	defer resp.Body.Close()

	results <- result{inputUrl, resp.Status, elapsed}
}

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://go.dev",
			"https://google.com",
			"https://example.com",
			"https://thisurldoesnotexist.example",
			"not-a-url",
		}
		fmt.Println("No URLs provided — using demo list.")
		fmt.Println("Usage: rlc <url1> <url2> ...")
		fmt.Println()
	}

	var wg sync.WaitGroup
	results := make(chan result)

	start := time.Now()

	for _, url := range urls {
		wg.Add(1)
		go checkLink(url, &wg, results)
	}

	// Close the channel once all checks finish, so the range loop below can exit.
	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Printf("%-50s %-20s %s\n", "URL", "STATUS", "TIME")
	fmt.Println("--------------------------------------------------------------------------------")

	alive, dead, invalid := 0, 0, 0
	for r := range results {
		fmt.Printf("%-50s %-20s %v\n", r.url, r.status, r.elapsed.Round(time.Millisecond))

		switch r.status {
		case allowedStatuses.DEAD:
			dead++
		case allowedStatuses.INVALID:
			invalid++
		default:
			alive++
		}
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Checked %d links (%d alive, %d dead, %d invalid) in %v\n",
		len(urls), alive, dead, invalid, time.Since(start).Round(time.Millisecond))
}
