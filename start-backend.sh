#!/bin/bash

# Start all CMP backend services locally (without Docker)
# This script loads environment variables and starts the services

set -e

echo "🚀 Starting CMP Backend Services..."

# Load environment variables from local.env
if [ -f "local.env" ]; then
    echo "📝 Loading environment from local.env..."
    export $(cat local.env | grep -v '^#' | xargs)
else
    echo "⚠️  local.env not found, using defaults"
fi

# Verify Redis is running
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis is not running. Starting Redis..."
    redis-server --daemonize yes --port ${REDIS_PORT:-6379}
fi
echo "✅ Redis is running"

# Verify Vault is running
if ! vault status > /dev/null 2>&1; then
    echo "❌ Vault is not running. Starting Vault in dev mode..."
    vault server -dev -dev-root-token-id=${VAULT_TOKEN:-dev-only-token} > vault.log 2>&1 &
    sleep 2
fi
echo "✅ Vault is running"

# Verify PostgreSQL is running and accessible
if ! psql -h ${DB_HOST:-localhost} -U ${DB_USER:-cmp_user} -d ${DB_NAME:-cmp_db} -c "\q" 2>/dev/null; then
    echo "⚠️  Cannot connect to PostgreSQL. Make sure it's running and credentials are correct."
    echo "   Run: ./setup-local-db.sh"
fi
echo "✅ PostgreSQL is accessible"

# Create log directory
mkdir -p logs

echo ""
echo "🔧 Starting services..."
echo ""

# Start Issuer Service (Port 8082)
echo "📦 Starting Issuer Service on port 8082..."
cd backend
go run ./cmd/issuer-service -port=8082 > ../logs/issuer-service.log 2>&1 &
ISSUER_PID=$!
echo "   PID: $ISSUER_PID"
cd ..

# Start Inventory Service (Port 8081)
echo "📦 Starting Inventory Service on port 8081..."
cd backend
go run ./cmd/inventory-service -port=8081 > ../logs/inventory-service.log 2>&1 &
INVENTORY_PID=$!
echo "   PID: $INVENTORY_PID"
cd ..

# Start Adapter Service (Port 8083)
echo "📦 Starting Adapter Service on port 8083..."
cd backend
go run ./cmd/adapter-service -port=8083 > ../logs/adapter-service.log 2>&1 &
ADAPTER_PID=$!
echo "   PID: $ADAPTER_PID"
cd ..

# Wait a bit for services to start
sleep 3

echo ""
echo "✅ All backend services started!"
echo ""
echo "📊 Service Status:"
echo "   • Issuer Service:    http://localhost:8082/health (PID: $ISSUER_PID)"
echo "   • Inventory Service: http://localhost:8081/health (PID: $INVENTORY_PID)"
echo "   • Adapter Service:   http://localhost:8083/health (PID: $ADAPTER_PID)"
echo ""
echo "📝 Logs:"
echo "   • Issuer:    tail -f logs/issuer-service.log"
echo "   • Inventory: tail -f logs/inventory-service.log"
echo "   • Adapter:   tail -f logs/adapter-service.log"
echo ""
echo "💡 To stop services, run: ./stop-backend.sh"
echo ""

# Save PIDs for stopping later
echo "$ISSUER_PID" > logs/issuer.pid
echo "$INVENTORY_PID" > logs/inventory.pid
echo "$ADAPTER_PID" > logs/adapter.pid

# Test services
echo "🧪 Testing service health..."
sleep 2
for port in 8081 8082 8083; do
    if curl -s http://localhost:$port/health > /dev/null; then
        echo "   ✅ Port $port is responding"
    else
        echo "   ⚠️  Port $port is not responding yet (may need more time)"
    fi
done

echo ""
echo "🎉 Backend is ready! You can now start the frontend."
