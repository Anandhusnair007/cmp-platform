#!/bin/bash
set -euo pipefail

# Enterprise Monitoring Setup Script for CMP Platform
# Sets up Prometheus, Grafana, and alerting

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

detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
    else
        error "Cannot detect OS"
    fi
}

# Install Prometheus
install_prometheus() {
    log "Installing Prometheus..."
    
    PROM_VERSION="2.48.0"
    PROM_DIR="/opt/prometheus"
    PROM_USER="prometheus"
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get install -y -qq prometheus prometheus-node-exporter
    else
        # Install from binary
        mkdir -p "$PROM_DIR"
        cd /tmp
        wget -q "https://github.com/prometheus/prometheus/releases/download/v${PROM_VERSION}/prometheus-${PROM_VERSION}.linux-amd64.tar.gz"
        tar -xzf "prometheus-${PROM_VERSION}.linux-amd64.tar.gz"
        mv prometheus-${PROM_VERSION}.linux-amd64/* "$PROM_DIR/"
        rm -rf prometheus-${PROM_VERSION}.linux-amd64*
        
        useradd -r -s /bin/false "$PROM_USER" || true
        chown -R "$PROM_USER:$PROM_USER" "$PROM_DIR"
        
        # Create systemd service
        cat > /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus Monitoring System
After=network.target

[Service]
Type=simple
User=$PROM_USER
Group=$PROM_USER
WorkingDirectory=$PROM_DIR
ExecStart=$PROM_DIR/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/var/lib/prometheus
Restart=always

[Install]
WantedBy=multi-user.target
EOF
    fi
    
    # Create Prometheus configuration
    mkdir -p /etc/prometheus
    cat > /etc/prometheus/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'cmp-production'
    environment: 'production'

rule_files:
  - '/etc/prometheus/cmp_rules.yml'

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'cmp-services'
    static_configs:
      - targets:
        - 'localhost:9091'  # inventory-service
        - 'localhost:9092'  # issuer-service
        - 'localhost:9093'  # adapter-service
        - 'localhost:9094'  # auth-service
    metrics_path: '/metrics'

  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
EOF

    # Create alert rules
    cat > /etc/prometheus/cmp_rules.yml <<EOF
groups:
  - name: cmp_alerts
    interval: 30s
    rules:
      - alert: ServiceDown
        expr: up{job="cmp-services"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "CMP Service is down"
          description: "{{ \$labels.job }} has been down for more than 1 minute"

      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is above 90%"

      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage"
          description: "CPU usage is above 80%"
EOF

    chmod 644 /etc/prometheus/*.yml
    systemctl daemon-reload
    systemctl enable prometheus
    systemctl start prometheus
    
    log "Prometheus installed and configured"
}

# Install Grafana
install_grafana() {
    log "Installing Grafana..."
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get install -y -qq software-properties-common
        add-apt-repository -y "deb https://packages.grafana.com/oss/deb stable main"
        wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -
        apt-get update -qq
        apt-get install -y -qq grafana
    elif [ "$OS" == "rhel" ] || [ "$OS" == "centos" ]; then
        cat > /etc/yum.repos.d/grafana.repo <<EOF
[grafana]
name=grafana
baseurl=https://packages.grafana.com/oss/rpm
repo_gpgcheck=1
enabled=1
gpgcheck=1
gpgkey=https://packages.grafana.com/gpg.key
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
EOF
        yum install -y -q grafana
    fi
    
    # Configure Grafana
    sed -i 's/;admin_user = admin/admin_user = admin/' /etc/grafana/grafana.ini
    sed -i 's/;admin_password = admin/admin_password = CHANGE_THIS_PASSWORD/' /etc/grafana/grafana.ini
    
    # Setup Prometheus datasource
    mkdir -p /etc/grafana/provisioning/datasources
    cat > /etc/grafana/provisioning/datasources/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://localhost:9090
    isDefault: true
    editable: true
EOF

    systemctl enable grafana-server
    systemctl start grafana-server
    
    log "Grafana installed and configured"
    warning "Change Grafana admin password: http://localhost:3000 (admin/admin)"
}

# Install Node Exporter
install_node_exporter() {
    log "Installing Node Exporter..."
    
    if [ "$OS" == "ubuntu" ] || [ "$OS" == "debian" ]; then
        apt-get install -y -qq prometheus-node-exporter
        systemctl enable prometheus-node-exporter
        systemctl start prometheus-node-exporter
    else
        NODE_VERSION="1.7.0"
        cd /tmp
        wget -q "https://github.com/prometheus/node_exporter/releases/download/v${NODE_VERSION}/node_exporter-${NODE_VERSION}.linux-amd64.tar.gz"
        tar -xzf "node_exporter-${NODE_VERSION}.linux-amd64.tar.gz"
        mv node_exporter-${NODE_VERSION}.linux-amd64/node_exporter /usr/local/bin/
        
        cat > /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Prometheus Node Exporter
After=network.target

[Service]
Type=simple
User=nobody
ExecStart=/usr/local/bin/node_exporter
Restart=always

[Install]
WantedBy=multi-user.target
EOF
        
        systemctl daemon-reload
        systemctl enable node_exporter
        systemctl start node_exporter
    fi
    
    log "Node Exporter installed"
}

# Setup log aggregation
setup_logging() {
    log "Setting up centralized logging..."
    
    # Configure rsyslog for CMP
    cat > /etc/rsyslog.d/cmp.conf <<EOF
# CMP Platform Logging
:programname, isequal, "cmp-inventory" /var/log/cmp/inventory.log
:programname, isequal, "cmp-issuer" /var/log/cmp/issuer.log
:programname, isequal, "cmp-adapter" /var/log/cmp/adapter.log
:programname, isequal, "cmp-auth" /var/log/cmp/auth.log
& stop
EOF

    systemctl restart rsyslog
    
    # Setup log rotation
    cat > /etc/logrotate.d/cmp-monitoring <<EOF
/var/log/cmp/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 cmp cmp
    sharedscripts
    postrotate
        systemctl reload rsyslog > /dev/null 2>&1 || true
    endscript
}
EOF
    
    log "Logging configured"
}

# Main function
main() {
    log "=========================================="
    log "CMP Platform Monitoring Setup"
    log "Enterprise-Grade Configuration"
    log "=========================================="
    
    check_root
    detect_os
    
    install_prometheus
    install_node_exporter
    install_grafana
    setup_logging
    
    log "=========================================="
    log "Monitoring setup completed!"
    log "=========================================="
    info ""
    info "Access Points:"
    info "  - Prometheus: http://localhost:9090"
    info "  - Grafana: http://localhost:3000 (admin/admin)"
    info ""
    warning "Change Grafana admin password immediately!"
    warning "Configure firewall rules to restrict access"
}

main "$@"

