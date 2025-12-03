-- Add approval_workflows table
CREATE TABLE IF NOT EXISTS approval_workflows (
    id VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255) UNIQUE NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    requester_id VARCHAR(255) NOT NULL,
    approvers JSONB NOT NULL,
    required_approvals INTEGER NOT NULL DEFAULT 1,
    current_step INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approval_workflows_entity ON approval_workflows(entity_type, entity_id);
CREATE INDEX idx_approval_workflows_status ON approval_workflows(status);
CREATE INDEX idx_approval_workflows_requester ON approval_workflows(requester_id);

-- Add tenant support
ALTER TABLE owners ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255);
ALTER TABLE certificates ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255);
ALTER TABLE issuance_requests ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255);
ALTER TABLE agents ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_owners_tenant ON owners(tenant_id);
CREATE INDEX IF NOT EXISTS idx_certificates_tenant ON certificates(tenant_id);
CREATE INDEX IF NOT EXISTS idx_issuance_requests_tenant ON issuance_requests(tenant_id);
CREATE INDEX IF NOT EXISTS idx_agents_tenant ON agents(tenant_id);

-- Add tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    enabled BOOLEAN NOT NULL DEFAULT true,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add encryption keys table for encryption at rest
CREATE TABLE IF NOT EXISTS encryption_keys (
    id VARCHAR(255) PRIMARY KEY,
    key_type VARCHAR(50) NOT NULL,
    key_data BYTEA NOT NULL,
    key_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rotated_at TIMESTAMPTZ,
    active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_encryption_keys_active ON encryption_keys(active, key_type);
