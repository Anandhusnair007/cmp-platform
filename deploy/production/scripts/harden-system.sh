#!/bin/bash
set -euo pipefail

# Enterprise Security Hardening Script for CMP Platform
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

# Configure kernel parameters
harden_kernel() {
    log "Hardening kernel parameters..."
    
    cat >> /etc/sysctl.d/99-cmp-security.conf <<EOF
# CMP Platform Security Hardening

# Network security
net.ipv4.ip_forward = 0
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.default.send_redirects = 0
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0
net.ipv4.conf.all.secure_redirects = 0
net.ipv4.conf.default.secure_redirects = 0
net.ipv4.conf.all.log_martians = 1
net.ipv4.conf.default.log_martians = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.icmp_ignore_bogus_error_responses = 1
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.tcp_syncookies = 1

# IPv6 security
net.ipv6.conf.all.accept_redirects = 0
net.ipv6.conf.default.accept_redirects = 0
net.ipv6.conf.all.accept_ra = 0
net.ipv6.conf.default.accept_ra = 0

# Memory protection
kernel.exec-shield = 1
kernel.randomize_va_space = 2
vm.mmap_rnd_bits = 32
vm.mmap_rnd_compat_bits = 16

# Process limits
kernel.dmesg_restrict = 1
kernel.kptr_restrict = 2
kernel.unprivileged_bpf_disabled = 1
kernel.unprivileged_userns_clone = 0
EOF

    sysctl -p /etc/sysctl.d/99-cmp-security.conf
    log "Kernel parameters hardened"
}

# Configure system limits
configure_limits() {
    log "Configuring system limits..."
    
    cat >> /etc/security/limits.d/cmp.conf <<EOF
# CMP Platform Limits
cmp soft nofile 65536
cmp hard nofile 65536
cmp soft nproc 4096
cmp hard nproc 4096
cmp soft memlock unlimited
cmp hard memlock unlimited
EOF

    log "System limits configured"
}

# Configure SSH hardening
harden_ssh() {
    log "Hardening SSH configuration..."
    
    # Backup original config
    cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup
    
    # Apply security settings
    sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config
    
    cat >> /etc/ssh/sshd_config <<EOF

# CMP Platform SSH Hardening
Protocol 2
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
PermitEmptyPasswords no
MaxAuthTries 3
ClientAliveInterval 300
ClientAliveCountMax 2
X11Forwarding no
AllowTcpForwarding no
PermitTunnel no
AllowAgentForwarding no
EOF

    systemctl restart sshd || warning "SSH restart failed, check configuration"
    log "SSH hardened"
}

# Configure fail2ban
configure_fail2ban() {
    log "Configuring fail2ban..."
    
    cat > /etc/fail2ban/jail.d/cmp.conf <<EOF
[cmp-nginx]
enabled = true
port = http,https
filter = cmp-nginx
logpath = /var/log/nginx/cmp-access.log
maxretry = 5
bantime = 3600
findtime = 600

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 7200
findtime = 600
EOF

    cat > /etc/fail2ban/filter.d/cmp-nginx.conf <<EOF
[Definition]
failregex = ^<HOST>.*\"(GET|POST|PUT|DELETE).*\" (4[0-9][0-9]|5[0-9][0-9])
ignoreregex =
EOF

    systemctl enable fail2ban
    systemctl restart fail2ban
    log "fail2ban configured"
}

# Configure auditd
configure_auditd() {
    log "Configuring audit logging..."
    
    if command -v auditd &> /dev/null; then
        cat > /etc/audit/rules.d/cmp.rules <<EOF
# CMP Platform Audit Rules

# Watch CMP directories
-w /opt/cmp -p wa -k cmp_changes
-w /etc/cmp -p wa -k cmp_config_changes
-w /var/log/cmp -p wa -k cmp_log_changes

# Watch system binaries
-w /usr/bin/systemctl -p x -k systemd_control
-w /usr/bin/nginx -p x -k nginx_exec

# Audit system calls
-a always,exit -F arch=b64 -S execve -k exec
-a always,exit -F arch=b32 -S execve -k exec

# Audit network
-a always,exit -F arch=b64 -S socket -k network
-a always,exit -F arch=b64 -S bind -k network_bind
EOF

        systemctl enable auditd
        systemctl restart auditd
        log "Audit logging configured"
    else
        warning "auditd not installed, skipping"
    fi
}

# Generate secure passwords
generate_passwords() {
    log "Generating secure passwords..."
    
    # Generate random passwords
    DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    REDIS_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    JWT_SECRET=$(openssl rand -hex 32)
    
    cat > /etc/cmp/.passwords <<EOF
# Generated Secure Passwords
# Store these securely and update configuration files

DB_PASSWORD=$DB_PASSWORD
REDIS_PASSWORD=$REDIS_PASSWORD
JWT_SECRET=$JWT_SECRET

Generated: $(date)
EOF

    chmod 600 /etc/cmp/.passwords
    chown root:root /etc/cmp/.passwords
    
    log "Passwords generated and saved to /etc/cmp/.passwords"
    warning "IMPORTANT: Save these passwords securely and update configuration files"
}

# Configure SSL/TLS certificates
setup_ssl() {
    log "Setting up SSL/TLS certificates..."
    
    CERT_DIR="/etc/cmp/ssl"
    mkdir -p "$CERT_DIR"
    
    info "For production, use Let's Encrypt or your CA certificates"
    info "Run: certbot certonly --nginx -d api.cmp.example.com -d app.cmp.example.com"
    
    # Create self-signed cert for testing (replace in production)
    if [ ! -f "$CERT_DIR/fullchain.pem" ]; then
        openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
            -keyout "$CERT_DIR/privkey.pem" \
            -out "$CERT_DIR/fullchain.pem" \
            -subj "/C=US/ST=State/L=City/O=CMP/CN=api.cmp.example.com" || warning "Self-signed cert generation failed"
        
        cp "$CERT_DIR/fullchain.pem" "$CERT_DIR/chain.pem"
        chmod 600 "$CERT_DIR"/*.pem
        chown root:root "$CERT_DIR"/*.pem
    fi
    
    log "SSL certificates configured"
}

# Main function
main() {
    log "=========================================="
    log "CMP Platform Security Hardening"
    log "Enterprise-Grade Configuration"
    log "=========================================="
    
    check_root
    
    harden_kernel
    configure_limits
    harden_ssh
    configure_fail2ban
    configure_auditd
    generate_passwords
    setup_ssl
    
    log "=========================================="
    log "Security hardening completed!"
    log "=========================================="
    warning "Review and test all configurations before deploying to production"
}

main "$@"

