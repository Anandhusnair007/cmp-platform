# CMP Platform - Production Deployment Checklist

Use this checklist to ensure all steps are completed for a successful production deployment.

## Pre-Deployment

### System Requirements
- [ ] Server with Ubuntu 20.04+ / Debian 11+ / RHEL 8+ / CentOS 8+
- [ ] Minimum 4 CPU cores available
- [ ] Minimum 8GB RAM available (16GB+ recommended)
- [ ] Minimum 100GB disk space available
- [ ] Static IP address assigned
- [ ] DNS records configured:
  - [ ] api.cmp.example.com → Server IP
  - [ ] app.cmp.example.com → Server IP

### Network Configuration
- [ ] Port 80 open (HTTP)
- [ ] Port 443 open (HTTPS)
- [ ] Port 22 open (SSH)
- [ ] Port 5432 restricted (PostgreSQL - localhost only)
- [ ] Port 6379 restricted (Redis - localhost only)
- [ ] Port 8200 restricted (Vault - localhost only)

## Installation Phase

### Step 1: System Preparation
- [ ] System updated: `apt-get update && apt-get upgrade -y`
- [ ] Repository cloned to `/opt/cmp-platform`
- [ ] Scripts made executable

### Step 2: Installation
- [ ] Installation script executed: `./deploy/production/install.sh`
- [ ] All dependencies installed successfully
- [ ] CMP user created
- [ ] Directories created with proper permissions
- [ ] PostgreSQL installed and configured
- [ ] Redis installed and configured
- [ ] Vault installed and configured
- [ ] Backend services built successfully
- [ ] Frontend application built successfully
- [ ] Systemd services installed

## Configuration Phase

### Step 3: Database Configuration
- [ ] PostgreSQL running and accessible
- [ ] Database `cmp_db` created
- [ ] User `cmp_user` created
- [ ] Database password set (strong, 32+ characters)
- [ ] SSL enabled for PostgreSQL connections
- [ ] `pg_hba.conf` configured for secure authentication

### Step 4: Redis Configuration
- [ ] Redis running and accessible
- [ ] Redis password set (strong password)
- [ ] Redis bind to localhost only
- [ ] Protected mode enabled

### Step 5: Vault Configuration
- [ ] Vault service created and configured
- [ ] Vault initialized
- [ ] Unseal keys saved securely (3 keys)
- [ ] Root token saved securely
- [ ] Vault unsealed and operational
- [ ] Vault PKI engine configured (if needed)

### Step 6: Service Configuration
- [ ] Configuration files copied from templates
- [ ] `cmp-inventory.env` configured:
  - [ ] Database password set
  - [ ] Redis password set
  - [ ] Vault token set
- [ ] `cmp-issuer.env` configured:
  - [ ] Database password set
  - [ ] Redis password set
  - [ ] Vault token set
- [ ] `cmp-adapter.env` configured:
  - [ ] Database password set
  - [ ] Redis password set
  - [ ] Vault token set
- [ ] `cmp-auth.env` configured:
  - [ ] Database password set
  - [ ] Redis password set
  - [ ] JWT secret generated (32+ bytes)
  - [ ] OIDC configured (if applicable)

### Step 7: SSL/TLS Configuration
- [ ] SSL setup script executed
- [ ] Let's Encrypt certificates obtained OR
- [ ] Custom CA certificates installed
- [ ] Certificates located in `/etc/cmp/ssl/`
- [ ] Certificate permissions set (600)
- [ ] Auto-renewal configured (if Let's Encrypt)

### Step 8: Nginx Configuration
- [ ] Nginx configuration updated with correct domains
- [ ] SSL certificate paths configured
- [ ] Configuration tested: `nginx -t`
- [ ] Upstream services configured correctly

## Security Phase

### Step 9: Security Hardening
- [ ] Security hardening script executed
- [ ] Kernel parameters hardened
- [ ] System limits configured
- [ ] SSH hardened:
  - [ ] Root login disabled
  - [ ] Password authentication disabled
  - [ ] Key-based authentication enabled
- [ ] Fail2ban configured and enabled
- [ ] Audit logging configured
- [ ] Secure passwords generated and saved
- [ ] Firewall configured:
  - [ ] UFW enabled OR
  - [ ] Firewalld configured
  - [ ] Only required ports open

### Step 10: Additional Security
- [ ] All default passwords changed
- [ ] SSH keys configured for access
- [ ] Automatic security updates enabled
- [ ] File integrity monitoring configured (AIDE)
- [ ] Security audit completed

## Deployment Phase

### Step 11: Database Migrations
- [ ] Migrate tool installed
- [ ] Database migrations executed: `make migrate-up`
- [ ] Schema verified in database
- [ ] Initial data loaded (if applicable)

### Step 12: Service Deployment
- [ ] Deployment script executed: `./deploy/production/scripts/deploy.sh`
- [ ] All services started successfully:
  - [ ] cmp-inventory service
  - [ ] cmp-issuer service
  - [ ] cmp-adapter service
  - [ ] cmp-auth service
- [ ] Nginx started successfully
- [ ] All services enabled on boot

### Step 13: Health Checks
- [ ] Inventory service health check: `curl http://localhost:8081/health`
- [ ] Issuer service health check: `curl http://localhost:8082/health`
- [ ] Adapter service health check: `curl http://localhost:8083/health`
- [ ] Auth service health check: `curl http://localhost:8084/health`
- [ ] API endpoint check: `curl https://api.cmp.example.com/health`
- [ ] Frontend accessible: `https://app.cmp.example.com`

## Monitoring Phase

### Step 14: Monitoring Setup (Optional)
- [ ] Monitoring setup script executed
- [ ] Prometheus installed and configured
- [ ] Grafana installed and configured
- [ ] Node Exporter installed
- [ ] Alert rules configured
- [ ] Grafana admin password changed
- [ ] Dashboards imported
- [ ] Log aggregation configured

### Step 15: Logging Configuration
- [ ] Log rotation configured
- [ ] Centralized logging setup (if applicable)
- [ ] Log retention policy set

## Post-Deployment

### Step 16: Testing
- [ ] API endpoints tested
- [ ] Frontend application functional
- [ ] Authentication working
- [ ] Certificate request flow tested
- [ ] Database operations tested
- [ ] Monitoring dashboards reviewed

### Step 17: Documentation
- [ ] Configuration documented
- [ ] Passwords securely stored
- [ ] Access credentials documented
- [ ] Runbooks created
- [ ] Team trained

### Step 18: Backup Strategy
- [ ] Database backup script created
- [ ] Configuration backup script created
- [ ] Backup schedule configured (cron)
- [ ] Backup location verified
- [ ] Backup restoration tested

### Step 19: Maintenance Plan
- [ ] Update procedure documented
- [ ] Maintenance window scheduled
- [ ] Rollback procedure documented
- [ ] Emergency contacts listed

## Production Readiness Verification

### Security Verification
- [ ] No default passwords in use
- [ ] SSL/TLS configured and working
- [ ] Firewall rules active
- [ ] Security hardening applied
- [ ] Audit logging active
- [ ] Fail2ban active

### Service Verification
- [ ] All services running
- [ ] Services restart on failure
- [ ] Services start on boot
- [ ] Logs being generated
- [ ] Health endpoints responding

### Network Verification
- [ ] DNS resolving correctly
- [ ] SSL certificates valid
- [ ] HTTPS redirect working
- [ ] API accessible externally
- [ ] Frontend accessible externally

### Performance Verification
- [ ] Response times acceptable
- [ ] Resource usage normal
- [ ] No memory leaks
- [ ] Database queries optimized

## Final Sign-Off

- [ ] All checklist items completed
- [ ] Production deployment verified
- [ ] Team notified
- [ ] Documentation updated
- [ ] Monitoring alerts configured
- [ ] Support contacts updated

---

**Deployment Completed:** _______________  
**Verified By:** _______________  
**Date:** _______________

---

## Emergency Contacts

- **On-Call Engineer:** _______________
- **Database Administrator:** _______________
- **Network Administrator:** _______________
- **Security Team:** _______________

## Quick Reference

```bash
# Service management
systemctl status cmp-*
systemctl restart cmp-*

# View logs
journalctl -u cmp-* -f

# Check health
curl http://localhost:8081/health

# Database backup
pg_dump cmp_db > backup.sql
```

---

**Ready for Production!** ✅

