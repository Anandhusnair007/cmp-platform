# Backend Services

This directory contains the Go backend services for the Certificate Management Platform.

## Services

- **inventory-service** (Port 8081): Certificate discovery and inventory
- **issuer-service** (Port 8082): Certificate issuance, renewal, and revocation
- **adapter-service** (Port 8083): Pluggable CA adapter layer

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis (optional)
- HashiCorp Vault

### Running Locally

```bash
# Start infrastructure
docker-compose -f ../deploy/docker-compose.yml up -d postgres redis vault

# Run migrations
make migrate-up

# Run a service
make run-inventory
make run-issuer
make run-adapter
```

### Environment Variables

- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database user (default: cmp_user)
- `DB_PASSWORD` - Database password (default: cmp_pass)
- `DB_NAME` - Database name (default: cmp_db)
- `REDIS_HOST` - Redis host (default: localhost)
- `VAULT_ADDR` - Vault address (default: http://localhost:8200)
- `VAULT_TOKEN` - Vault token
- `SERVER_PORT` - Service port

### Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Coverage
make test-coverage
```

### Building

```bash
# Build all services
make build

# Build Docker images
docker build -f Dockerfile.inventory -t cmp-inventory .
docker build -f Dockerfile.issuer -t cmp-issuer .
docker build -f Dockerfile.adapter -t cmp-adapter .
```

## API Documentation

See `api/openapi.yaml` for the complete OpenAPI 3.0 specification.
