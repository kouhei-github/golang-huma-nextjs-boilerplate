-- Drop indexes
DROP INDEX IF EXISTS idx_tenants_subdomain;
DROP INDEX IF EXISTS idx_tenants_organization_id;
DROP INDEX IF EXISTS idx_tenant_users_tenant_user;
DROP INDEX IF EXISTS idx_tenant_users_user_id;
DROP INDEX IF EXISTS idx_tenant_users_tenant_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_cognito_id;

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS tenant_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS organizations;