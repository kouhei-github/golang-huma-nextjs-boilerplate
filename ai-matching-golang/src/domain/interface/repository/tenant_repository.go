package repository

import (
	"ai-matching/db/sqlc"
	"context"
)

type TenantRepository interface {
	GetTenant(ctx context.Context, id int64) (db.Tenant, error)
	GetTenantBySubdomain(ctx context.Context, subdomain string) (db.Tenant, error)
	ListTenantsByOrganization(ctx context.Context, organizationID int64, limit, offset int32) ([]db.Tenant, error)
	CreateTenant(ctx context.Context, params db.CreateTenantParams) (db.Tenant, error)
	UpdateTenant(ctx context.Context, params db.UpdateTenantParams) (db.Tenant, error)
	DeleteTenant(ctx context.Context, id int64) error
}
