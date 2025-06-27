package repository

import (
	"ai-matching/db/sqlc"
	"context"
	
	"github.com/google/uuid"
)

type TenantUserRepository interface {
	// Add user to tenant with role
	AddUserToTenant(ctx context.Context, params db.AddUserToTenantParams) (db.TenantUser, error)
	
	// Remove user from tenant
	RemoveUserFromTenant(ctx context.Context, tenantID, userID uuid.UUID) error
	
	// Get all users in a tenant
	GetUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]db.User, error)
	
	// Get all tenants for a user
	GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]db.Tenant, error)
	
	// Get specific tenant-user relationship
	GetTenantUser(ctx context.Context, tenantID, userID uuid.UUID) (db.TenantUser, error)
	
	// Update user role in tenant
	UpdateUserRoleInTenant(ctx context.Context, tenantID, userID uuid.UUID, role string) (db.TenantUser, error)
	
	// List all users in tenant with details
	ListTenantUsers(ctx context.Context, tenantID uuid.UUID) ([]db.ListTenantUsersRow, error)
}