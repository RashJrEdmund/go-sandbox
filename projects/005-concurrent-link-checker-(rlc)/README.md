# rlc

A concurrent link checker. Checks a list of URLs and reports which are alive, their HTTP status, and how long each took.

Statuses:

- **HTTP status** (e.g. `200 OK`) — the URL responded
- **DEAD** — the request failed (timeout, DNS, network error)
- **INVALID** — the input is not a valid URL (e.g. missing a scheme)

## Install

Requires Go installed ([download](https://go.dev/dl/)).

```bash
go install github.com/orashus/rlc@latest
```

This puts the `rlc` binary in your Go bin directory (usually `~/go/bin`).

<!-- Make sure that directory is on your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
``` -->

## Usage

```bash
rlc <url1> <url2> ...
```

Run with no arguments to check a built-in demo list:

```bash
rlc
```

## Example

```bash
rlc https://go.dev https://google.com https://example.com https://thisurldoesnotexist.example not-a-url
```

Output:

```bash
URL                                           STATUS               TIME
--------------------------------------------------------------------------------
not-a-url                                     INVALID              0s
https://thisurldoesnotexist.example           DEAD                 190ms
https://example.com                           200 OK               300ms
https://go.dev                                200 OK               250ms
https://google.com                            200 OK               180ms
--------------------------------------------------------------------------------
Checked 5 links (3 alive, 1 dead, 1 invalid) in 300ms
```

Results print as each check finishes (fast links first). Total time is roughly the slowest link — not the sum — because checks run concurrently.

## Develop locally

From this directory:

```bash
go run .
# or install your local copy
go install .
```
