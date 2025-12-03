# Linux Agent

The Linux agent deploys certificates to Linux hosts and manages certificate installation.

## Features

- Register with CMP
- Accept installation jobs
- Download and install certificates
- Execute reload commands
- Report status

## Building

```bash
make build
```

## Running

```bash
./bin/agent \
  -api-url=http://localhost:8082 \
  -api-token=your-token \
  -agent-id=agent-1 \
  -hostname=$(hostname) \
  -cert-dir=/var/lib/cmp-agent/certs
```

## Environment Variables

- `CMP_API_URL` - CMP API URL
- `CMP_API_TOKEN` - Authentication token
- `AGENT_ID` - Unique agent identifier
- `AGENT_HOSTNAME` - Hostname for this agent
- `CERT_DIR` - Directory to store certificates

## Installation

### Systemd Service

Create `/etc/systemd/system/cmp-agent.service`:

```ini
[Unit]
Description=CMP Linux Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/cmp-agent \
  -api-url=http://cmp-api:8082 \
  -cert-dir=/var/lib/cmp-agent/certs
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable cmp-agent
sudo systemctl start cmp-agent
```

## Security

- Agent authenticates via API token or Vault token
- Certificates stored securely on filesystem
- Minimal permissions required
