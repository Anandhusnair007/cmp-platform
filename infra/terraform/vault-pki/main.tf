terraform {
  required_version = ">= 1.0"
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 3.0"
    }
  }
}

provider "vault" {
  address = var.vault_address
  token   = var.vault_token
}

# Enable PKI secrets engine
resource "vault_mount" "pki" {
  path        = var.pki_path
  type        = "pki"
  description = "CMP PKI for certificate issuance"
  
  default_lease_ttl_seconds = 31536000  # 1 year
  max_lease_ttl_seconds     = 31536000  # 1 year
}

# Generate root CA certificate
resource "vault_pki_secret_backend_root_cert" "root" {
  backend = vault_mount.pki.path

  type    = "internal"
  common_name = var.ca_common_name
  ttl     = "87600h"  # 10 years
  
  key_bits = 4096
}

# Configure URLs
resource "vault_pki_secret_backend_config_urls" "urls" {
  backend                 = vault_mount.pki.path
  issuing_certificates    = ["${var.vault_address}/v1/${vault_mount.pki.path}/ca"]
  crl_distribution_points = ["${var.vault_address}/v1/${vault_mount.pki.path}/crl"]
}

# Create role for certificate signing
resource "vault_pki_secret_backend_role" "cmp_role" {
  backend = vault_mount.pki.path
  name    = var.role_name

  allowed_domains    = var.allowed_domains
  allow_subdomains   = true
  allow_bare_domains = true
  allow_glob_domains = false

  max_ttl         = "8760h"  # 1 year
  ttl             = "8760h"
  
  key_bits        = 2048
  key_type        = "rsa"
  allow_any_name  = false
  enforce_hostnames = true
}

# Output CA certificate
output "ca_certificate" {
  value     = vault_pki_secret_backend_root_cert.root.certificate
  sensitive = true
}

output "pki_path" {
  value = vault_mount.pki.path
}

output "role_name" {
  value = vault_pki_secret_backend_role.cmp_role.name
}
