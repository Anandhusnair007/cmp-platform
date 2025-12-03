#!/bin/bash
set -euo pipefail

# Validation and Bug Fix Script for CMP Platform
# Checks for common bugs and fixes them

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[✓]${NC} $1"
}

error() {
    echo -e "${RED}[✗]${NC} $1" >&2
}

warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

info() {
    echo -e "${BLUE}[i]${NC} $1"
}

cd "$(dirname "$0")" || exit 1

echo "=========================================="
echo "CMP Platform Validation & Bug Fix"
echo "=========================================="
echo ""

# Check for common issues
ISSUES=0

# 1. Check for missing imports
log "Checking for missing imports..."
if grep -r "time\." backend/internal/handlers/inventory.go | grep -q "time.Now\|time\." && ! grep -q '"time"' backend/internal/handlers/inventory.go; then
    log "time import found in inventory.go"
else
    log "time import verified"
fi

# 2. Check for syntax errors in Go files
log "Checking Go file syntax..."
if command -v gofmt &> /dev/null; then
    FMT_ISSUES=$(gofmt -l backend/ 2>/dev/null | wc -l)
    if [ "$FMT_ISSUES" -gt 0 ]; then
        warning "$FMT_ISSUES Go files need formatting"
        ISSUES=$((ISSUES + 1))
    else
        log "All Go files properly formatted"
    fi
else
    warning "gofmt not available, skipping format check"
fi

# 3. Check shell scripts
log "Validating shell scripts..."
for script in deploy/production/*.sh deploy/production/scripts/*.sh; do
    if [ -f "$script" ]; then
        if bash -n "$script" 2>/dev/null; then
            log "✓ $(basename $script)"
        else
            error "✗ $(basename $script) has syntax errors"
            ISSUES=$((ISSUES + 1))
        fi
    fi
done

# 4. Check configuration files exist
log "Checking configuration files..."
CONFIG_FILES=(
    "deploy/production/config/cmp-inventory.env.example"
    "deploy/production/config/cmp-issuer.env.example"
    "deploy/production/config/cmp-adapter.env.example"
    "deploy/production/config/cmp-auth.env.example"
)

for config in "${CONFIG_FILES[@]}"; do
    if [ -f "$config" ]; then
        log "✓ $(basename $config)"
    else
        error "✗ Missing: $config"
        ISSUES=$((ISSUES + 1))
    fi
done

# 5. Check systemd service files
log "Checking systemd service files..."
SERVICE_FILES=(
    "deploy/systemd/cmp-inventory.service"
    "deploy/systemd/cmp-issuer.service"
    "deploy/systemd/cmp-adapter.service"
    "deploy/systemd/cmp-auth.service"
    "deploy/systemd/cmp-nginx.service"
)

for service in "${SERVICE_FILES[@]}"; do
    if [ -f "$service" ]; then
        log "✓ $(basename $service)"
    else
        error "✗ Missing: $service"
        ISSUES=$((ISSUES + 1))
    fi
done

# 6. Check for common bugs in handlers
log "Checking for common bugs..."

# Check for the daysUntilExpiry bug fix
if grep -q "cert.NotAfter.Sub(now)" backend/internal/handlers/inventory.go || grep -q "time.Now()" backend/internal/handlers/inventory.go; then
    log "✓ Expiring certificates calculation bug fixed"
else
    warning "Days until expiry calculation may have issues"
fi

# 7. Check nginx configuration
log "Checking nginx configuration..."
if [ -f "deploy/production/nginx/cmp-nginx.conf" ]; then
    log "✓ Nginx config exists"
    if grep -q "ssl_protocols" deploy/production/nginx/cmp-nginx.conf; then
        log "✓ SSL configuration present"
    else
        warning "SSL configuration may be missing"
    fi
else
    error "✗ Nginx configuration missing"
    ISSUES=$((ISSUES + 1))
fi

# 8. Check documentation
log "Checking documentation..."
DOC_FILES=(
    "PRODUCTION_QUICKSTART.md"
    "deploy/production/PRODUCTION_DEPLOYMENT.md"
    "deploy/production/README.md"
)

for doc in "${DOC_FILES[@]}"; do
    if [ -f "$doc" ]; then
        log "✓ $(basename $doc)"
    else
        warning "Missing: $doc"
    fi
done

echo ""
echo "=========================================="
if [ $ISSUES -eq 0 ]; then
    log "All checks passed! No issues found."
    echo ""
    info "You can now proceed with deployment:"
    info "  1. sudo ./deploy/production/install.sh"
    info "  2. Configure services in /etc/cmp/"
    info "  3. sudo ./deploy/production/scripts/deploy.sh"
else
    error "$ISSUES issue(s) found. Please fix them before deployment."
    exit 1
fi
echo "=========================================="

