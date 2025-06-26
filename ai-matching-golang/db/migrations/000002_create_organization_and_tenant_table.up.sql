CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD COLUMN organization_id BIGINT REFERENCES organizations(id);
ALTER TABLE users ADD COLUMN tenant_id BIGINT REFERENCES tenants(id);

CREATE INDEX idx_tenants_organization_id ON tenants(organization_id);
CREATE INDEX idx_tenants_subdomain ON tenants(subdomain);
CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);