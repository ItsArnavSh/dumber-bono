# dumber-bono

A real-time F1 2025 telemetry ingestion and in-cab radio assistant. It listens
to the game's UDP telemetry stream, parses every packet type, persists the
data, and broadcasts race information to the driver over an in-cab "radio"
using text-to-speech.

![CI](https://github.com/ItsArnavSh/dumber-bono/actions/workflows/ci.yml/badge.svg)

## Features

- **F1 2025 UDP telemetry** — parses all 16 packet types (motion, lap, session,
  event, participants, car setup, telemetry, car status, lobby, car damage,
  session history, tyre sets, and more).
- **Streaming ingestion pipeline** — packets are parsed, mapped to domain
  entities, throttled, and persisted into a local analytics DB (DuckDB), a
  SQLite store, and a Badger key/value cache.
- **Driver pressure model** — computes a 0–5 "pressure" score from lateral/long
  G-force, steering and braking; drives which radio messages get spoken.
- **Event monitor** — turns race events (fastest laps, overtakes, safety cars,
  penalties, collisions…) into spoken messages.
- **In-cab radio** — a priority queue delivers messages to a streaming TTS
  pipeline (Wyoming/piper), with push-to-talk speech-to-text (Groq Whisper).
- **Global hotkeys** — push-to-talk (R), copy affirmation (C), mute (M).
- **LLM streaming** — a Groq-backed chat model (with a raw `StreamLLM` SSE
  implementation and an Eino adapter) answers driver questions over the radio.

## Architecture

```
main.go
 └─ boot.go                 zap logger
 └─ app/api
     ├─ server.go           wires hotkeys, services, and the UDP listener
     ├─ udp/                F1 telemetry UDP server
     │   ├─ parsers/        raw binary → parser structs
     │   ├─ mappers/        parser structs → domain tel-entity structs
     │   └─ packethandler   routes by packet ID → ingestion service
     ├─ hotkeys/            Hyprland / cross-platform global hotkeys
     └─ mcp/                (WIP) MCP tool server
 └─ app/service
     ├─ ingestionservice/   receives packets, caches + accumulates frames
     ├─ monitorservice/     event → radio message; driver pressure; random callouts
     ├─ radio/              priority queue, hotkey handling, LLM + TTS + STT
     └─ internal/
         ├─ badger/         key/value cache
         ├─ duckdb/         OLAP telemetry frames + pressure queries
         ├─ sqlite/         relational store
         ├─ audio/          portaudio speaker, STT (VAD + Whisper), TTS (Wyoming)
         └─ llm/            Groq chat completions (streaming + non-streaming)
 └─ app/types/entity        domain entities + constants (F1 driver/track/team maps)
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full deep-dive.

## Requirements

- **Go 1.26+** and CGO toolchain (portaudio + VAD are C libraries).
- Linux packages: `libasound2-dev`, `portaudio19-dev` (Ubuntu/Debian).
- A running [wyoming-piper](https://github.com/rhasspy/wyoming-piper) TTS
  server on `127.0.0.1:10200` (see `docker-compose.yml`).
- A [Groq](https://console.groq.com) API key for STT/LLM:
  - `GROQ_API_KEY`
  - `GROQ_STT_MODEL` (default `whisper-large-v3`)

## Setup

```sh
# 1. Install Go dependencies
go mod download

# 2. Start the TTS server
docker compose up -d

# 3. Configure environment
cp .env.example .env   # or export GROQ_API_KEY=...

# 4. Build & run
just run               # or: go run .
```

## Usage

The app listens for F1 2025 UDP telemetry on **port 4345** (configure the game
to send UDP telemetry to this port).

| Hotkey | Action |
|--------|--------|
| `R` (hold) | Push-to-talk — speak, release to have your message answered via LLM + TTS |
| `C`        | Copy affirmation (WIP) |
| `M`        | Mute / unmute the radio (always lets max-priority messages through) |

## Commands (Justfile)

```sh
just run          # build & run the app
just build        # compile to bin/dumber-bono
just test         # run all tests
just test-race    # run tests with the race detector
just cover        # test + coverage report
just lint         # run golangci-lint
just lint-fix     # auto-fix lint issues
just fmt          # gofmt + goimports
just vet          # go vet
just ci           # fmt-check → vet → lint → test (matches CI)
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for more.

## CI/CD

- [`.github/workflows/ci.yml`](.github/workflows/ci.yml) — lint, vet, test,
  race detector on every push/PR.
- [`.github/workflows/build.yml`](.github/workflows/build.yml) — cross-platform
  linux binaries (amd64/arm64) uploaded as artifacts.

## Project layout

```
Justfile          task runner
.golangci.yml      golangci-lint v2 config
docker-compose.yml piper TTS server
priority.md        scratchpad / bug notes
```

## License

TBD.
