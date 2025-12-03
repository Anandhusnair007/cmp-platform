#!/bin/bash
set -e

echo "Initializing Vault PKI for CMP..."

export VAULT_ADDR="${VAULT_ADDR:-http://localhost:8200}"
export VAULT_TOKEN="${VAULT_TOKEN:-dev-only-token}"

# Wait for Vault to be ready
echo "Waiting for Vault to be ready..."
for i in {1..30}; do
  if vault status &>/dev/null; then
    break
  fi
  sleep 1
done

# Enable PKI if not already enabled
if ! vault secrets list | grep -q "^cmp-pki/"; then
  echo "Enabling PKI secrets engine..."
  vault secrets enable -path=cmp-pki pki
  
  # Tune PKI
  vault secrets tune -max-lease-ttl=87600h cmp-pki
  
  # Generate root CA
  echo "Generating root CA..."
  vault write -field=certificate cmp-pki/root/generate/internal \
    common_name="CMP Root CA" \
    ttl=87600h > /tmp/ca.crt
  
  # Configure URLs
  echo "Configuring PKI URLs..."
  vault write cmp-pki/config/urls \
    issuing_certificates="${VAULT_ADDR}/v1/cmp-pki/ca" \
    crl_distribution_points="${VAULT_ADDR}/v1/cmp-pki/crl"
  
  # Create role for certificate signing
  echo "Creating signing role..."
  vault write cmp-pki/roles/cmp-role \
    allowed_domains="staging.example.com,example.com" \
    allow_subdomains=true \
    allow_bare_domains=true \
    max_ttl="8760h" \
    ttl="8760h" \
    key_bits=2048 \
    key_type="rsa"
  
  echo "Vault PKI initialized successfully!"
  echo "CA Certificate saved to /tmp/ca.crt"
else
  echo "PKI already enabled, skipping initialization"
fi

# Register adapter in database (requires database connection)
echo "Note: Register adapter in database manually or via API"
echo "Adapter config:"
echo "  ID: vault-staging"
echo "  Type: vault-pki"
echo "  Config: {\"address\":\"$VAULT_ADDR\",\"token\":\"$VAULT_TOKEN\",\"mount_path\":\"cmp-pki\",\"role_name\":\"cmp-role\"}"
