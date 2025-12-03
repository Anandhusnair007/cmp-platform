# Certificate Management Platform (CMP)

A production-grade Certificate Management Platform that automates discovery, issuance, renewal, revocation, deployment, and monitoring of TLS/SSL certificates across cloud, on-prem, and hybrid environments.

## 🚀 Host This Application

**Quick Hosting:** Clone this repository on your server and run:
```bash
sudo ./deploy/production/install.sh
sudo ./deploy/production/scripts/deploy.sh
```

📖 **See [HOSTING_ON_GITHUB.md](HOSTING_ON_GITHUB.md) for complete hosting instructions**

🔗 **Repository**: https://github.com/Anandhusnair007/cmp-platform

## 🚀 Quick Start (10-Minute Demo)

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 18+ and npm
- Make

### Start the Full Stack

```bash
# Clone and navigate
cd cmp-platform

# Start all services
docker-compose -f deploy/docker-compose.yml up --build

# In a separate terminal, run database migrations
make migrate-up

# Wait for services to be healthy (~30 seconds)
# Access:
# - Frontend: http://localhost:3000
# - API: http://localhost:8080/api/v1
# - Vault UI: http://localhost:8200 (token: dev-only-token)
```

### Run End-to-End Test

```bash
# Run the full E2E scenario
./tests/e2e/run.sh

# This will:
# 1. Request a certificate via API
# 2. Issue it via Vault PKI adapter
# 3. Install it to nginx via agent
# 4. Verify HTTPS connection works
```

## 📋 Architecture Overview

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Frontend  │────▶│  REST API    │────▶│   Services  │
│  (React)    │     │  (OpenAPI)   │     │    (Go)     │
└─────────────┘     └──────────────┘     └─────────────┘
                                               │
                    ┌──────────────────────────┼──────────────────────────┐
                    │                          │                          │
              ┌─────▼──────┐          ┌────────▼──────┐          ┌───────▼──────┐
              │ PostgreSQL │          │     Redis     │          │    Vault     │
              │  (State)   │          │   (Cache)     │          │  (Secrets)   │
              └────────────┘          └───────────────┘          └──────────────┘
                                               │
                                    ┌──────────▼──────────┐
                                    │   Linux Agents      │
                                    │  (Certificate       │
                                    │   Installation)     │
                                    └─────────────────────┘
```

## 🏗️ Components

### Backend Services

- **inventory-service**: Discovery and inventory of certificates across environments
- **issuer-service**: Certificate issuance, renewal, and revocation
- **adapter-service**: Pluggable CA adapter layer (Vault PKI, ACME)

### Agents

- **linux-agent**: Deploys certificates to Linux hosts and reloads services

### Frontend

- **webapp**: React + TypeScript dashboard for certificate management

### Infrastructure

- **Kubernetes Operator**: Manages Certificate CRDs and Kubernetes Secrets
- **Helm Charts**: K8s deployment manifests
- **Terraform**: Cloud infrastructure provisioning

## 📁 Repository Structure

```
/
├── README.md                  # This file
├── .github/workflows/         # CI/CD pipelines
├── infra/terraform/           # Infrastructure as Code
├── charts/cmp/                # Helm charts
├── deploy/                    # Docker Compose configs
├── backend/                   # Go backend services
│   ├── cmd/                   # Service entry points
│   ├── internal/              # Internal packages
│   └── api/                   # OpenAPI spec
├── agents/                    # Agent implementations
├── k8s-operator/              # Kubernetes operator
├── frontend/webapp/           # React frontend
├── tests/                     # Test suites
└── docs/                      # Documentation
```

## 🔧 Development

### Local Development Setup

```bash
# Install dependencies
cd backend && go mod download
cd ../frontend/webapp && npm install

# Start infrastructure services only
docker-compose -f deploy/docker-compose.yml up postgres redis vault -d

# Run services locally
make run-inventory      # Port 8081
make run-issuer         # Port 8082
make run-adapter        # Port 8083

# Run frontend dev server
cd frontend/webapp && npm run dev  # Port 3000
```

### Database Migrations

```bash
# Create a new migration
make migrate-create NAME=add_audit_logs

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down
```

### Testing

```bash
# Unit tests
make test

# Integration tests
make test-integration

# E2E tests
./tests/e2e/run.sh

# Test coverage
make test-coverage
```

### Building

```bash
# Build all services
make build

# Build Docker images
make docker-build

# Build agent binary
cd agents/linux-agent && make build
```

## 🔐 Security

- **Secrets Management**: HashiCorp Vault (dev/staging)
- **Service Communication**: mTLS between services
- **RBAC**: Role-based access control with OIDC (Keycloak)
- **Audit Logging**: All actions logged to append-only audit_logs table
- **Key Generation**: Agent-side preferred, or Vault-encrypted storage

## 📊 Monitoring

- **Metrics**: Prometheus endpoints at `/metrics` on each service
- **Dashboards**: Grafana dashboard templates in `docs/grafana/`
- **Logging**: Structured JSON logs (Elasticsearch/OpenSearch)

## 🚢 Deployment

### Docker Compose (Development)

```bash
docker-compose -f deploy/docker-compose.yml up
```

### Kubernetes (Production)

```bash
# Install Helm chart
helm install cmp ./charts/cmp --namespace cmp --create-namespace

# Configure values
helm upgrade cmp ./charts/cmp --set image.tag=v0.1.0
```

## 📚 Documentation

- [Architecture](./docs/architecture.md) - System architecture and design decisions
- [Runbooks](./docs/runbooks.md) - Operational procedures
- [Onboarding](./docs/onboarding.md) - Developer onboarding guide
- [API Documentation](./backend/api/openapi.yaml) - OpenAPI 3.0 specification

## 🔌 Adapters

### Implemented

- **Vault PKI**: HashiCorp Vault PKI engine integration
- **ACME**: Let's Encrypt compatible ACME protocol (test server)

### Planned

- Venafi integration
- HSM PKCS#11 integration
- Multi-CA policy engine

## 🧪 Acceptance Criteria

✅ Certificate can be requested via API and issued via Vault adapter  
✅ Agent automatically installs certificate to nginx and reloads service  
✅ Dashboard displays issued certificate and expiry  
✅ RBAC prevents unauthorized actions  
✅ Audit logs capture all issuance and installation events  
✅ CI pipeline runs and passes unit+integration tests  

## 🤝 Contributing

1. Create a feature branch from `main`
2. Make changes with semantic commits (`feat/`, `fix/`, `chore/`, etc.)
3. Ensure tests pass: `make test`
4. Open a Pull Request with description

## 📄 License

MIT License - See LICENSE file for details

## 🗺️ Roadmap

- [ ] Multi-cloud CA adapter support
- [ ] Automated certificate renewal with job scheduler
- [ ] ServiceNow integration for approval workflows
- [ ] Certificate compliance scanning
- [ ] Multi-tenant isolation
- [ ] HSM integration for key management
# cmp-platform
