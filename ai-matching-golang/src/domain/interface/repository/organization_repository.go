package repository

import (
	"ai-matching/db/sqlc"
	"context"
)

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, id int64) (db.Organization, error)
	ListOrganizations(ctx context.Context, limit, offset int32) ([]db.Organization, error)
	CreateOrganization(ctx context.Context, params db.CreateOrganizationParams) (db.Organization, error)
	UpdateOrganization(ctx context.Context, params db.UpdateOrganizationParams) (db.Organization, error)
	DeleteOrganization(ctx context.Context, id int64) error
	
	// Relationship methods
	GetOrganizationWithTenants(ctx context.Context, id int64) (db.GetOrganizationWithTenantsRow, error)
	GetTenantsByOrganization(ctx context.Context, organizationID int64) ([]db.Tenant, error)
	GetOrganizationByTenant(ctx context.Context, tenantID int64) (db.Organization, error)
}
