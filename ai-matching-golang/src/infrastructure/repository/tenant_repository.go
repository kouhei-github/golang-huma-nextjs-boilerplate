package repository

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
)

type tenantRepository struct {
	queries db.Querier
}

func NewTenantRepository(queries db.Querier) repository.TenantRepository {
	return &tenantRepository{
		queries: queries,
	}
}

func (r *tenantRepository) GetTenant(ctx context.Context, id int64) (db.Tenant, error) {
	return r.queries.GetTenant(ctx, id)
}

func (r *tenantRepository) GetTenantBySubdomain(ctx context.Context, subdomain string) (db.Tenant, error) {
	return r.queries.GetTenantBySubdomain(ctx, subdomain)
}

func (r *tenantRepository) ListTenantsByOrganization(ctx context.Context, organizationID int64, limit, offset int32) ([]db.Tenant, error) {
	params := db.ListTenantsByOrganizationParams{
		OrganizationID: organizationID,
		Limit:          limit,
		Offset:         offset,
	}
	return r.queries.ListTenantsByOrganization(ctx, params)
}

func (r *tenantRepository) CreateTenant(ctx context.Context, params db.CreateTenantParams) (db.Tenant, error) {
	return r.queries.CreateTenant(ctx, params)
}

func (r *tenantRepository) UpdateTenant(ctx context.Context, params db.UpdateTenantParams) (db.Tenant, error) {
	return r.queries.UpdateTenant(ctx, params)
}

func (r *tenantRepository) DeleteTenant(ctx context.Context, id int64) error {
	return r.queries.DeleteTenant(ctx, id)
}
