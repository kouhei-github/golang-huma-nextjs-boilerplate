package repository

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type tenantUserRepository struct {
	queries db.Querier
}

func NewTenantUserRepository(queries db.Querier) repository.TenantUserRepository {
	return &tenantUserRepository{
		queries: queries,
	}
}

func (r *tenantUserRepository) AddUserToTenant(ctx context.Context, params db.AddUserToTenantParams) (db.TenantUser, error) {
	return r.queries.AddUserToTenant(ctx, params)
}

func (r *tenantUserRepository) RemoveUserFromTenant(ctx context.Context, tenantID, userID uuid.UUID) error {
	return r.queries.RemoveUserFromTenant(ctx, db.RemoveUserFromTenantParams{
		TenantID: tenantID,
		UserID:   userID,
	})
}

func (r *tenantUserRepository) GetUsersByTenant(ctx context.Context, tenantID uuid.UUID) ([]db.User, error) {
	return r.queries.GetUsersByTenant(ctx, tenantID)
}

func (r *tenantUserRepository) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]db.Tenant, error) {
	return r.queries.GetTenantsByUser(ctx, userID)
}

func (r *tenantUserRepository) GetTenantUser(ctx context.Context, tenantID, userID uuid.UUID) (db.TenantUser, error) {
	params := db.GetTenantUserParams{
		TenantID: tenantID,
		UserID:   userID,
	}
	return r.queries.GetTenantUser(ctx, params)
}

func (r *tenantUserRepository) UpdateUserRoleInTenant(ctx context.Context, tenantID, userID uuid.UUID, role string) (db.TenantUser, error) {
	params := db.UpdateUserRoleInTenantParams{
		TenantID: tenantID,
		UserID:   userID,
		Role:     sql.NullString{String: role, Valid: true},
	}
	return r.queries.UpdateUserRoleInTenant(ctx, params)
}

func (r *tenantUserRepository) ListTenantUsers(ctx context.Context, tenantID uuid.UUID) ([]db.ListTenantUsersRow, error) {
	return r.queries.ListTenantUsers(ctx, tenantID)
}