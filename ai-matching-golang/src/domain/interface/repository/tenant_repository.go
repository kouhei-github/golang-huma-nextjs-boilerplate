package repository

import (
	"ai-matching/db/sqlc"
	"context"
	
	"github.com/google/uuid"
)

type TenantRepository interface {
	GetTenant(ctx context.Context, id uuid.UUID) (db.Tenant, error)
	GetTenantBySubdomain(ctx context.Context, subdomain string) (db.Tenant, error)
	ListTenantsByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int32) ([]db.Tenant, error)
	CreateTenant(ctx context.Context, params db.CreateTenantParams) (db.Tenant, error)
	UpdateTenant(ctx context.Context, params db.UpdateTenantParams) (db.Tenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	
	// Relationship methods
	GetTenantWithUserCount(ctx context.Context, id uuid.UUID) (db.GetTenantWithUserCountRow, error)
	GetTenantsByUserID(ctx context.Context, userID uuid.UUID) ([]db.GetTenantsByUserIDRow, error)
	CheckUserBelongsToTenant(ctx context.Context, tenantID, userID uuid.UUID) (bool, error)
}
