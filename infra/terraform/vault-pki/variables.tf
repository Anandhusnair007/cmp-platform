variable "vault_address" {
  description = "Vault server address"
  type        = string
  default     = "http://localhost:8200"
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "pki_path" {
  description = "Path where PKI secrets engine will be mounted"
  type        = string
  default     = "cmp-pki"
}

variable "ca_common_name" {
  description = "Common name for the root CA"
  type        = string
  default     = "CMP Root CA"
}

variable "role_name" {
  description = "Name of the role for certificate signing"
  type        = string
  default     = "cmp-role"
}

variable "allowed_domains" {
  description = "List of allowed domains for certificate signing"
  type        = list(string)
  default     = ["staging.example.com", "example.com"]
}
