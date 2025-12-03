#!/bin/bash
set -e

echo "=== CMP End-to-End Test ==="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8082}"
AGENT_ID="${AGENT_ID:-agent-1}"
NGINX_URL="${NGINX_URL:-https://localhost:8443}"

echo "API URL: $API_URL"
echo "Agent ID: $AGENT_ID"
echo "Nginx URL: $NGINX_URL"

# Step 1: Request a certificate
echo -e "\n${GREEN}[1/4] Requesting certificate...${NC}"
REQUEST_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/certs/request" \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": "default-owner",
    "common_name": "app.staging.example.com",
    "sans": ["app.staging.example.com", "www.staging.example.com"],
    "key_algorithm": "rsa",
    "key_size": 2048,
    "adapter_id": "vault-staging",
    "install_targets": [{
      "agent_id": "'"$AGENT_ID"'",
      "path": "/etc/nginx/ssl/app.pem",
      "reload_cmd": "systemctl reload nginx || nginx -s reload"
    }]
  }')

REQUEST_ID=$(echo $REQUEST_RESPONSE | jq -r '.request_id')
if [ -z "$REQUEST_ID" ] || [ "$REQUEST_ID" = "null" ]; then
  echo -e "${RED}Failed to create certificate request${NC}"
  echo "Response: $REQUEST_RESPONSE"
  exit 1
fi

echo "Request ID: $REQUEST_ID"

# Step 2: Wait for certificate to be issued (poll)
echo -e "\n${GREEN}[2/4] Waiting for certificate issuance...${NC}"
MAX_WAIT=60
WAIT_COUNT=0
CERT_ID=""

while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
  CERT_RESPONSE=$(curl -s "$API_URL/api/v1/certs?limit=10")
  CERT_ID=$(echo $CERT_RESPONSE | jq -r '.certificates[0].id // empty')
  
  if [ -n "$CERT_ID" ] && [ "$CERT_ID" != "null" ]; then
    CERT_STATUS=$(echo $CERT_RESPONSE | jq -r '.certificates[0].status // empty')
    if [ "$CERT_STATUS" = "active" ]; then
      echo "Certificate issued: $CERT_ID"
      break
    fi
  fi
  
  echo "Waiting for certificate... ($WAIT_COUNT/$MAX_WAIT)"
  sleep 2
  WAIT_COUNT=$((WAIT_COUNT + 2))
done

if [ -z "$CERT_ID" ] || [ "$CERT_ID" = "null" ]; then
  echo -e "${RED}Certificate was not issued in time${NC}"
  exit 1
fi

# Step 3: Install certificate via agent
echo -e "\n${GREEN}[3/4] Installing certificate via agent...${NC}"
INSTALL_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/agents/$AGENT_ID/install" \
  -H "Content-Type: application/json" \
  -d '{
    "cert_id": "'"$CERT_ID"'",
    "path": "/etc/nginx/ssl/app.pem",
    "reload_cmd": "nginx -s reload || true"
  }')

INSTALL_STATUS=$(echo $INSTALL_RESPONSE | jq -r '.status // empty')
if [ "$INSTALL_STATUS" != "installed" ] && [ "$INSTALL_STATUS" != "pending" ]; then
  echo -e "${RED}Failed to install certificate${NC}"
  echo "Response: $INSTALL_RESPONSE"
  exit 1
fi

echo "Installation status: $INSTALL_STATUS"
sleep 5  # Wait for installation to complete

# Step 4: Verify HTTPS connection
echo -e "\n${GREEN}[4/4] Verifying HTTPS connection...${NC}"

# Test HTTPS connection (skip cert verification for self-signed in dev)
HTTP_CODE=$(curl -s -k -o /dev/null -w "%{http_code}" "$NGINX_URL" || echo "000")

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ HTTPS connection successful (HTTP $HTTP_CODE)${NC}"
  
  # Verify certificate
  CERT_INFO=$(echo | openssl s_client -connect localhost:8443 -servername app.staging.example.com 2>/dev/null | openssl x509 -noout -subject -dates 2>/dev/null || echo "")
  if [ -n "$CERT_INFO" ]; then
    echo -e "${GREEN}✓ Certificate details:${NC}"
    echo "$CERT_INFO"
  fi
else
  echo -e "${RED}✗ HTTPS connection failed (HTTP $HTTP_CODE)${NC}"
  exit 1
fi

echo -e "\n${GREEN}=== E2E Test Passed ===${NC}"
