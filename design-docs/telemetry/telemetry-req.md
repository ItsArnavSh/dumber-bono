# Telemetry Requirement

## Problem Statement
F1 25 streams all their telemetry data via UDP on specified port from the in game menu. The task is to build a server to receive that data, parse it and save into clickhouse db so analysis can be performed later.

## Scope
- Read the live UDP data from F1 25 
- Convert them to entity Structs
- Ingest all incoming data in Clickhouse, also define its schema

## Out of Scope
- No analysis at this stage
- No GRPC based lookups
