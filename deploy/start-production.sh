#!/bin/bash
set -e

echo "🚀 Starting CMP Platform - Production Deployment"
echo "=================================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "Please copy deploy/production.env.example to .env and configure it"
    exit 1
fi

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Check required variables
if [ -z "$DB_PASSWORD" ] || [ "$DB_PASSWORD" = "CHANGE-THIS-SECURE-PASSWORD" ]; then
    echo "❌ Error: Please set DB_PASSWORD in .env file"
    exit 1
fi

if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "CHANGE-THIS-GENERATE-32-BYTE-RANDOM-SECRET" ]; then
    echo "❌ Error: Please set JWT_SECRET in .env file"
    exit 1
fi

echo "✅ Environment variables loaded"

# Start infrastructure
echo ""
echo "📦 Starting infrastructure services..."
docker-compose -f deploy/docker-compose.ha.yml up -d postgres redis vault

echo "⏳ Waiting for services to be ready..."
sleep 30

# Check if database is ready
echo "🔍 Checking database connection..."
until docker-compose -f deploy/docker-compose.ha.yml exec -T postgres pg_isready -U cmp_user > /dev/null 2>&1; do
    echo "   Waiting for database..."
    sleep 2
done
echo "✅ Database is ready"

# Run migrations
echo ""
echo "🗄️  Running database migrations..."
export DB_HOST=postgres
export DB_PASSWORD=$DB_PASSWORD
make migrate-up || echo "⚠️  Migration may have already been applied"

# Initialize Vault (if needed)
echo ""
echo "🔐 Initializing Vault PKI..."
./deploy/vault-init.sh || echo "⚠️  Vault may already be initialized"

# Start all services
echo ""
echo "🚀 Starting all services..."
docker-compose -f deploy/docker-compose.ha.yml up -d

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📊 Access Points:"
echo "   Frontend:    http://localhost:3000"
echo "   API Gateway: http://localhost:80"
echo "   Grafana:     http://localhost:3001 (admin/admin)"
echo "   Prometheus:  http://localhost:9090"
echo ""
echo "🔍 Check status:"
echo "   docker-compose -f deploy/docker-compose.ha.yml ps"
echo ""
echo "📝 View logs:"
echo "   docker-compose -f deploy/docker-compose.ha.yml logs -f"
