# 🚀 Quick Start Guide

## Start Everything

### Option 1: Docker Compose (Recommended)

```bash
cd cmp-platform
docker-compose -f deploy/docker-compose.yml up --build
```

### Option 2: Local Development

#### Backend Services

```bash
# Terminal 1 - Auth Service
cd backend
go run ./cmd/auth-service

# Terminal 2 - Issuer Service  
cd backend
go run ./cmd/issuer-service

# Terminal 3 - Inventory Service
cd backend
go run ./cmd/inventory-service

# Terminal 4 - Adapter Service
cd backend
go run ./cmd/adapter-service
```

#### Database Setup

```bash
# Apply migrations
make migrate-up

# Initialize Vault PKI
./deploy/vault-init.sh
```

#### Frontend

```bash
cd frontend/webapp
npm install
npm run dev
```

## Access Points

- **Frontend**: http://localhost:3000
- **API**: http://localhost:8082/api/v1
- **Vault UI**: http://localhost:8200 (token: `dev-only-token`)
- **Auth Service**: http://localhost:8084

## Login Credentials

- **Email**: `admin@example.com`
- **Password**: `admin`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user
- `POST /api/v1/auth/logout` - Logout

### Certificates
- `GET /api/v1/certs` - List certificates
- `GET /api/v1/certs/{id}` - Get certificate details
- `POST /api/v1/certs/request` - Request certificate
- `POST /api/v1/certs/{id}/revoke` - Revoke certificate

### Inventory
- `GET /api/v1/inventory` - Get inventory
- `GET /api/v1/inventory/expiring?days=30` - Get expiring certificates

### Agents
- `GET /api/v1/agents` - List agents
- `POST /api/v1/agents/{id}/install` - Install certificate

## Features

✅ JWT Authentication
✅ Dark Theme UI
✅ Certificate Management
✅ Agent Management
✅ Inventory Tracking
✅ Auto-refresh
✅ Toast Notifications
✅ Loading States

## Troubleshooting

### Database Connection Failed
- Check PostgreSQL is running
- Verify connection string in environment variables

### 401 Unauthorized
- Login again to get new token
- Check token is being sent in Authorization header

### Frontend Can't Connect
- Verify backend services are running
- Check API_URL in frontend environment

## Next Steps

1. Customize adapter configurations
2. Add more agents
3. Request certificates
4. Monitor expiring certificates
5. Deploy to production

Enjoy your Certificate Management Platform! 🎉
