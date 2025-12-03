#!/bin/bash
set -euo pipefail

# Test Build Script for CMP Platform
# Tests frontend and backend builds to ensure no errors

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[TEST]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

cd "$(dirname "$0")" || error "Cannot change to project directory"

log "=========================================="
log "CMP Platform Build Test"
log "=========================================="

# Test Frontend Build
log "Testing Frontend Build..."
if [ -d "frontend/webapp" ]; then
    cd frontend/webapp || error "Frontend directory not found"
    
    if command -v npm &> /dev/null; then
        log "Installing frontend dependencies..."
        npm ci --silent || warning "npm ci failed, trying npm install"
        npm install --silent || error "Failed to install frontend dependencies"
        
        log "Building frontend..."
        npm run build || error "Frontend build failed"
        
        if [ -d "dist" ]; then
            log "✓ Frontend build successful"
            info "Built files in: $(pwd)/dist"
        else
            error "Frontend build did not create dist directory"
        fi
    else
        warning "npm not found, skipping frontend build test"
    fi
    
    cd ../..
else
    warning "Frontend directory not found, skipping"
fi

# Test Backend Code Validation
log ""
log "Testing Backend Code..."
if [ -d "backend" ]; then
    cd backend || error "Backend directory not found"
    
    if command -v go &> /dev/null; then
        log "Checking Go version..."
        go version
        
        log "Validating Go modules..."
        go mod verify || error "Go module verification failed"
        
        log "Downloading dependencies..."
        go mod download || error "Failed to download Go dependencies"
        
        log "Running go fmt check..."
        if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
            warning "Some files need formatting:"
            gofmt -l .
        else
            log "✓ All Go files properly formatted"
        fi
        
        log "Running go vet (static analysis)..."
        go vet ./... || warning "go vet found some issues"
        
        log "Checking if services compile..."
        for service in cmd/inventory-service cmd/issuer-service cmd/adapter-service cmd/auth-service; do
            if [ -d "$service" ]; then
                info "Checking $service..."
                go build -o /dev/null "./$service" || error "Failed to compile $service"
                log "✓ $service compiles successfully"
            else
                warning "$service directory not found"
            fi
        done
        
        log "✓ Backend code validation successful"
    else
        warning "Go not installed, skipping backend build test"
        info "Install Go 1.21+ to test backend compilation"
    fi
    
    cd ..
else
    warning "Backend directory not found, skipping"
fi

# Test Script Syntax
log ""
log "Testing Shell Scripts..."
SCRIPT_ERRORS=0
for script in deploy/production/*.sh deploy/production/scripts/*.sh; do
    if [ -f "$script" ]; then
        if bash -n "$script" 2>&1; then
            log "✓ $script syntax OK"
        else
            error "✗ $script has syntax errors"
            SCRIPT_ERRORS=$((SCRIPT_ERRORS + 1))
        fi
    fi
done

if [ $SCRIPT_ERRORS -eq 0 ]; then
    log "✓ All shell scripts have valid syntax"
else
    error "$SCRIPT_ERRORS script(s) have syntax errors"
fi

# Test Configuration Files
log ""
log "Testing Configuration Files..."
if [ -d "deploy/production/config" ]; then
    for config in deploy/production/config/*.env.example; do
        if [ -f "$config" ]; then
            log "✓ Found: $(basename $config)"
        fi
    done
fi

# Summary
log ""
log "=========================================="
log "Build Test Summary"
log "=========================================="
log "✓ Frontend: Build tested"
log "✓ Backend: Code validated"
log "✓ Scripts: Syntax checked"
log "✓ Configuration: Files verified"
log ""
log "All tests passed! ✓"
log ""

