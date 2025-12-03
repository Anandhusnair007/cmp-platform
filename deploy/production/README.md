# CMP Platform - Production Deployment

Enterprise-grade production deployment for the Certificate Management Platform without Docker.

## Quick Start

### 1. Install System

```bash
# Run as root
sudo ./deploy/production/install.sh
```

### 2. Configure Services

```bash
# Copy and edit configuration files
sudo cp deploy/production/config/*.env.example /etc/cmp/
sudo nano /etc/cmp/cmp-inventory.env
sudo nano /etc/cmp/cmp-issuer.env
sudo nano /etc/cmp/cmp-adapter.env
sudo nano /etc/cmp/cmp-auth.env
```

### 3. Setup SSL/TLS

```bash
# Run SSL setup (Let's Encrypt recommended)
sudo ./deploy/production/scripts/setup-ssl.sh
```

### 4. Harden Security

```bash
# Run security hardening
sudo ./deploy/production/scripts/harden-system.sh
```

### 5. Deploy Services

```bash
# Run deployment
sudo ./deploy/production/scripts/deploy.sh
```

### 6. Setup Monitoring (Optional)

```bash
# Setup Prometheus and Grafana
sudo ./deploy/production/scripts/setup-monitoring.sh
```

## Directory Structure

```
/opt/cmp/
├── bin/                    # Service binaries
├── backend/                # Backend code
└── frontend/               # Frontend build

/etc/cmp/
├── cmp-inventory.env       # Inventory service config
├── cmp-issuer.env          # Issuer service config
├── cmp-adapter.env         # Adapter service config
└── cmp-auth.env            # Auth service config

/etc/nginx/cmp/
└── cmp-nginx.conf          # Nginx configuration

/var/log/cmp/               # Application logs
/var/lib/cmp/               # Application data
```

## Service Management

```bash
# Start service
sudo systemctl start cmp-inventory

# Stop service
sudo systemctl stop cmp-inventory

# Restart service
sudo systemctl restart cmp-inventory

# Check status
sudo systemctl status cmp-inventory

# View logs
sudo journalctl -u cmp-inventory -f
```

## Configuration Files

All configuration files are located in `/etc/cmp/`:

- `cmp-inventory.env` - Inventory service configuration
- `cmp-issuer.env` - Issuer service configuration
- `cmp-adapter.env` - Adapter service configuration
- `cmp-auth.env` - Auth service configuration

Copy from examples:

```bash
sudo cp deploy/production/config/*.env.example /etc/cmp/
```

## Security Hardening

The platform includes enterprise-grade security:

- Kernel parameter hardening
- System limits configuration
- SSH hardening
- Fail2ban integration
- Audit logging
- Firewall configuration
- SSL/TLS encryption

Run security hardening:

```bash
sudo ./deploy/production/scripts/harden-system.sh
```

## Monitoring

Monitor the platform using:

- **Systemd Journal**: `journalctl -u cmp-*`
- **Prometheus**: Metrics at `/metrics` endpoint
- **Grafana**: Dashboards for visualization
- **Node Exporter**: System metrics

Setup monitoring:

```bash
sudo ./deploy/production/scripts/setup-monitoring.sh
```

## Documentation

- [Full Deployment Guide](PRODUCTION_DEPLOYMENT.md) - Comprehensive deployment instructions
- [Architecture Documentation](../docs/architecture.md) - System architecture
- [Runbooks](../docs/runbooks.md) - Operational procedures

## Support

For issues:

1. Check logs: `sudo journalctl -u cmp-*`
2. Review configuration: `/etc/cmp/`
3. Check service status: `sudo systemctl status cmp-*`

## Requirements

- Ubuntu 20.04+ / Debian 11+ / RHEL 8+ / CentOS 8+
- 4+ CPU cores
- 8GB+ RAM
- 100GB+ disk space
- Static IP address
- DNS records configured

## Production Checklist

- [ ] All passwords changed from defaults
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules configured
- [ ] Security hardening applied
- [ ] Database backups configured
- [ ] Monitoring setup complete
- [ ] Log aggregation configured
- [ ] Documentation reviewed

## License

MIT License - See LICENSE file for details

