package repository

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
	"github.com/google/uuid"
)

type userRepository struct {
	queries db.Querier
}

func NewUserRepository(queries db.Querier) repository.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) GetUser(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *userRepository) GetUserByCognitoID(ctx context.Context, cognitoID string) (db.User, error) {
	return r.queries.GetUserByCognitoID(ctx, cognitoID)
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int32) ([]db.User, error) {
	params := db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}
	return r.queries.ListUsers(ctx, params)
}

func (r *userRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *userRepository) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	return r.queries.UpdateUser(ctx, params)
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

// Relationship methods

func (r *userRepository) GetUsersNotInTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]db.User, error) {
	params := db.GetUsersNotInTenantParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
	}
	return r.queries.GetUsersNotInTenant(ctx, params)
}

func (r *userRepository) GetUserWithTenants(ctx context.Context, id uuid.UUID) (db.GetUserWithTenantsRow, error) {
	return r.queries.GetUserWithTenants(ctx, id)
}

// Count methods

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) CountUsersNotInTenant(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return r.queries.CountUsersNotInTenant(ctx, tenantID)
}
