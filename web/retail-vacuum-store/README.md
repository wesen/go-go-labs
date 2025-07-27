# Vacuum Store - E-commerce Platform

A comprehensive e-commerce platform for vacuum cleaner sales built with Go, following go-go-golems patterns.

## Features

- 🌪️ Product catalog management
- 🛒 Shopping cart system
- 📦 Order management
- 👤 Admin panel
- 📱 Responsive design with Bootstrap
- 🔐 Secure configuration management
- 📊 Structured logging with zerolog

## Tech Stack

- **Backend**: Go 1.23+ with Cobra CLI
- **Frontend**: Bootstrap 5, vanilla JavaScript
- **Database**: SQLite (configurable)
- **Logging**: Zerolog
- **Configuration**: Viper

## Quick Start

### Prerequisites

- Go 1.23 or later
- Make (optional, for convenience commands)

### Setup

1. **Install dependencies**:
   ```bash
   make setup
   # or manually:
   go mod tidy
   ```

2. **Start development server**:
   ```bash
   make dev
   # or manually:
   ./scripts/dev.sh
   ```

3. **Open browser**: http://localhost:8080

### Available Commands

```bash
make dev      # Start development server with hot reload
make build    # Build production binary
make test     # Run tests and code checks
make clean    # Clean build artifacts
make run      # Run production binary
make fmt      # Format Go code
make lint     # Run linters (requires golangci-lint)
```

## Project Structure

```
web/retail-vacuum-store/
├── cmd/server/          # Application entry point
├── pkg/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers and routing
│   ├── models/          # Data models
│   └── services/        # Business logic
├── static/
│   ├── css/            # Stylesheets
│   ├── js/             # JavaScript files
│   └── images/         # Static images
├── templates/          # HTML templates (future)
├── db/                 # Database files
├── scripts/            # Build and development scripts
├── docs/               # Documentation
├── config.yaml         # Default configuration
├── Makefile           # Development commands
└── README.md          # This file
```

## Configuration

The application uses YAML configuration with environment variable override support:

```yaml
server:
  port: 8080
  host: "localhost"
  log_level: "info"

database:
  driver: "sqlite3"
  dsn: "./db/vacuum_store.db"

static:
  dir: "./static"
  prefix: "/static/"
```

### Environment Variables

You can override configuration using environment variables:

```bash
export PORT=3000
export LOG_LEVEL=debug
export DATABASE_DSN="./custom.db"
```

## Development

### Running in Development Mode

```bash
# Start with debug logging
./scripts/dev.sh

# Or with custom port
go run ./cmd/server --port=3000 --log-level=debug
```

### Building for Production

```bash
make build
cd build
./vacuum-store --config=config.yaml
```

### Running Tests

```bash
make test
# Runs:
# - Unit tests
# - Race condition tests
# - Code formatting checks
# - Go vet analysis
```

## API Endpoints

### Public Routes

- `GET /` - Home page
- `GET /products` - Product catalog
- `GET /cart` - Shopping cart
- `GET /health` - Health check
- `GET /static/*` - Static assets

### Admin Routes

- `GET /admin` - Admin panel

### Future API Routes

- `GET /api/products` - Product API
- `POST /api/orders` - Order submission
- `GET /api/cart` - Cart management

## Logging

The application uses structured logging with zerolog:

```bash
# Development (console output)
go run ./cmd/server --log-level=debug

# Production (JSON output)
LOG_LEVEL=info ./vacuum-store
```

## Contributing

1. Follow go-go-golems patterns
2. Use zerolog for logging
3. Follow Go naming conventions
4. Add tests for new features
5. Run `make test` before committing

## License

This project is part of the go-go-golems ecosystem.
