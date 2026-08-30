# Architecture

`dumber-bono` is a layered Go application. Data flows from the F1 game through
the API layer (ingress), into the service layer (business logic), and down to
the persistence/audio/LLM internals.

```
        ┌──────────────────────────────────────────────────────────────┐
        │                        main.go                                │
        │              getLogger() → api.NewServer() → StartServer()     │
        └──────────────────────────────┬───────────────────────────────┘
                                       │
        ┌──────────────────────────────▼───────────────────────────────┐
        │                       app/api (ingress)                      │
        │                                                              │
        │  udp/ ───── F1 UDP telemetry (port 4345)                     │
        │  ├─ parsers/    binary.Read → parser structs                 │
        │  ├─ mappers/    parser → domain tel-entity                   │
        │  └─ packethandler → routes by PacketID → ingestion service    │
        │                                                              │
        │  hotkeys/ ───── global hotkeys (Hyprland / cross-platform)    │
        │  mcp/ (WIP) ─── MCP tool server                              │
        └──────────────────────────────┬───────────────────────────────┘
                                       │  types.Ingestion / types.Monitor
        ┌──────────────────────────────▼───────────────────────────────┐
        │                    app/service (business logic)              │
        │                                                              │
        │  ingestionservice/  caches packets, accumulates motion+tele+lap│
        │  monitorservice/    events → radio msgs, driver pressure,     │
        │                     random callouts                          │
        │  radio/             priority queue, hotkeys, LLM→TTS→speaker  │
        │                                                              │
        │  internal/          badger (cache), duckdb (OLAP),           │
        │                     sqlite (relational), audio (STT/TTS),     │
        │                     llm (Groq)                               │
        └──────────────────────────────────────────────────────────────┘
```

## 1. Ingress: `app/api`

### UDP server (`app/api/udp`)

`udp.ListenUDP(ctx, logger, 4345, ingestion)` binds a UDP socket and spawns a
`listenUDP` goroutine. Each datagram is handed to `handle_packet`, which:

1. Parses the fixed **28-byte header** (`ParseHeader`) to learn the packet type.
2. Parses the payload into the matching `parsers.PacketXxxData` struct.
3. Maps it into a domain `telentity.Xxx` struct via the mappers.
4. Calls the corresponding `Ingest*` method on the ingestion service.

**Parsers** (`app/api/udp/parsers`) are a 1:1 mirror of the C structs in the
official F1 2025 UDP spec. They only do `binary.Read` — no logic.

**Mappers** (`app/api/udp/mappers`) translate raw parser structs into
`app/types/entity/tel-entity` domain structs, converting numeric enum codes
into human-readable strings via the `consts` package (track names, driver
names, weather, etc.) and booleans. They also validate `nil` inputs.

### Hotkeys (`app/api/hotkeys`)

`NewHotKeyListner()` returns a `HotKeyHandler`:

- `hyperland.HyperlandHotkeys` — reads `/dev/input/by-id/event-kbd*` directly
  for push-to-talk on Linux/Hyprland.
- `general.CrossPlatformHotkeys` — cross-platform hotkey library (used on
  non-Hyprland setups).

Hotkey events flow through a shared `chan entity.HotKeyEvent`:
`RADIO_PRESS`, `RADIO_RELEASE`, `COPY_AFFIRMATION`, `MUTE_TOGGLE`.

## 2. Business logic: `app/service`

### Ingestion service (`ingestionservice`)

Implements `types.Ingestion`. Receives every mapped packet and:

- **Header / session**: caches `sessionid`, `playerindex`, `lastupdated` in
  Badger. Detects session changes.
- **Participants**: caches the player's car/driver/team/race-index so other
  services know "my car".
- **Lap/session/setup/status/etc.**: throttled writes into Badger.
- **Motion / telemetry / lap**: pushed into a `TeleAccumulator`. When all three
  frames are fresh (within ~50ms via `utility.CheckFramesFresh`), it builds a
  `entity.TelemetryFrame` per car and batches them into DuckDB.

### Monitor service (`monitorservice`)

Implements `types.Monitor`. Two background loops:

1. **Pressure monitor** (`monitorPressure`) — every second, queries the latest
   frame from DuckDB and computes a 0–5 driver-pressure score from G-force,
   steering and braking (`AlterConfidence`).
2. **Random stats** (`RandomStatsMonitor`) — every 30s, if a live session is
   detected, runs a random subset of informational callouts (position, gaps,
   lap time, warnings) via the `Shuffler`.

`EventMonitor` turns race events into radio messages:

| Event | Behavior |
|-------|----------|
| FastestLap / RaceWinner / Retirement / Penalty | name lookup via driver consts → `DirectMessage` |
| Overtake | only if near the player or in the top 3 |
| Collision | only if the player isn't involved |
| SafetyCar / DRSDisabled / TeamMateInPits | direct callouts |
| SpeedTrap / StartLights | called out via `handleSpeedTrap` / `handleStartLights` |

### Radio service (`radio`)

The in-cab radio. Owns:

- **Priority queue** (`msg.go`) — `prio_sorted_vc` maps priority → expiry queue.
  `GetMessageByMinPriority` pops the highest-priority message at or above the
  current driver pressure. If muted, only max-priority messages surface.
- **Hotkey handling** (`controls.go`) — push-to-talk records audio via STT,
  sends the transcript to the streaming LLM, and enqueues the `IOPipe` reply
  at max priority so the driver hears the answer.
- **Streaming TTS** (`stream.go`) — `speakStream` buffers text until a sentence
  boundary (`.!?`), converts it to PCM via the Wyoming/piper server, and plays
  it through portaudio.

## 3. Internals: `app/service/internal`

| Package | Role |
|---------|------|
| `badger` | Key/value cache: `Set/Get/BulkSet/List` over namespaced keys. |
| `duckdb` | OLAP store for `telemetry_frames`; batch appender + `GetPressureFactors`. Embed-migrations applied on startup. |
| `sqlite` | Relational store (WAL mode, foreign keys on). |
| `audio/speaker` | Portaudio playback of PCM. |
| `audio/stt` | VAD (WebRTC) speech detection, RMS thresholding, WAV encoding, Groq Whisper transcription. |
| `audio/tts` | Wyoming protocol client (`synthesize` → `audio-start` → `audio-chunk` → `audio-stop`). |
| `llm` | Groq chat completions: `QueryLLM` (non-streaming) and `StreamLLM` (SSE streaming). |

## 4. Domain types: `app/types`

- `entity/` — shared domain: cache namespaces/keys, hotkey events, radio
  messages, `LapTimeStamp` formatting, coordinates, and the F1 const maps
  (drivers, teams, tracks, weather…).
- `entity/tel-entity/` — the "mapped" telemetry structs used across services.
- `interfaces/` — generic serialization constraint.

## Data flow example: a lap packet

```
UDP datagram (PacketLapData, 1285 bytes)
  → parsers.ParsePacket[PacketLapData]         (binary.Read)
  → mappers.MapToLapDataPacket                  (enums → strings, nil checks)
  → ingestion.IngestLapPacket
      → TeleAccumulator.UpsertLapData
      → throttled → repo.Cache.BulkSet(LAPDATA, kv)
```

When motion + telemetry + lap are all present and fresh, `SignalPush` writes a
`TelemetryFrame` per car to DuckDB, which later feeds the pressure monitor and
any analytics queries.
