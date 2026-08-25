# Ledger Config

> Backend service responsible for managing accounting configuration for the Ledger platform.

## Overview

Ledger Config is a backend service responsible for providing and managing accounting-related configuration used by the Ledger ecosystem.

The project is part of a distributed backend architecture composed of independent services, with a dedicated API responsible for centralizing configuration concerns.

The service was developed as a Go backend application with a structure designed to separate application, domain and infrastructure responsibilities.

## Architecture

The project follows a modular backend structure, separating the application entry point, configuration, documentation, internal implementation and infrastructure resources.

```text
ledger-config/
├── cmd/
├── config/
├── docs/
├── internal/
├── scripts/
├── docker-compose.yaml
├── Makefile
├── go.mod
└── go.sum
```

### Main components

#### `cmd/`

Contains the application entry point and service initialization.

#### `config/`

Contains application configuration and configuration-related resources.

#### `internal/`

Contains the application's internal implementation, keeping implementation details encapsulated within the service.

#### `docs/`

Contains project documentation and architectural resources.

#### `scripts/`

Contains supporting scripts and infrastructure resources used during development and deployment.

## Responsibilities

The main responsibility of Ledger Config is to provide a centralized API for accounting configuration.

The service is designed to:

- Manage accounting-related configuration
- Expose configuration through an API
- Centralize configuration rules used by Ledger services
- Isolate configuration concerns from other business services
- Provide a dedicated backend component within the Ledger ecosystem

## Architecture Principles

The project follows principles commonly used in modern backend systems:

- **Separation of concerns** — application responsibilities are separated from infrastructure concerns.
- **Encapsulation** — internal implementation details remain inside the service.
- **Service isolation** — accounting configuration is managed independently from other Ledger components.
- **Reusability** — configuration can be consumed by other services through the API.
- **Cloud-native development** — the repository includes container and infrastructure resources for local and deployment environments.

## Local Development

### Requirements

- Go
- Docker
- Docker Compose

### Clone

```bash
git clone https://github.com/clodoaldomarques/ledger-config.git
cd ledger-config
```

### Install dependencies

```bash
go mod download
```

### Run the application

```bash
go run ./cmd/...
```

### Run tests

```bash
go test ./...
```

### Run with Docker Compose

```bash
docker compose up
```

## Project Structure

```text
cmd/
    Application entry point

config/
    Application configuration

docs/
    Documentation and architectural resources

internal/
    Application implementation

scripts/
    Development and infrastructure scripts
```

## Development

A `Makefile` is provided to simplify common development and operational tasks.

Check the available commands with:

```bash
make help
```

If the project does not provide a `help` target, inspect the available targets directly:

```bash
make
```

## Technology Stack

| Technology | Purpose |
|---|---|
| Go | Backend service |
| Docker | Containerization |
| Docker Compose | Local development environment |
| Make | Development automation |

## Role in the Ledger Ecosystem

Ledger Config is designed as an independent service within the Ledger platform.

```text
                 ┌──────────────────────┐
                 │    Ledger Services   │
                 └──────────┬───────────┘
                            │
                            │ Configuration
                            ▼
                 ┌──────────────────────┐
                 │    Ledger Config     │
                 │                      │
                 │      Go API          │
                 └──────────────────────┘
```

By isolating accounting configuration into a dedicated service, other components of the platform can consume configuration without coupling their internal implementation to the configuration management logic.

## Engineering Concepts

This project demonstrates practical backend engineering concepts including:

- Go backend development
- REST API development
- Service-oriented architecture
- Separation of concerns
- Configuration management
- Containerization
- Local development automation
- Modular project organization

## Project Status

This project is part of my Go backend engineering portfolio and is intended to demonstrate the design and implementation of a dedicated configuration service within a distributed backend ecosystem.

## Author

**Clodoaldo Marques**

Backend Software Engineer focused on Go, Microservices, Distributed Systems and Cloud-Native architectures.

- GitHub: https://github.com/clodoaldomarques
- LinkedIn: https://www.linkedin.com/in/clodoaldomarques/