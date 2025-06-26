DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_organization_id;
DROP INDEX IF EXISTS idx_tenants_subdomain;
DROP INDEX IF EXISTS idx_tenants_organization_id;

ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;

DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS organizations;