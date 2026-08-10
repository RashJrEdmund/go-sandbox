## 1. CLI To-Do List

  A command-line app to add, list, complete, and delete tasks, persisting them to a local JSON file.

- I'll learn: structs, slices, reading/writing files (os, encoding/json), command-line args (os.Args or the flag package).
Stretch goals: add due dates, filter by status, use the cobra library for subcommands.

## ~2. Unit Converter / Calculator CLI~

  Convert between units (temperature, length, currency) or evaluate simple math from arguments.

- I'll learn: strconv for parsing strings to numbers, switch statements, functions, error handling with error return values.
Stretch goals: live currency rates by calling an HTTP API.

## ~3. URL Shortener (HTTP Server)~

  A web service where you POST a long URL and get back a short code that redirects.

- I'll learn: the net/http package, routing/handlers, HTTP methods, storing data in a map (then a real store later).
Stretch goals: persist to SQLite, add a tiny HTML frontend, generate collision-free codes.

## 4. Weather CLI (API Client)

  Fetch and display the current weather for a city using a free API (e.g. OpenWeather).

- I'll learn: making HTTP requests, parsing JSON responses into structs, environment variables for API keys, formatting output.
Stretch goals: cache results, support multiple cities, add a --forecast flag.

## ~5. Concurrent Link Checker (rlc)~

  Given a list of URLs, check which are alive and how fast they respond—concurrently.

- I'll learn: goroutines, channels, sync.WaitGroup, time measurement. This is the concept that makes Go special.
Stretch goals: limit concurrency with a worker pool, output a summary report, set request timeouts.
