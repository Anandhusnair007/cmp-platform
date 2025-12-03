# Developer Onboarding Guide

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Node.js 18+ and npm
- Make
- PostgreSQL client tools (optional)
- Terraform (for infrastructure work)

## Initial Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd cmp-platform
```

### 2. Start Development Environment

```bash
# Start all services
docker-compose -f deploy/docker-compose.yml up --build

# Wait for services to be healthy (~30 seconds)
# Check status
docker-compose -f deploy/docker-compose.yml ps
```

### 3. Run Database Migrations

```bash
# Install golang-migrate if needed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply migrations
make migrate-up
```

### 4. Initialize Vault PKI

```bash
# Start Vault (if not already running)
docker-compose -f deploy/docker-compose.yml up -d vault

# Wait for Vault to be ready
sleep 5

# Initialize PKI (manual steps)
docker exec -it cmp-vault sh

# Inside Vault container
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=dev-only-token

# Enable PKI
vault secrets enable -path=cmp-pki pki

# Configure PKI
vault secrets tune -max-lease-ttl=87600h cmp-pki

# Generate root CA
vault write -field=certificate cmp-pki/root/generate/internal \
    common_name="CMP Root CA" \
    ttl=87600h > /tmp/ca.crt

# Configure URLs
vault write cmp-pki/config/urls \
    issuing_certificates="http://127.0.0.1:8200/v1/cmp-pki/ca" \
    crl_distribution_points="http://127.0.0.1:8200/v1/cmp-pki/crl"

# Create role
vault write cmp-pki/roles/cmp-role \
    allowed_domains="staging.example.com,example.com" \
    allow_subdomains=true \
    max_ttl="8760h"
```

Or use Terraform:

```bash
cd infra/terraform/vault-pki
terraform init
terraform plan
terraform apply
```

## Running Services Locally

### Backend Services

```bash
# Start infrastructure only
docker-compose -f deploy/docker-compose.yml up -d postgres redis vault

# Run services locally (in separate terminals)
make run-inventory      # Port 8081
make run-issuer         # Port 8082
make run-adapter        # Port 8083
```

### Frontend

```bash
cd frontend/webapp
npm install
npm run dev  # Port 3000
```

### Agent

```bash
cd agents/linux-agent
go run ./cmd/agent \
  -api-url=http://localhost:8082 \
  -agent-id=agent-1 \
  -hostname=$(hostname)
```

## Development Workflow

### Creating Database Migrations

```bash
make migrate-create NAME=add_new_column

# Edit the generated files in backend/migrations/
# Then apply:
make migrate-up
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# With coverage
make test-coverage
```

### Code Style

- Use `gofmt` for Go formatting
- Follow Go style guide
- Run `golangci-lint` before committing

## Project Structure

```
cmp-platform/
├── backend/              # Go backend services
│   ├── cmd/             # Service entry points
│   ├── internal/        # Internal packages
│   ├── migrations/      # Database migrations
│   └── api/             # OpenAPI spec
├── agents/              # Agent implementations
├── frontend/webapp/     # React frontend
├── k8s-operator/        # Kubernetes operator
├── tests/               # Test suites
├── deploy/              # Docker Compose configs
├── infra/               # Infrastructure as Code
└── docs/                # Documentation
```

## Key Concepts

### Certificate Lifecycle

1. **Request** - User/API requests certificate
2. **CSR Generation** - Create certificate signing request
3. **Signing** - CA adapter signs certificate
4. **Storage** - Store in database and Vault
5. **Installation** - Agent installs to target system
6. **Monitoring** - Track expiry and renew

### Adapter Pattern

Adapters abstract CA-specific logic:
- Vault PKI adapter uses Vault API
- ACME adapter uses ACME protocol
- Easy to add new CAs

### Agent Pattern

Agents run on target hosts:
- Register with CMP
- Poll for installation jobs
- Download and install certificates
- Execute reload commands
- Report status

## Common Development Tasks

### Add New API Endpoint

1. Update `backend/api/openapi.yaml`
2. Add handler in `backend/internal/handlers/`
3. Add route in service main.go
4. Add tests
5. Update frontend if needed

### Add New Adapter

1. Implement adapter interface in `backend/internal/adapters/`
2. Register adapter in database
3. Add configuration schema
4. Add tests
5. Document in architecture.md

### Debug Issues

1. Check service logs: `docker-compose logs <service>`
2. Query database: `psql -h localhost -U cmp_user -d cmp_db`
3. Check Vault: `vault read cmp-pki/config/urls`
4. Review audit logs: `SELECT * FROM audit_logs ORDER BY timestamp DESC LIMIT 10;`

## Useful Commands

```bash
# View logs
docker-compose -f deploy/docker-compose.yml logs -f issuer-service

# Access database
docker exec -it cmp-postgres psql -U cmp_user -d cmp_db

# Access Vault
docker exec -it cmp-vault vault status

# Run E2E test
./tests/e2e/run.sh

# Build all services
make build
```

## Next Steps

1. Read [Architecture](./architecture.md) for system design
2. Review [Runbooks](./runbooks.md) for operations
3. Explore the codebase
4. Run the E2E test to see full flow
5. Pick a task from the backlog

## Getting Help

- Check documentation in `docs/`
- Review code comments
- Ask team in Slack/Discord
- Create GitHub issue for bugs

## Resources

- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [Vault PKI Documentation](https://developer.hashicorp.com/vault/docs/secrets/pki)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
