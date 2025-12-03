# 🚀 Deployment Guide - Host Your Enterprise CMP Platform

## ✅ **YES - You Can Host This Application!**

The platform is fully configured and ready for deployment. Choose your deployment method below.

---

## 🎯 **Deployment Options**

### Option 1: Docker Compose (Recommended for Production)

#### Standard Deployment
```bash
# 1. Clone/navigate to project
cd cmp-platform

# 2. Start all services
docker-compose -f deploy/docker-compose.yml up -d

# 3. Run database migrations
make migrate-up

# 4. Initialize Vault PKI
./deploy/vault-init.sh

# 5. Access services
# Frontend: http://localhost:3000
# API: http://localhost:8082/api/v1
# Vault: http://localhost:8200
```

#### High Availability Deployment
```bash
# Start HA stack with load balancer
docker-compose -f deploy/docker-compose.ha.yml up -d

# Access points:
# Frontend: http://localhost:3000
# API Gateway: http://localhost:80
# Grafana: http://localhost:3001
# Prometheus: http://localhost:9090
```

---

### Option 2: Kubernetes Deployment

#### Prerequisites
```bash
# Install kubectl and helm
kubectl version --client
helm version
```

#### Deploy with Helm
```bash
# 1. Install Helm chart
helm install cmp ./charts/cmp \
  --namespace cmp \
  --create-namespace \
  --set inventory.replicaCount=2 \
  --set issuer.replicaCount=2

# 2. Check deployment
kubectl get pods -n cmp

# 3. Port forward to access
kubectl port-forward -n cmp svc/cmp-frontend 3000:3000
```

#### Configure Ingress
```bash
# Apply ingress configuration
kubectl apply -f charts/cmp/ingress.yaml

# Access via ingress
# https://cmp.yourdomain.com
```

---

### Option 3: Cloud Platform Deployment

#### AWS Deployment

```bash
# 1. Build and push Docker images
aws ecr create-repository --repository-name cmp-platform

# Build images
docker build -t cmp-inventory:latest ./backend -f Dockerfile.inventory
docker tag cmp-inventory:latest <account>.dkr.ecr.<region>.amazonaws.com/cmp-inventory:latest
docker push <account>.dkr.ecr.<region>.amazonaws.com/cmp-inventory:latest

# 2. Deploy with ECS/EKS
# Use provided Terraform configurations in infra/terraform/
```

#### Azure Deployment

```bash
# Create Azure Container Registry
az acr create --resource-group cmp-rg --name cmpregistry --sku Basic

# Build and push
az acr build --registry cmpregistry --image cmp-inventory:latest ./backend

# Deploy to AKS
az aks create --resource-group cmp-rg --name cmp-cluster
```

#### GCP Deployment

```bash
# Enable GKE
gcloud container clusters create cmp-cluster

# Deploy
kubectl apply -f charts/cmp/
```

---

## 🔧 **Environment Configuration**

### Required Environment Variables

#### Backend Services
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=cmp_user
DB_PASSWORD=your-secure-password
DB_NAME=cmp_db
DB_SSLMODE=require

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Vault
VAULT_ADDR=http://vault:8200
VAULT_TOKEN=your-vault-token

# OIDC (Production)
OIDC_ISSUER_URL=https://keycloak.yourdomain.com/realms/cmp
OIDC_CLIENT_ID=cmp-client
OIDC_CLIENT_SECRET=your-client-secret

# JWT Secret (CHANGE IN PRODUCTION!)
JWT_SECRET=your-super-secret-jwt-key-change-this

# Service Configuration
SERVER_PORT=8082
LOG_LEVEL=info
```

#### Frontend
```bash
REACT_APP_API_URL=https://api.yourdomain.com
REACT_APP_OIDC_ISSUER_URL=https://keycloak.yourdomain.com
```

---

## 📋 **Pre-Deployment Checklist**

### 1. Security Setup
- [ ] Change all default passwords
- [ ] Generate secure JWT secret
- [ ] Configure OIDC provider (Keycloak/Auth0)
- [ ] Set up SSL certificates for production
- [ ] Configure firewall rules
- [ ] Enable encryption at rest

### 2. Database Setup
- [ ] Create PostgreSQL database
- [ ] Run migrations: `make migrate-up`
- [ ] Configure database backups
- [ ] Set up replication (for HA)

### 3. Vault Setup
- [ ] Initialize Vault (production mode)
- [ ] Configure PKI engine
- [ ] Set up unseal keys (secure storage)
- [ ] Configure access policies

### 4. Monitoring
- [ ] Deploy Prometheus
- [ ] Deploy Grafana
- [ ] Configure dashboards
- [ ] Set up alerting

### 5. Logging
- [ ] Deploy ELK stack
- [ ] Configure log shipping
- [ ] Set up log retention

---

## 🌐 **Domain & DNS Configuration**

### Production Domain Setup

1. **Configure DNS Records**
```
api.yourdomain.com     -> Load Balancer IP
app.yourdomain.com     -> Frontend IP
vault.yourdomain.com   -> Vault IP
```

2. **SSL Certificates**
```bash
# Use Let's Encrypt or your CA
certbot certonly --nginx -d api.yourdomain.com
certbot certonly --nginx -d app.yourdomain.com
```

3. **Update Configuration**
```bash
# Update API URL in frontend
export REACT_APP_API_URL=https://api.yourdomain.com

# Update CORS in backend
export ALLOWED_ORIGINS=https://app.yourdomain.com
```

---

## 🗄️ **Database Initialization**

### Production Database Setup

```bash
# 1. Create database
psql -h your-db-host -U postgres
CREATE DATABASE cmp_db;
CREATE USER cmp_user WITH PASSWORD 'secure-password';
GRANT ALL PRIVILEGES ON DATABASE cmp_db TO cmp_user;

# 2. Run migrations
export DB_HOST=your-db-host
export DB_PASSWORD=secure-password
make migrate-up

# 3. Initialize default data
psql -h your-db-host -U cmp_user -d cmp_db -f scripts/init-data.sql
```

---

## 🔐 **Vault Production Setup**

```bash
# 1. Initialize Vault (production mode)
vault operator init -key-shares=5 -key-threshold=3

# 2. Store unseal keys securely
# 3. Unseal Vault
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>

# 4. Enable PKI
vault secrets enable -path=cmp-pki pki

# 5. Use Terraform to configure
cd infra/terraform/vault-pki
terraform init
terraform apply
```

---

## 📊 **Monitoring Setup**

### Start Monitoring Stack

```bash
# Start Prometheus & Grafana
docker-compose -f deploy/docker-compose.ha.yml up -d prometheus grafana

# Access Grafana
# URL: http://localhost:3001
# Login: admin/admin (CHANGE PASSWORD!)

# Import dashboards
# Upload dashboards from deploy/grafana/dashboards/
```

---

## 🚢 **Production Deployment Steps**

### Step-by-Step Production Deployment

1. **Prepare Environment**
```bash
# Clone repository
git clone <repo-url>
cd cmp-platform

# Set environment variables
cp .env.example .env
# Edit .env with production values
```

2. **Build Images**
```bash
# Build all service images
docker-compose -f deploy/docker-compose.ha.yml build

# Or build individually
cd backend && docker build -f Dockerfile.issuer -t cmp-issuer:latest .
```

3. **Deploy Infrastructure**
```bash
# Start infrastructure first
docker-compose -f deploy/docker-compose.ha.yml up -d postgres redis vault

# Wait for services to be ready
sleep 30

# Initialize database
make migrate-up

# Initialize Vault
./deploy/vault-init.sh
```

4. **Deploy Services**
```bash
# Start all services
docker-compose -f deploy/docker-compose.ha.yml up -d

# Check status
docker-compose -f deploy/docker-compose.ha.yml ps
```

5. **Verify Deployment**
```bash
# Check health endpoints
curl http://localhost/health

# Check services
curl http://localhost:8082/health
curl http://localhost:8081/health
```

---

## 🔒 **Security Hardening for Production**

### 1. Update Default Passwords
```bash
# Change database password
# Change Vault root token
# Change Grafana admin password
# Generate new JWT secret
```

### 2. Enable HTTPS
```bash
# Configure SSL/TLS
# Update nginx.conf with SSL certificates
# Enable HSTS
# Configure cipher suites
```

### 3. Network Security
```bash
# Configure firewall
# Restrict service access
# Enable mTLS between services
# Configure network policies (K8s)
```

### 4. Secrets Management
```bash
# Use Vault for all secrets
# Never commit secrets to git
# Rotate secrets regularly
# Use environment variables from secure store
```

---

## 📈 **Scaling Configuration**

### Horizontal Scaling

```bash
# Scale services
docker-compose -f deploy/docker-compose.ha.yml up -d --scale issuer-service=3
docker-compose -f deploy/docker-compose.ha.yml up -d --scale inventory-service=3
```

### Kubernetes Scaling

```bash
# Auto-scaling
kubectl autoscale deployment cmp-issuer -n cmp --min=2 --max=10 --cpu-percent=80

# Manual scaling
kubectl scale deployment cmp-issuer -n cmp --replicas=5
```

---

## 🔄 **Updates & Maintenance**

### Update Application

```bash
# 1. Pull latest code
git pull origin main

# 2. Rebuild images
docker-compose -f deploy/docker-compose.ha.yml build

# 3. Rolling update
docker-compose -f deploy/docker-compose.ha.yml up -d --no-deps issuer-service
```

### Database Migrations

```bash
# Run new migrations
make migrate-up

# Rollback if needed
make migrate-down
```

---

## 🆘 **Troubleshooting**

### Common Issues

1. **Services Not Starting**
```bash
# Check logs
docker-compose logs issuer-service
docker-compose logs postgres

# Check resource usage
docker stats
```

2. **Database Connection Issues**
```bash
# Verify database is running
docker-compose ps postgres

# Test connection
psql -h localhost -U cmp_user -d cmp_db
```

3. **Vault Issues**
```bash
# Check Vault status
vault status

# Check if sealed
vault operator unseal <key>
```

---

## 📞 **Support**

For deployment issues:
- Check logs: `docker-compose logs <service>`
- Review documentation: `docs/`
- Check health endpoints: `/health`

---

## ✅ **Deployment Complete!**

Once deployed, your platform will be accessible at:
- **Frontend**: http://your-domain:3000
- **API**: http://your-domain:8082
- **Grafana**: http://your-domain:3001

**You're ready to host and deploy!** 🚀
