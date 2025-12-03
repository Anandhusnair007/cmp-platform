#!/bin/bash
set -euo pipefail

# SSL/TLS Certificate Setup Script for CMP Platform
# Supports Let's Encrypt and custom CA certificates

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

# Setup Let's Encrypt certificates
setup_letsencrypt() {
    log "Setting up Let's Encrypt certificates..."
    
    read -p "Enter API domain (e.g., api.cmp.example.com): " API_DOMAIN
    read -p "Enter App domain (e.g., app.cmp.example.com): " APP_DOMAIN
    read -p "Enter email for Let's Encrypt notifications: " EMAIL
    
    CERT_DIR="/etc/cmp/ssl"
    mkdir -p "$CERT_DIR"
    mkdir -p /var/www/certbot
    
    # Install certbot if not installed
    if ! command -v certbot &> /dev/null; then
        info "Installing certbot..."
        apt-get update -qq
        apt-get install -y -qq certbot python3-certbot-nginx || error "Failed to install certbot"
    fi
    
    # Obtain certificates
    certbot certonly --nginx \
        -d "$API_DOMAIN" \
        -d "$APP_DOMAIN" \
        --email "$EMAIL" \
        --agree-tos \
        --non-interactive \
        --keep-until-expiring || error "Failed to obtain certificates"
    
    # Copy certificates to CMP directory
    cp "/etc/letsencrypt/live/$API_DOMAIN/fullchain.pem" "$CERT_DIR/fullchain.pem"
    cp "/etc/letsencrypt/live/$API_DOMAIN/privkey.pem" "$CERT_DIR/privkey.pem"
    cp "/etc/letsencrypt/live/$API_DOMAIN/chain.pem" "$CERT_DIR/chain.pem"
    
    chmod 600 "$CERT_DIR"/*.pem
    chown root:root "$CERT_DIR"/*.pem
    
    # Setup auto-renewal
    (crontab -l 2>/dev/null; echo "0 2 * * * certbot renew --quiet --deploy-hook 'systemctl reload cmp-nginx'") | crontab -
    
    log "Let's Encrypt certificates installed successfully"
}

# Setup custom CA certificates
setup_custom_certs() {
    log "Setting up custom CA certificates..."
    
    CERT_DIR="/etc/cmp/ssl"
    mkdir -p "$CERT_DIR"
    
    read -p "Enter path to certificate chain file (.crt/.pem): " CERT_FILE
    read -p "Enter path to private key file (.key): " KEY_FILE
    read -p "Enter path to CA chain file (optional): " CHAIN_FILE
    
    if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
        error "Certificate files not found"
    fi
    
    cp "$CERT_FILE" "$CERT_DIR/fullchain.pem"
    cp "$KEY_FILE" "$CERT_DIR/privkey.pem"
    
    if [ -n "$CHAIN_FILE" ] && [ -f "$CHAIN_FILE" ]; then
        cp "$CHAIN_FILE" "$CERT_DIR/chain.pem"
    else
        cp "$CERT_FILE" "$CERT_DIR/chain.pem"
    fi
    
    chmod 600 "$CERT_DIR"/*.pem
    chown root:root "$CERT_DIR"/*.pem
    
    log "Custom certificates installed successfully"
}

# Generate self-signed certificates (testing only)
generate_self_signed() {
    warning "Generating self-signed certificates for TESTING ONLY"
    warning "DO NOT use in production!"
    
    read -p "Enter domain name: " DOMAIN
    
    CERT_DIR="/etc/cmp/ssl"
    mkdir -p "$CERT_DIR"
    
    openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
        -keyout "$CERT_DIR/privkey.pem" \
        -out "$CERT_DIR/fullchain.pem" \
        -subj "/C=US/ST=State/L=City/O=CMP/CN=$DOMAIN"
    
    cp "$CERT_DIR/fullchain.pem" "$CERT_DIR/chain.pem"
    
    chmod 600 "$CERT_DIR"/*.pem
    chown root:root "$CERT_DIR"/*.pem
    
    log "Self-signed certificates generated (TESTING ONLY)"
}

# Main function
main() {
    log "=========================================="
    log "CMP Platform SSL/TLS Certificate Setup"
    log "=========================================="
    
    check_root
    
    echo ""
    echo "Select certificate type:"
    echo "1) Let's Encrypt (Recommended for production)"
    echo "2) Custom CA certificates"
    echo "3) Self-signed (Testing only)"
    echo ""
    read -p "Enter choice [1-3]: " choice
    
    case $choice in
        1)
            setup_letsencrypt
            ;;
        2)
            setup_custom_certs
            ;;
        3)
            generate_self_signed
            ;;
        *)
            error "Invalid choice"
            ;;
    esac
    
    log "=========================================="
    log "SSL/TLS setup completed!"
    log "=========================================="
    info "Update nginx configuration with certificate paths if needed"
}

main "$@"

