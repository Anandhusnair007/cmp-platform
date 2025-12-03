#!/bin/bash
set -euo pipefail

# Enterprise Production Deployment Script for CMP Platform
# Designed for Wipro, TCS, Microsoft Security, Cisco-grade deployments

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
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

check_root() {
    if [ "$EUID" -ne 0 ]; then 
        error "Please run as root (use sudo)"
    fi
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if binaries exist
    for service in inventory-service issuer-service adapter-service auth-service; do
        if [ ! -f "/opt/cmp/bin/$service" ]; then
            error "$service binary not found. Run install.sh first."
        fi
    done
    
    # Check configuration files
    for env_file in cmp-inventory.env cmp-issuer.env cmp-adapter.env cmp-auth.env; do
        if [ ! -f "/etc/cmp/$env_file" ]; then
            warning "$env_file not found, using defaults"
        fi
    done
    
    # Check database connectivity
    if ! pg_isready -h localhost -U cmp_user > /dev/null 2>&1; then
        error "Cannot connect to PostgreSQL. Please check database configuration."
    fi
    
    # Check Redis connectivity
    if ! redis-cli -h localhost ping > /dev/null 2>&1; then
        error "Cannot connect to Redis. Please check Redis configuration."
    fi
    
    log "Prerequisites check passed"
}

# Run database migrations
run_migrations() {
    log "Running database migrations..."
    
    cd /opt/cmp/backend || error "Backend directory not found"
    
    # Check if migrate tool is available
    if ! command -v migrate &> /dev/null; then
        warning "migrate tool not found, installing..."
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    fi
    
    # Load database password from config
    DB_PASSWORD=$(grep DB_PASSWORD /etc/cmp/cmp-inventory.env 2>/dev/null | cut -d'=' -f2 || echo "")
    
    if [ -n "$DB_PASSWORD" ]; then
        export DB_PASSWORD
        migrate -path migrations -database "postgres://cmp_user:${DB_PASSWORD}@localhost:5432/cmp_db?sslmode=require" up || warning "Migration may have already been applied"
    else
        warning "Database password not found, skipping migrations"
    fi
    
    log "Migrations completed"
}

# Initialize Vault
init_vault() {
    log "Initializing Vault..."
    
    if systemctl is-active --quiet vault; then
        info "Vault is running"
        
        # Check if already initialized
        if ! vault status > /dev/null 2>&1; then
            warning "Vault needs initialization. Please run: vault operator init"
        else
            log "Vault is already initialized"
        fi
    else
        warning "Vault is not running. Please start it first: systemctl start vault"
    fi
}

# Start services
start_services() {
    log "Starting CMP services..."
    
    services=("cmp-inventory" "cmp-issuer" "cmp-adapter" "cmp-auth")
    
    for service in "${services[@]}"; do
        if systemctl is-active --quiet "$service"; then
            info "$service is already running, reloading..."
            systemctl reload "$service" || systemctl restart "$service"
        else
            systemctl enable "$service"
            systemctl start "$service"
        fi
        
        # Wait for service to be healthy
        sleep 2
        
        if systemctl is-active --quiet "$service"; then
            log "$service started successfully"
        else
            error "$service failed to start. Check logs: journalctl -u $service"
        fi
    done
}

# Start Nginx
start_nginx() {
    log "Starting Nginx reverse proxy..."
    
    # Test configuration
    if nginx -t -c /etc/nginx/cmp/cmp-nginx.conf; then
        systemctl enable cmp-nginx
        systemctl start cmp-nginx
        
        if systemctl is-active --quiet cmp-nginx; then
            log "Nginx started successfully"
        else
            error "Nginx failed to start. Check logs: journalctl -u cmp-nginx"
        fi
    else
        error "Nginx configuration test failed"
    fi
}

# Health check
health_check() {
    log "Performing health checks..."
    
    sleep 5
    
    # Check services
    services=("cmp-inventory" "cmp-issuer" "cmp-adapter" "cmp-auth")
    for service in "${services[@]}"; do
        if systemctl is-active --quiet "$service"; then
            log "$service: ✓ Running"
        else
            error "$service: ✗ Not running"
        fi
    done
    
    # Check HTTP endpoints
    if curl -f -s http://localhost:8081/health > /dev/null; then
        log "Inventory service health: ✓ Healthy"
    else
        warning "Inventory service health: ✗ Unhealthy"
    fi
    
    if curl -f -s http://localhost:8082/health > /dev/null; then
        log "Issuer service health: ✓ Healthy"
    else
        warning "Issuer service health: ✗ Unhealthy"
    fi
    
    log "Health checks completed"
}

# Display status
display_status() {
    log "=========================================="
    log "CMP Platform Deployment Status"
    log "=========================================="
    
    echo ""
    info "Services:"
    systemctl status cmp-inventory --no-pager -l || true
    systemctl status cmp-issuer --no-pager -l || true
    systemctl status cmp-adapter --no-pager -l || true
    systemctl status cmp-auth --no-pager -l || true
    systemctl status cmp-nginx --no-pager -l || true
    
    echo ""
    info "Access Points:"
    info "  - Frontend: https://app.cmp.example.com"
    info "  - API: https://api.cmp.example.com"
    info "  - Health: http://localhost:8081/health"
    
    echo ""
    info "Useful commands:"
    info "  - View logs: journalctl -u cmp-{service-name} -f"
    info "  - Restart service: systemctl restart cmp-{service-name}"
    info "  - Check status: systemctl status cmp-{service-name}"
}

# Main function
main() {
    log "=========================================="
    log "CMP Platform Production Deployment"
    log "Enterprise-Grade Setup"
    log "=========================================="
    
    check_root
    check_prerequisites
    run_migrations
    init_vault
    start_services
    start_nginx
    health_check
    display_status
    
    log "=========================================="
    log "Deployment completed successfully!"
    log "=========================================="
}

main "$@"

