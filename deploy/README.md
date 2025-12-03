# Deployment Directory

This directory contains all deployment configurations for hosting the CMP Platform.

## 🚀 Quick Deploy

### Development
```bash
docker-compose -f docker-compose.yml up
```

### Production
```bash
./start-production.sh
```

## 📁 Files

- `docker-compose.yml` - Standard development deployment
- `docker-compose.ha.yml` - High availability production deployment
- `start-production.sh` - Production deployment script
- `production.env.example` - Production environment template
- `nginx/` - Load balancer configuration
- `postgres/` - Database replication config
- `grafana/` - Monitoring dashboards
- `elk/` - Logging stack configuration
- `prometheus/` - Metrics configuration

## 🌐 Hosting Options

1. **Docker Compose** (VPS/Dedicated)
2. **Kubernetes** (Cloud: AWS/Azure/GCP)
3. **Cloud Platforms** (Heroku, DigitalOcean App Platform)

See `HOSTING_GUIDE.md` for detailed instructions.
