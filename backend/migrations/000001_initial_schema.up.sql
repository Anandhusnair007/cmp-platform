-- Create owners table
CREATE TABLE IF NOT EXISTS owners (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    team VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create adapters table
CREATE TABLE IF NOT EXISTS adapters (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id VARCHAR(255) PRIMARY KEY,
    fingerprint VARCHAR(255) UNIQUE NOT NULL,
    subject VARCHAR(255) NOT NULL,
    sans JSONB NOT NULL DEFAULT '[]',
    issuer VARCHAR(255) NOT NULL,
    not_before TIMESTAMPTZ NOT NULL,
    not_after TIMESTAMPTZ NOT NULL,
    key_algo VARCHAR(50) NOT NULL,
    key_size INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    owner_id VARCHAR(255) REFERENCES owners(id),
    last_scanned_at TIMESTAMPTZ,
    source VARCHAR(100),
    cert_pem TEXT,
    chain_pem TEXT,
    private_key_ref VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes on certificates
CREATE INDEX idx_certificates_owner_id ON certificates(owner_id);
CREATE INDEX idx_certificates_status ON certificates(status);
CREATE INDEX idx_certificates_not_after ON certificates(not_after);
CREATE INDEX idx_certificates_source ON certificates(source);

-- Create issuance_requests table
CREATE TABLE IF NOT EXISTS issuance_requests (
    id VARCHAR(255) PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL REFERENCES owners(id),
    csr TEXT,
    csr_attrs JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    adapter_id VARCHAR(255) NOT NULL REFERENCES adapters(id),
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    issued_cert_id VARCHAR(255) REFERENCES certificates(id),
    error_message TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes on issuance_requests
CREATE INDEX idx_issuance_requests_owner_id ON issuance_requests(owner_id);
CREATE INDEX idx_issuance_requests_status ON issuance_requests(status);
CREATE INDEX idx_issuance_requests_adapter_id ON issuance_requests(adapter_id);

-- Create agents table
CREATE TABLE IF NOT EXISTS agents (
    id VARCHAR(255) PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip VARCHAR(45),
    last_checkin TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    capabilities JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    auth_token_hash VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes on agents
CREATE INDEX idx_agents_hostname ON agents(hostname);
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_last_checkin ON agents(last_checkin);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    performed_by VARCHAR(255) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    details JSONB NOT NULL DEFAULT '{}',
    ip_address VARCHAR(45),
    user_agent TEXT
);

-- Create indexes on audit_logs
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_performed_by ON audit_logs(performed_by);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Create installation_jobs table
CREATE TABLE IF NOT EXISTS installation_jobs (
    id VARCHAR(255) PRIMARY KEY,
    agent_id VARCHAR(255) NOT NULL REFERENCES agents(id),
    cert_id VARCHAR(255) NOT NULL REFERENCES certificates(id),
    target_path VARCHAR(500) NOT NULL,
    reload_cmd VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes on installation_jobs
CREATE INDEX idx_installation_jobs_agent_id ON installation_jobs(agent_id);
CREATE INDEX idx_installation_jobs_cert_id ON installation_jobs(cert_id);
CREATE INDEX idx_installation_jobs_status ON installation_jobs(status);

-- Insert default owner for testing
INSERT INTO owners (id, name, email, team) VALUES 
    ('default-owner', 'Default Owner', 'admin@example.com', 'Platform')
ON CONFLICT (id) DO NOTHING;
