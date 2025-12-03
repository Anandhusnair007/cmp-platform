# CMP Platform - Enterprise Production Deployment Guide

**Version:** 2.0.0  
**Target:** Enterprise-grade deployments (Wipro, TCS, Microsoft Security, Cisco standards)  
**Deployment Method:** Native Linux services (No Docker)

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Security Hardening](#security-hardening)
6. [Deployment](#deployment)
7. [Monitoring & Logging](#monitoring--logging)
8. [Maintenance](#maintenance)
9. [Troubleshooting](#troubleshooting)

---

## Overview

This guide provides step-by-step instructions for deploying the Certificate Management Platform (CMP) as a production-grade system using native Linux services. The deployment follows enterprise security standards suitable for organizations like Wipro, TCS, Microsoft Security, and Cisco.

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Load Balancer                        │
│                    (Nginx/HAProxy)                       │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
   ┌────▼────┐  ┌───▼───┐  ┌────▼────┐
   │  Auth   │  │ Issuer│  │Inventory│
   │ Service │  │Service│  │ Service │
   └─────────┘  └───────┘  └─────────┘
        │            │            │
        └────────────┼────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
   ┌────▼────┐  ┌───▼───┐  ┌────▼────┐
   │PostgreSQL│  │ Redis │  │  Vault  │
   └──────────┘  └───────┘  └─────────┘
```

### Components

- **Backend Services**: Go-based microservices (inventory, issuer, adapter, auth)
- **Frontend**: React/TypeScript web application
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Secrets**: HashiCorp Vault
- **Reverse Proxy**: Nginx with SSL/TLS
- **Monitoring**: Prometheus + Grafana
- **Logging**: Systemd journal + centralized logging

---

## Prerequisites

### System Requirements

- **OS**: Ubuntu 20.04+, Debian 11+, RHEL 8+, CentOS 8+
- **CPU**: 4+ cores recommended
- **RAM**: 8GB+ recommended (16GB+ for production)
- **Disk**: 100GB+ SSD recommended
- **Network**: Static IP address, DNS records configured

### Software Requirements

- Go 1.21+
- Node.js 18+ and npm
- PostgreSQL 15+
- Redis 7+
- Nginx 1.20+
- Systemd
- OpenSSL

### Network Requirements

- Port 80 (HTTP - redirect to HTTPS)
- Port 443 (HTTPS - production traffic)
- Port 5432 (PostgreSQL - restrict to localhost)
- Port 6379 (Redis - restrict to localhost)
- Port 8200 (Vault - restrict to localhost)

---

## Installation

### Step 1: Prepare System

```bash
# Update system packages
sudo apt-get update && sudo apt-get upgrade -y

# Install basic tools
sudo apt-get install -y curl wget git build-essential

# Set hostname
sudo hostnamectl set-hostname cmp-production

# Configure timezone
sudo timedatectl set-timezone UTC
```

### Step 2: Clone Repository

```bash
# Clone the repository
cd /opt
sudo git clone <repository-url> cmp-platform
cd cmp-platform
```

### Step 3: Run Installation Script

```bash
# Make script executable
sudo chmod +x deploy/production/install.sh

# Run installation (as root)
sudo ./deploy/production/install.sh
```

The installation script will:
- Install all system dependencies
- Create CMP user and directories
- Install PostgreSQL, Redis, Vault
- Build backend services
- Build frontend application
- Install systemd services
- Configure Nginx
- Setup log rotation

**Expected Duration:** 15-30 minutes

---

## Configuration

### Step 1: Database Configuration

```bash
# Edit PostgreSQL configuration
sudo nano /etc/postgresql/15/main/postgresql.conf

# Set secure settings
listen_addresses = 'localhost'
ssl = on
ssl_cert_file = '/etc/ssl/certs/ssl-cert-snakeoil.pem'
ssl_key_file = '/etc/ssl/private/ssl-cert-snakeoil.key'

# Edit pg_hba.conf
sudo nano /etc/postgresql/15/main/pg_hba.conf

# Add secure authentication
host    cmp_db    cmp_user    127.0.0.1/32    scram-sha-256

# Restart PostgreSQL
sudo systemctl restart postgresql
```

### Step 2: Redis Configuration

```bash
# Edit Redis configuration
sudo nano /etc/redis/redis.conf

# Set secure settings
bind 127.0.0.1
protected-mode yes
requirepass YOUR_SECURE_REDIS_PASSWORD

# Restart Redis
sudo systemctl restart redis
```

### Step 3: Vault Configuration

```bash
# Create Vault configuration
sudo mkdir -p /etc/vault.d

# Create vault.hcl
sudo nano /etc/vault.d/vault.hcl
```

Add the following configuration:

```hcl
ui = true

storage "file" {
  path = "/var/lib/vault"
}

listener "tcp" {
  address     = "127.0.0.1:8200"
  tls_disable = 1
}

api_addr = "http://127.0.0.1:8200"
```

Initialize Vault:

```bash
# Start Vault
sudo systemctl start vault

# Initialize Vault
vault operator init

# Save unseal keys and root token securely
# Unseal Vault
vault operator unseal <unseal-key-1>
vault operator unseal <unseal-key-2>
vault operator unseal <unseal-key-3>
```

### Step 4: Service Configuration

Copy environment templates and configure:

```bash
# Copy configuration templates
sudo cp deploy/production/config/*.env.example /etc/cmp/

# Edit each service configuration
sudo nano /etc/cmp/cmp-inventory.env
sudo nano /etc/cmp/cmp-issuer.env
sudo nano /etc/cmp/cmp-adapter.env
sudo nano /etc/cmp/cmp-auth.env
```

**Important Settings:**

- `DB_PASSWORD`: Use strong password (32+ characters)
- `REDIS_PASSWORD`: Use strong password
- `JWT_SECRET`: Generate with `openssl rand -hex 32`
- `VAULT_TOKEN`: Use Vault root token or app token

### Step 5: SSL/TLS Configuration

```bash
# Make SSL setup script executable
sudo chmod +x deploy/production/scripts/setup-ssl.sh

# Run SSL setup (choose Let's Encrypt for production)
sudo ./deploy/production/scripts/setup-ssl.sh
```

For Let's Encrypt:
- Enter API domain: `api.cmp.example.com`
- Enter App domain: `app.cmp.example.com`
- Enter email for notifications

### Step 6: Nginx Configuration

```bash
# Update Nginx configuration with your domains
sudo nano /etc/nginx/cmp/cmp-nginx.conf

# Replace example.com with your domain
# Test configuration
sudo nginx -t -c /etc/nginx/cmp/cmp-nginx.conf
```

---

## Security Hardening

### Step 1: Run Security Hardening Script

```bash
# Make script executable
sudo chmod +x deploy/production/scripts/harden-system.sh

# Run hardening
sudo ./deploy/production/scripts/harden-system.sh
```

This script will:
- Harden kernel parameters
- Configure system limits
- Harden SSH configuration
- Setup fail2ban
- Configure audit logging
- Generate secure passwords
- Setup SSL certificates

### Step 2: Firewall Configuration

```bash
# Enable firewall
sudo ufw enable

# Allow required ports
sudo ufw allow 22/tcp comment 'SSH'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'

# Deny direct access to services
sudo ufw deny 5432/tcp comment 'PostgreSQL'
sudo ufw deny 6379/tcp comment 'Redis'
sudo ufw deny 8200/tcp comment 'Vault'
```

### Step 3: Additional Security

```bash
# Disable root login via SSH
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Setup automatic security updates
sudo apt-get install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades

# Install and configure AIDE (file integrity monitoring)
sudo apt-get install -y aide
sudo aideinit
sudo mv /var/lib/aide/aide.db.new /var/lib/aide/aide.db
```

---

## Deployment

### Step 1: Run Database Migrations

```bash
# Navigate to backend directory
cd /opt/cmp-platform/backend

# Run migrations
make migrate-up

# Verify schema
psql -h localhost -U cmp_user -d cmp_db -c "\dt"
```

### Step 2: Deploy Services

```bash
# Make deployment script executable
sudo chmod +x deploy/production/scripts/deploy.sh

# Run deployment
sudo ./deploy/production/scripts/deploy.sh
```

This script will:
- Check prerequisites
- Run database migrations
- Initialize Vault
- Start all services
- Start Nginx
- Perform health checks

### Step 3: Verify Deployment

```bash
# Check service status
sudo systemctl status cmp-inventory
sudo systemctl status cmp-issuer
sudo systemctl status cmp-adapter
sudo systemctl status cmp-auth
sudo systemctl status cmp-nginx

# Check health endpoints
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
curl http://localhost:8084/health

# Check API endpoint
curl -k https://api.cmp.example.com/health
```

### Step 4: Enable Services

```bash
# Enable all services to start on boot
sudo systemctl enable cmp-inventory
sudo systemctl enable cmp-issuer
sudo systemctl enable cmp-adapter
sudo systemctl enable cmp-auth
sudo systemctl enable cmp-nginx
```

---

## Monitoring & Logging

### Systemd Journal

View logs:

```bash
# View all CMP logs
sudo journalctl -u cmp-* -f

# View specific service
sudo journalctl -u cmp-inventory -f

# View logs with filters
sudo journalctl -u cmp-* --since "1 hour ago" --priority err
```

### Prometheus Setup

```bash
# Install Prometheus
sudo apt-get install -y prometheus

# Configure Prometheus
sudo nano /etc/prometheus/prometheus.yml
```

Add scrape configs:

```yaml
scrape_configs:
  - job_name: 'cmp-services'
    static_configs:
      - targets: ['localhost:9091', 'localhost:9092', 'localhost:9093', 'localhost:9094']
```

### Grafana Setup

```bash
# Install Grafana
sudo apt-get install -y grafana

# Start Grafana
sudo systemctl enable grafana-server
sudo systemctl start grafana-server

# Access at http://localhost:3000
# Default credentials: admin/admin
```

### Log Aggregation

Configure centralized logging:

```bash
# Install rsyslog
sudo apt-get install -y rsyslog

# Configure remote logging (if applicable)
sudo nano /etc/rsyslog.d/cmp.conf
```

---

## Maintenance

### Service Management

```bash
# Start service
sudo systemctl start cmp-inventory

# Stop service
sudo systemctl stop cmp-inventory

# Restart service
sudo systemctl restart cmp-inventory

# Reload configuration
sudo systemctl reload cmp-inventory

# Check status
sudo systemctl status cmp-inventory

# View logs
sudo journalctl -u cmp-inventory -f
```

### Updates

```bash
# Pull latest code
cd /opt/cmp-platform
sudo git pull

# Rebuild services
cd backend
go mod download
go build -o /opt/cmp/bin/inventory-service ./cmd/inventory-service
go build -o /opt/cmp/bin/issuer-service ./cmd/issuer-service
go build -o /opt/cmp/bin/adapter-service ./cmd/adapter-service
go build -o /opt/cmp/bin/auth-service ./cmd/auth-service

# Rebuild frontend
cd ../frontend/webapp
npm ci
npm run build
cp -r dist/* /opt/cmp/frontend/dist/

# Restart services
sudo systemctl restart cmp-*
```

### Backup

```bash
# Database backup
sudo -u postgres pg_dump cmp_db > /backup/cmp_db_$(date +%Y%m%d).sql

# Configuration backup
sudo tar -czf /backup/cmp_config_$(date +%Y%m%d).tar.gz /etc/cmp

# Vault backup
vault operator backup -address=http://127.0.0.1:8200 /backup/vault_$(date +%Y%m%d).bak
```

---

## Troubleshooting

### Service Won't Start

```bash
# Check logs
sudo journalctl -u cmp-inventory -n 50

# Check configuration
sudo -u cmp /opt/cmp/bin/inventory-service --help

# Check permissions
ls -la /opt/cmp/bin/
ls -la /etc/cmp/
```

### Database Connection Issues

```bash
# Test connection
psql -h localhost -U cmp_user -d cmp_db

# Check PostgreSQL status
sudo systemctl status postgresql

# Check logs
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

### SSL Certificate Issues

```bash
# Check certificate validity
openssl x509 -in /etc/cmp/ssl/fullchain.pem -text -noout

# Test SSL connection
openssl s_client -connect api.cmp.example.com:443

# Renew Let's Encrypt certificate
sudo certbot renew
```

### High Resource Usage

```bash
# Check system resources
htop

# Check service resource usage
systemctl status cmp-inventory --full
journalctl -u cmp-inventory | grep -i "error\|warning"

# Check database queries
sudo -u postgres psql -d cmp_db -c "SELECT * FROM pg_stat_activity;"
```

---

## Support

For issues and questions:
- Check logs: `sudo journalctl -u cmp-*`
- Review documentation: `/opt/cmp-platform/docs/`
- Contact support: [support@cmp-platform.com]

---

**Last Updated:** 2024-01-15  
**Maintained By:** CMP Platform Team

