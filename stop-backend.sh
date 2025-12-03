#!/bin/bash

# Stop all CMP backend services

echo "🛑 Stopping CMP Backend Services..."

if [ -d "logs" ]; then
    for pidfile in logs/*.pid; do
        if [ -f "$pidfile" ]; then
            PID=$(cat "$pidfile")
            SERVICE=$(basename "$pidfile" .pid)
            if ps -p $PID > /dev/null 2>&1; then
                echo "   Stopping $SERVICE (PID: $PID)..."
                kill $PID 2>/dev/null || true
            fi
            rm "$pidfile"
        fi
    done
fi

# Also kill any remaining go run processes for our services
pkill -f "go run ./cmd/issuer-service" || true
pkill -f "go run ./cmd/inventory-service" || true
pkill -f "go run ./cmd/adapter-service" || true

echo "✅ All backend services stopped!"
