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

// Relationship methods

func (r *tenantRepository) GetTenantWithUserCount(ctx context.Context, id int64) (db.GetTenantWithUserCountRow, error) {
	return r.queries.GetTenantWithUserCount(ctx, id)
}

func (r *tenantRepository) GetTenantsByUserID(ctx context.Context, userID int64) ([]db.GetTenantsByUserIDRow, error) {
	return r.queries.GetTenantsByUserID(ctx, userID)
}

func (r *tenantRepository) CheckUserBelongsToTenant(ctx context.Context, tenantID, userID int64) (bool, error) {
	belongs, err := r.queries.CheckUserBelongsToTenant(ctx, db.CheckUserBelongsToTenantParams{
		TenantID: tenantID,
		UserID:   userID,
	})
	if err != nil {
		return false, err
	}
	return belongs, nil
}
