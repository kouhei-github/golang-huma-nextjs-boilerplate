package repository

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
)

type authRepository struct {
	queries db.Querier
}

func NewAuthRepository(queries db.Querier) repository.AuthRepository {
	return &authRepository{
		queries: queries,
	}
}

func (r *authRepository) CreateUserAuth(ctx context.Context, params db.CreateUserAuthParams) (db.UserAuth, error) {
	return r.queries.CreateUserAuth(ctx, params)
}

func (r *authRepository) GetUserAuthByToken(ctx context.Context, refreshToken string) (db.UserAuth, error) {
	return r.queries.GetUserAuthByToken(ctx, sql.NullString{String: refreshToken, Valid: true})
}

func (r *authRepository) GetUserAuthByUserID(ctx context.Context, userID int64) (db.UserAuth, error) {
	return r.queries.GetUserAuthByUserID(ctx, userID)
}

func (r *authRepository) UpdateUserAuth(ctx context.Context, params db.UpdateUserAuthParams) (db.UserAuth, error) {
	return r.queries.UpdateUserAuth(ctx, params)
}

func (r *authRepository) DeleteUserAuth(ctx context.Context, id int64) error {
	return r.queries.DeleteUserAuth(ctx, id)
}

func (r *authRepository) DeleteExpiredUserAuth(ctx context.Context) error {
	return r.queries.DeleteExpiredUserAuth(ctx)
}
