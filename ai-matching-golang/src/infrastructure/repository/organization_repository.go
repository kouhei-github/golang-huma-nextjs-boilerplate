package repository

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
	"github.com/google/uuid"
)

type organizationRepository struct {
	queries db.Querier
}

func NewOrganizationRepository(queries db.Querier) repository.OrganizationRepository {
	return &organizationRepository{
		queries: queries,
	}
}

func (r *organizationRepository) GetOrganization(ctx context.Context, id uuid.UUID) (db.Organization, error) {
	return r.queries.GetOrganization(ctx, id)
}

func (r *organizationRepository) ListOrganizations(ctx context.Context, limit, offset int32) ([]db.Organization, error) {
	params := db.ListOrganizationsParams{
		Limit:  limit,
		Offset: offset,
	}
	return r.queries.ListOrganizations(ctx, params)
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, params db.CreateOrganizationParams) (db.Organization, error) {
	return r.queries.CreateOrganization(ctx, params)
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, params db.UpdateOrganizationParams) (db.Organization, error) {
	return r.queries.UpdateOrganization(ctx, params)
}

func (r *organizationRepository) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteOrganization(ctx, id)
}

// Relationship methods

func (r *organizationRepository) GetOrganizationWithTenants(ctx context.Context, id uuid.UUID) (db.GetOrganizationWithTenantsRow, error) {
	return r.queries.GetOrganizationWithTenants(ctx, id)
}

func (r *organizationRepository) GetTenantsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]db.Tenant, error) {
	return r.queries.GetTenantsByOrganization(ctx, organizationID)
}

func (r *organizationRepository) GetOrganizationByTenant(ctx context.Context, tenantID uuid.UUID) (db.Organization, error) {
	return r.queries.GetOrganizationByTenant(ctx, tenantID)
}

// Count methods

func (r *organizationRepository) CountOrganizations(ctx context.Context) (int64, error) {
	return r.queries.CountOrganizations(ctx)
}
