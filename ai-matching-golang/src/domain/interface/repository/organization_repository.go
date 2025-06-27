package repository

import (
	"ai-matching/db/sqlc"
	"context"
	
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, id uuid.UUID) (db.Organization, error)
	ListOrganizations(ctx context.Context, limit, offset int32) ([]db.Organization, error)
	CreateOrganization(ctx context.Context, params db.CreateOrganizationParams) (db.Organization, error)
	UpdateOrganization(ctx context.Context, params db.UpdateOrganizationParams) (db.Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	
	// Relationship methods
	GetOrganizationWithTenants(ctx context.Context, id uuid.UUID) (db.GetOrganizationWithTenantsRow, error)
	GetTenantsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]db.Tenant, error)
	GetOrganizationByTenant(ctx context.Context, tenantID uuid.UUID) (db.Organization, error)
}
