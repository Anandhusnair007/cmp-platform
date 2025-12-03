DROP TABLE IF EXISTS approval_workflows;
DROP TABLE IF EXISTS encryption_keys;
DROP TABLE IF EXISTS tenants;

ALTER TABLE owners DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE certificates DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE issuance_requests DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE agents DROP COLUMN IF EXISTS tenant_id;
