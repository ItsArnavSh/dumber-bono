# Development

This project uses [just](https://github.com/casey/just) as its task runner and
[golangci-lint](https://golangci-lint.run) (v2) for linting.

## Prerequisites

- Go 1.26+ with CGO (`CGO_ENABLED=1`).
- Linux C deps (portaudio + ALSA):
  ```sh
  sudo apt-get install -y libasound2-dev portaudio19-dev
  ```
- `just` and `golangci-lint`:
  ```sh
  cargo install just            # or via your package manager
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
  ```

## Common tasks

| Task | Command | Purpose |
|------|---------|---------|
| Run | `just run` | Build & run the app |
| Build | `just build` | Compile to `bin/dumber-bono` |
| Test | `just test` | `go test ./...` |
| Test (verbose) | `just test-verbose` | `go test -v ./...` |
| Test (race) | `just test-race` | `go test -race ./...` |
| Coverage | `just cover` | Test + coverage summary |
| Lint | `just lint` | `golangci-lint run ./...` |
| Auto-fix lint | `just lint-fix` | `golangci-lint run --fix ./...` |
| Format | `just fmt` | `gofmt -w` + `goimports -w` |
| Vet | `just vet` | `go vet ./...` |
| Tidy deps | `just tidy` | `go mod tidy` |
| CI replica | `just ci` | fmt-check → vet → lint → test |

## Adding unit tests

Tests live next to the code they cover (`*_test.go`). Existing coverage:

| Package | Notes |
|---------|-------|
| `app/utility` | Throttler, ExpiryQueue, RandomString, CheckFramesFresh |
| `app/api/udp/parsers` | header / event / generic packet parsing |
| `app/api/udp/mappers` | header / telemetry / lapdata / motion / event mapping |
| `app/service/radio` | priority queue + message listener |
| `app/service/monitorservice` | pressure model, shuffler, position checks |
| `app/service/ingestionservice` | telemetry accumulator |
| `app/service/internal/audio/stt` | RMS, WAV encode/decode, VAD state machine |
| `app/service/internal/badger` | cache set/get/bulk/list |
| `app/types/entity` | `LapTimeStamp` formatting |

Run the suite (including the race detector — important, the code is
concurrency-heavy):

```sh
just test-race
```

## Linting

Config: `.golangci.yml` (golangci-lint v2). Enabled linters: `errcheck`,
`govet`, `staticcheck`, `unused`, `ineffassign`, `misspell`, `unconvert`,
`prealloc`, `noctx`, `bodyclose`, `sqlclosecheck`; formatters `gofmt` +
`goimports`.

```sh
just lint
```

Two intentional exclusions in the config:

- `misspell` is skipped on `app/types/entity/consts/f125-consts.go` because
  "Sargeant" (F1 driver Logan Sargeant) is a real proper noun.
- `errcheck` / `bodyclose` are relaxed in `_test.go` files.

## CI

GitHub Actions mirrors `just ci`:

- `.github/workflows/ci.yml` — checkout → Go 1.26.1 → C deps → `go vet` →
  `golangci-lint` → `go test` → `go test -race`.
- `.github/workflows/build.yml` — linux/amd64 + linux/arm64 builds uploaded as
  artifacts.

## Useful tips

- The UDP server binds port `4345` and the TTS server is expected at
  `127.0.0.1:10200`. Both are hardcoded — adjust in `app/api/server.go`,
  `app/api/udp/udp.go`, and `app/service/internal/audio/tts/tts.go` if needed.
- `GROQ_API_KEY` is required for STT and LLM features.
- The badger cache, sqlite and duckdb files are written to the path passed to
  `api.NewServer` (default `/tmp` in `main.go`). Point it somewhere persistent
  for real sessions.
