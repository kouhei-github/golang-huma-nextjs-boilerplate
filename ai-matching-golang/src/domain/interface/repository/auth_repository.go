package repository

import (
	"ai-matching/db/sqlc"
	"context"
)

type AuthRepository interface {
	CreateUserAuth(ctx context.Context, params db.CreateUserAuthParams) (db.UserAuth, error)
	GetUserAuthByToken(ctx context.Context, refreshToken string) (db.UserAuth, error)
	GetUserAuthByUserID(ctx context.Context, userID int64) (db.UserAuth, error)
	UpdateUserAuth(ctx context.Context, params db.UpdateUserAuthParams) (db.UserAuth, error)
	DeleteUserAuth(ctx context.Context, id int64) error
	DeleteExpiredUserAuth(ctx context.Context) error
}
