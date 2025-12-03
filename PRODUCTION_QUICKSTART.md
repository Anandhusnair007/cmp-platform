# CMP Platform - Production Quick Start Guide

Enterprise-grade production deployment guide for organizations like Wipro, TCS, Microsoft Security, and Cisco.

## Overview

This guide provides a streamlined path to deploy the Certificate Management Platform in a production environment **without Docker**, using native Linux services with enterprise-grade security.

## Prerequisites

- Linux server (Ubuntu 20.04+, Debian 11+, RHEL 8+, CentOS 8+)
- Root or sudo access
- 4+ CPU cores, 8GB+ RAM, 100GB+ disk
- Static IP address
- Domain names configured (api.cmp.example.com, app.cmp.example.com)

## Quick Deployment (5 Steps)

### Step 1: Clone and Prepare

```bash
# Clone repository
cd /opt
sudo git clone <repository-url> cmp-platform
cd cmp-platform

# Make scripts executable
sudo chmod +x deploy/production/install.sh
sudo chmod +x deploy/production/scripts/*.sh
```

### Step 2: Install System

```bash
# Run installation (15-30 minutes)
sudo ./deploy/production/install.sh
```

This installs:
- All dependencies (Go, Node.js, PostgreSQL, Redis, Vault, Nginx)
- Builds all services
- Creates systemd services
- Sets up directories and permissions

### Step 3: Configure

```bash
# Copy configuration templates
sudo cp deploy/production/config/*.env.example /etc/cmp/

# Generate secure passwords
sudo openssl rand -base64 32  # For DB_PASSWORD
sudo openssl rand -hex 32     # For JWT_SECRET

# Edit configuration files
sudo nano /etc/cmp/cmp-inventory.env
sudo nano /etc/cmp/cmp-issuer.env
sudo nano /etc/cmp/cmp-adapter.env
sudo nano /etc/cmp/cmp-auth.env

# Update database password
sudo -u postgres psql -c "ALTER USER cmp_user WITH PASSWORD 'YOUR_SECURE_PASSWORD';"
```

### Step 4: Setup SSL and Security

```bash
# Setup SSL certificates (Let's Encrypt)
sudo ./deploy/production/scripts/setup-ssl.sh

# Hardening security
sudo ./deploy/production/scripts/harden-system.sh
```

### Step 5: Deploy

```bash
# Deploy all services
sudo ./deploy/production/scripts/deploy.sh
```

## Verify Deployment

```bash
# Check service status
sudo systemctl status cmp-inventory cmp-issuer cmp-adapter cmp-auth cmp-nginx

# Test health endpoints
curl http://localhost:8081/health
curl http://localhost:8082/health
curl https://api.cmp.example.com/health
```

## Access Your Platform

- **Frontend**: https://app.cmp.example.com
- **API**: https://api.cmp.example.com
- **Grafana**: http://localhost:3000 (if monitoring is setup)

## Common Commands

### Service Management

```bash
# Start/stop/restart services
sudo systemctl start|stop|restart cmp-inventory
sudo systemctl start|stop|restart cmp-issuer
sudo systemctl start|stop|restart cmp-adapter
sudo systemctl start|stop|restart cmp-auth
sudo systemctl start|stop|restart cmp-nginx

# Enable services on boot
sudo systemctl enable cmp-*
```

### View Logs

```bash
# All CMP logs
sudo journalctl -u cmp-* -f

# Specific service
sudo journalctl -u cmp-inventory -f

# Recent errors
sudo journalctl -u cmp-* --since "1 hour ago" --priority err
```

### Database Operations

```bash
# Connect to database
sudo -u postgres psql -d cmp_db

# Run migrations
cd /opt/cmp-platform/backend
make migrate-up

# Backup database
sudo -u postgres pg_dump cmp_db > backup.sql
```

## Security Checklist

Before going live:

- [ ] All default passwords changed
- [ ] SSL/TLS certificates installed (Let's Encrypt or custom CA)
- [ ] Firewall configured (ports 80, 443 open)
- [ ] SSH hardened (key-based auth only)
- [ ] Security hardening script run
- [ ] Database SSL enabled
- [ ] Redis password set
- [ ] Vault initialized and unsealed
- [ ] Fail2ban configured
- [ ] Monitoring and alerting setup

## Troubleshooting

### Service Won't Start

```bash
# Check logs
sudo journalctl -u cmp-inventory -n 50

# Check configuration
cat /etc/cmp/cmp-inventory.env

# Test binary
sudo -u cmp /opt/cmp/bin/inventory-service --help
```

### Database Connection Issues

```bash
# Test connection
sudo -u postgres psql -d cmp_db -U cmp_user

# Check PostgreSQL
sudo systemctl status postgresql
sudo tail -f /var/log/postgresql/postgresql-*.log
```

### SSL Issues

```bash
# Check certificate
openssl x509 -in /etc/cmp/ssl/fullchain.pem -text -noout

# Test SSL
openssl s_client -connect api.cmp.example.com:443

# Renew Let's Encrypt
sudo certbot renew
```

## Optional: Setup Monitoring

```bash
# Install Prometheus and Grafana
sudo ./deploy/production/scripts/setup-monitoring.sh

# Access Grafana
# http://localhost:3000
# Default: admin/admin (CHANGE PASSWORD!)
```

## Production Recommendations

1. **High Availability**: Deploy behind load balancer with multiple instances
2. **Backup**: Setup automated database and configuration backups
3. **Monitoring**: Enable Prometheus + Grafana + alerting
4. **Logging**: Configure centralized log aggregation (ELK stack)
5. **Updates**: Schedule regular security updates
6. **Documentation**: Maintain runbooks for your team

## Support

- Full Documentation: `deploy/production/PRODUCTION_DEPLOYMENT.md`
- Architecture: `docs/architecture.md`
- Runbooks: `docs/runbooks.md`

## Next Steps

1. Configure your domain names
2. Setup SSL certificates
3. Review and customize configuration
4. Enable monitoring
5. Setup backups
6. Train your team

---

**Ready for Enterprise Production!** 🚀

For detailed information, see [PRODUCTION_DEPLOYMENT.md](deploy/production/PRODUCTION_DEPLOYMENT.md)

