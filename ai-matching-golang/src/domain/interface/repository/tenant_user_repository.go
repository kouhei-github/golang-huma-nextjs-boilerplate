package repository

import (
	"ai-matching/db/sqlc"
	"context"
)

type TenantUserRepository interface {
	// Add user to tenant with role
	AddUserToTenant(ctx context.Context, params db.AddUserToTenantParams) (db.TenantUser, error)
	
	// Remove user from tenant
	RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error
	
	// Get all users in a tenant
	GetUsersByTenant(ctx context.Context, tenantID int64) ([]db.User, error)
	
	// Get all tenants for a user
	GetTenantsByUser(ctx context.Context, userID int64) ([]db.Tenant, error)
	
	// Get specific tenant-user relationship
	GetTenantUser(ctx context.Context, tenantID, userID int64) (db.TenantUser, error)
	
	// Update user role in tenant
	UpdateUserRoleInTenant(ctx context.Context, tenantID, userID int64, role string) (db.TenantUser, error)
	
	// List all users in tenant with details
	ListTenantUsers(ctx context.Context, tenantID int64) ([]db.ListTenantUsersRow, error)
}