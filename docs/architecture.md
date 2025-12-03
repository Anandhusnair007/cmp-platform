# CMP Architecture

## Overview

The Certificate Management Platform (CMP) is a distributed system designed to automate TLS/SSL certificate lifecycle management across hybrid cloud environments.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                        │
│                  Port 3000 - Dashboard & UI                     │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTPS/REST
┌───────────────────────────▼─────────────────────────────────────┐
│                    API Gateway / Load Balancer                  │
└───────────┬───────────┬───────────┬───────────┬────────────────┘
            │           │           │           │
    ┌───────▼───┐ ┌─────▼────┐ ┌───▼─────┐ ┌──▼──────┐
    │Inventory  │ │  Issuer  │ │ Adapter │ │ Agents  │
    │ Service   │ │ Service  │ │ Service │ │ Handler │
    │ :8081     │ │ :8082    │ │ :8083   │ │ :8082   │
    └───────┬───┘ └─────┬────┘ └───┬─────┘ └─────────┘
            │           │           │
    ┌───────┴───────────┴───────────┴──────────────┐
    │            Shared Database (PostgreSQL)       │
    │         - certificates, requests, agents      │
    │         - audit_logs, adapters, owners        │
    └───────────────────────────────────────────────┘
            │           │           │
    ┌───────▼─────┐ ┌───▼──────┐ ┌──▼──────────┐
    │    Redis    │ │  Vault   │ │   Agents    │
    │   (Cache)   │ │ (Secrets)│ │  (Linux)    │
    └─────────────┘ └──────────┘ └─────────────┘
                                          │
                              ┌───────────▼──────────┐
                              │  Target Systems      │
                              │  - nginx, apache     │
                              │  - Kubernetes        │
                              │  - Load Balancers    │
                              └──────────────────────┘
```

## Components

### Backend Services

#### Inventory Service
- **Purpose**: Certificate discovery and inventory management
- **Port**: 8081
- **Responsibilities**:
  - Scan and catalog existing certificates
  - Track certificate expiry
  - Provide inventory API endpoints

#### Issuer Service
- **Purpose**: Certificate issuance, renewal, and revocation
- **Port**: 8082
- **Responsibilities**:
  - Handle certificate requests
  - Coordinate with adapter service
  - Manage certificate lifecycle
  - Agent management and installation jobs

#### Adapter Service
- **Purpose**: Pluggable CA adapter layer
- **Port**: 8083
- **Responsibilities**:
  - Interface with Certificate Authorities (Vault PKI, ACME, etc.)
  - Sign certificate requests
  - Revoke certificates
  - Adapter configuration management

### Agents

#### Linux Agent
- **Purpose**: Deploy certificates to Linux hosts
- **Responsibilities**:
  - Register with CMP
  - Receive installation jobs
  - Download and install certificates
  - Execute reload commands
  - Report status

### Data Layer

#### PostgreSQL
- Primary database for all state
- Tables: certificates, issuance_requests, agents, audit_logs, adapters, owners

#### Redis
- Caching layer
- Session storage
- Distributed locks

#### HashiCorp Vault
- Secret storage (private keys, certificates)
- PKI engine for certificate signing
- Token management for agents

## Security Architecture

### Authentication & Authorization
- OIDC integration (Keycloak/Auth0)
- RBAC roles: admin, security, developer, agent
- Service-to-service: mTLS (planned)

### Secret Management
- Private keys stored encrypted in Vault
- Agent-side key generation preferred
- Never log secrets in plaintext

### Audit Logging
- All actions logged to audit_logs table
- Append-only storage (object store with WORM)
- Includes: entity type, action, performer, timestamp, details

## Adapter Architecture

Adapters implement a common interface:
- `CreateCertificate()` - Issue new certificate
- `RevokeCertificate()` - Revoke existing certificate
- `GetStatus()` - Get adapter/certificate status

### Supported Adapters
1. **Vault PKI** - HashiCorp Vault PKI engine
2. **ACME** - Let's Encrypt compatible (test server)

### Future Adapters
- Venafi
- AWS Certificate Manager
- Azure Key Vault
- HSM (PKCS#11)

## Deployment Architecture

### Development
- Docker Compose with all services
- Vault dev mode
- Single-node PostgreSQL

### Production
- Kubernetes deployment via Helm
- High availability (multiple replicas)
- PostgreSQL with replicas
- Redis cluster
- External Vault cluster
- Multi-AZ deployment

## Network Security

- All services communicate via internal network
- mTLS between services (planned)
- Agents authenticate via Vault tokens or mTLS
- Frontend communicates via HTTPS to API

## Monitoring & Observability

- Prometheus metrics endpoints on each service
- Grafana dashboards for:
  - Certificate expiry timeline
  - Issuance latency
  - Agent health
  - Error rates
- Structured JSON logs to Elasticsearch/OpenSearch

## Scalability Considerations

- Stateless services for horizontal scaling
- Database connection pooling
- Redis for distributed caching
- Agent registration and job queuing
- Async certificate processing

## Disaster Recovery

- Database backups (daily)
- Vault unseal keys stored securely
- Certificate backups in object store
- Agent recovery procedures
