package repository

import (
	"ai-matching/db/sqlc"
	"context"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByCognitoID(ctx context.Context, cognitoID string) (db.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]db.User, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	DeleteUser(ctx context.Context, id int64) error
	
	// Relationship methods
	GetUsersNotInTenant(ctx context.Context, tenantID int64, limit, offset int32) ([]db.User, error)
	GetUserWithTenants(ctx context.Context, id int64) (db.GetUserWithTenantsRow, error)
}
