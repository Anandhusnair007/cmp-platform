#!/bin/bash
set -euo pipefail

# Enterprise Production Installation Script for CMP Platform
# Designed for Wipro, TCS, Microsoft Security, Cisco-grade deployments
# Version: 2.0.0

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Installation directories
CMP_HOME="/opt/cmp"
CMP_USER="cmp"
CMP_GROUP="cmp"
SERVICE_USER_HOME="/var/lib/cmp"
LOG_DIR="/var/log/cmp"
CONFIG_DIR="/etc/cmp"
NGINX_CONFIG_DIR="/etc/nginx/cmp"
CERT_DIR="/etc/cmp/ssl"

# Logging function
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

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then 
        error "Please run as root (use sudo)"
    fi
}

# Detect OS
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        OS_VERSION=$VERSION_ID
        info "Detected OS: $OS $OS_VERSION"
    else
        error "Cannot detect OS. Only Ubuntu/Debian/RHEL/CentOS supported."
    fi
}

# Install system dependencies
install_dependencies() {
    log "Installing system dependencies..."
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get update -qq
        apt-get install -y -qq \
            postgresql-client \
            redis-tools \
            nginx \
            curl \
            wget \
            git \
            build-essential \
            ca-certificates \
            gnupg \
            lsb-release \
            systemd \
            logrotate \
            fail2ban \
            ufw \
            openssl \
            certbot \
            python3-certbot-nginx \
            prometheus \
            prometheus-node-exporter \
            nodejs \
            npm \
            golang-go || error "Failed to install dependencies"
    elif [ "$OS" == "rhel" ] || [ "$OS" == "centos" ] || [ "$OS" == "fedora" ]; then
        if [ "$OS" == "centos" ] || [ "$OS" == "rhel" ]; then
            yum install -y -q epel-release
        fi
        yum install -y -q \
            postgresql \
            redis \
            nginx \
            curl \
            wget \
            git \
            gcc \
            make \
            ca-certificates \
            systemd \
            logrotate \
            fail2ban \
            firewalld \
            openssl \
            certbot \
            python3-certbot-nginx \
            nodejs \
            npm \
            golang || error "Failed to install dependencies"
    else
        error "Unsupported OS: $OS"
    fi
    
    log "Dependencies installed successfully"
}

# Create CMP user and directories
create_user_and_directories() {
    log "Creating CMP user and directories..."
    
    # Create user if doesn't exist
    if ! id "$CMP_USER" &>/dev/null; then
        useradd -r -s /bin/false -d "$SERVICE_USER_HOME" -m "$CMP_USER" || error "Failed to create user"
        log "Created user: $CMP_USER"
    else
        warning "User $CMP_USER already exists"
    fi
    
    # Create directories
    mkdir -p "$CMP_HOME"/{bin,backend,frontend}
    mkdir -p "$SERVICE_USER_HOME"
    mkdir -p "$LOG_DIR"
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$NGINX_CONFIG_DIR"
    mkdir -p "$CERT_DIR"
    mkdir -p /var/lib/postgresql/data
    mkdir -p /var/lib/redis
    
    # Set permissions
    chown -R "$CMP_USER:$CMP_GROUP" "$CMP_HOME"
    chown -R "$CMP_USER:$CMP_GROUP" "$SERVICE_USER_HOME"
    chown -R "$CMP_USER:$CMP_GROUP" "$LOG_DIR"
    chmod 750 "$CONFIG_DIR"
    chmod 750 "$CERT_DIR"
    
    log "User and directories created"
}

# Install PostgreSQL
setup_postgresql() {
    log "Setting up PostgreSQL..."
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get install -y -qq postgresql postgresql-contrib
        systemctl enable postgresql
        systemctl start postgresql
    elif [ "$OS" == "rhel" ] || [ "$OS" == "centos" ] || [ "$OS" == "fedora" ]; then
        yum install -y -q postgresql-server postgresql-contrib
        if [ "$OS" == "rhel" ] || [ "$OS" == "centos" ]; then
            postgresql-setup --initdb
        fi
        systemctl enable postgresql
        systemctl start postgresql
    fi
    
    log "PostgreSQL installed and started"
}

# Install Redis
setup_redis() {
    log "Setting up Redis..."
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get install -y -qq redis-server
    elif [ "$OS" == "rhel" ] || [ "$OS" == "centos" ] || [ "$OS" == "fedora" ]; then
        yum install -y -q redis
    fi
    
    systemctl enable redis
    systemctl start redis
    
    log "Redis installed and started"
}

# Install HashiCorp Vault
setup_vault() {
    log "Setting up HashiCorp Vault..."
    
    VAULT_VERSION="1.15.4"
    VAULT_DIR="/opt/vault"
    
    mkdir -p "$VAULT_DIR"
    
    # Download and install Vault
    if [ ! -f "$VAULT_DIR/vault" ]; then
        cd /tmp
        wget -q "https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip"
        unzip -q "vault_${VAULT_VERSION}_linux_amd64.zip"
        mv vault "$VAULT_DIR/vault"
        chmod +x "$VAULT_DIR/vault"
        rm "vault_${VAULT_VERSION}_linux_amd64.zip"
    fi
    
    # Create systemd service for Vault (production mode)
    cat > /etc/systemd/system/vault.service <<EOF
[Unit]
Description=HashiCorp Vault - Secrets Management
After=network.target
Requires=network.target

[Service]
Type=notify
User=vault
Group=vault
WorkingDirectory=/opt/vault
ExecStart=$VAULT_DIR/vault server -config=/etc/vault.d/vault.hcl
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=process
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    
    # Create Vault directories
    mkdir -p /etc/vault.d
    mkdir -p /var/lib/vault
    useradd -r -s /bin/false -d /var/lib/vault vault || true
    chown -R vault:vault /var/lib/vault
    
    log "Vault installed (configuration required)"
}

# Build backend services
build_backend() {
    log "Building backend services..."
    
    cd "$(dirname "$0")/../.." || error "Cannot find project root"
    
    # Install Go dependencies
    cd backend || error "Backend directory not found"
    go mod download || error "Failed to download Go dependencies"
    
    # Build services
    info "Building inventory-service..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o "$CMP_HOME/bin/inventory-service" ./cmd/inventory-service || error "Failed to build inventory-service"
    
    info "Building issuer-service..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o "$CMP_HOME/bin/issuer-service" ./cmd/issuer-service || error "Failed to build issuer-service"
    
    info "Building adapter-service..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o "$CMP_HOME/bin/adapter-service" ./cmd/adapter-service || error "Failed to build adapter-service"
    
    info "Building auth-service..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o "$CMP_HOME/bin/auth-service" ./cmd/auth-service || error "Failed to build auth-service"
    
    chmod +x "$CMP_HOME/bin"/*
    chown -R "$CMP_USER:$CMP_GROUP" "$CMP_HOME/bin"
    
    # Copy backend code
    cp -r . "$CMP_HOME/backend/" || error "Failed to copy backend code"
    
    log "Backend services built successfully"
}

# Build frontend
build_frontend() {
    log "Building frontend application..."
    
    cd "$(dirname "$0")/../.." || error "Cannot find project root"
    cd frontend/webapp || error "Frontend directory not found"
    
    # Install dependencies
    npm ci --production || error "Failed to install frontend dependencies"
    
    # Build for production
    npm run build || error "Failed to build frontend"
    
    # Copy built files
    mkdir -p "$CMP_HOME/frontend/dist"
    cp -r dist/* "$CMP_HOME/frontend/dist/" || error "Failed to copy frontend build"
    
    chown -R "$CMP_USER:$CMP_GROUP" "$CMP_HOME/frontend"
    
    log "Frontend built successfully"
}

# Install systemd services
install_systemd_services() {
    log "Installing systemd services..."
    
    cd "$(dirname "$0")" || error "Cannot find deploy directory"
    
    cp systemd/*.service /etc/systemd/system/ || error "Failed to copy systemd services"
    systemctl daemon-reload
    
    log "Systemd services installed"
}

# Setup Nginx
setup_nginx() {
    log "Setting up Nginx reverse proxy..."
    
    cd "$(dirname "$0")" || error "Cannot find deploy directory"
    
    # Copy nginx configuration
    cp -r nginx/* "$NGINX_CONFIG_DIR/" || error "Failed to copy nginx config"
    
    # Create nginx PID directory
    mkdir -p /var/run/nginx
    chown nginx:nginx /var/run/nginx
    
    # Test nginx configuration
    nginx -t -c "$NGINX_CONFIG_DIR/cmp-nginx.conf" || warning "Nginx config test failed, check configuration"
    
    log "Nginx configured"
}

# Setup database
setup_database() {
    log "Setting up database..."
    
    # Create database and user (if not exists)
    sudo -u postgres psql <<EOF || warning "Database may already exist"
CREATE USER cmp_user WITH PASSWORD 'CHANGE_THIS_PASSWORD_IN_PROD';
CREATE DATABASE cmp_db OWNER cmp_user;
GRANT ALL PRIVILEGES ON DATABASE cmp_db TO cmp_user;
EOF
    
    log "Database created (please change password in production!)"
    info "Run migrations: make migrate-up"
}

# Setup configuration files
setup_configuration() {
    log "Setting up configuration files..."
    
    cd "$(dirname "$0")" || error "Cannot find deploy directory"
    
    # Copy environment templates
    cp production/*.env.example "$CONFIG_DIR/" || warning "No env.example files found"
    
    log "Configuration files ready"
    warning "IMPORTANT: Edit configuration files in $CONFIG_DIR before starting services"
}

# Setup log rotation
setup_logrotate() {
    log "Setting up log rotation..."
    
    cat > /etc/logrotate.d/cmp <<EOF
$LOG_DIR/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 $CMP_USER $CMP_GROUP
    sharedscripts
    postrotate
        systemctl reload cmp-* 2>/dev/null || true
    endscript
}
EOF
    
    log "Log rotation configured"
}

# Setup firewall
setup_firewall() {
    log "Configuring firewall..."
    
    if command -v ufw &> /dev/null; then
        ufw allow 22/tcp comment 'SSH'
        ufw allow 80/tcp comment 'HTTP'
        ufw allow 443/tcp comment 'HTTPS'
        ufw --force enable || warning "Firewall enable failed"
    elif command -v firewall-cmd &> /dev/null; then
        systemctl enable firewalld
        systemctl start firewalld
        firewall-cmd --permanent --add-service=ssh
        firewall-cmd --permanent --add-service=http
        firewall-cmd --permanent --add-service=https
        firewall-cmd --reload
    fi
    
    log "Firewall configured"
}

# Main installation
main() {
    log "=========================================="
    log "CMP Platform Production Installation"
    log "Enterprise-Grade Deployment"
    log "=========================================="
    
    check_root
    detect_os
    
    install_dependencies
    create_user_and_directories
    setup_postgresql
    setup_redis
    setup_vault
    build_backend
    build_frontend
    install_systemd_services
    setup_nginx
    setup_database
    setup_configuration
    setup_logrotate
    setup_firewall
    
    log "=========================================="
    log "Installation completed successfully!"
    log "=========================================="
    info ""
    info "Next steps:"
    info "1. Edit configuration files in $CONFIG_DIR"
    info "2. Run database migrations: make migrate-up"
    info "3. Initialize Vault: ./deploy/vault-init.sh"
    info "4. Start services: systemctl start cmp-inventory cmp-issuer cmp-adapter cmp-auth"
    info "5. Start Nginx: systemctl start cmp-nginx"
    info ""
    warning "Remember to change all default passwords!"
}

# Run main function
main "$@"

