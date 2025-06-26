package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/user/requests"
	"ai-matching/src/api/auth/user/response"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) GetUser(ctx context.Context, id int64) (*response.UserResponse, error) {
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName.String,
		LastName:       user.LastName.String,
		OrganizationID: nullInt64ToPtr(user.OrganizationID),
		TenantID:       nullInt64ToPtr(user.TenantID),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func (u *UserUsecase) ListUsers(ctx context.Context, page, pageSize int) (*response.UserListResponse, error) {
	offset := (page - 1) * pageSize
	users, err := u.userRepo.ListUsers(ctx, int32(pageSize), int32(offset))
	if err != nil {
		return nil, err
	}

	userResponses := make([]response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = response.UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			FirstName:      user.FirstName.String,
			LastName:       user.LastName.String,
			OrganizationID: nullInt64ToPtr(user.OrganizationID),
			TenantID:       nullInt64ToPtr(user.TenantID),
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}
	}

	return &response.UserListResponse{
		Users:    userResponses,
		Total:    len(userResponses),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (u *UserUsecase) CreateUser(ctx context.Context, req requests.CreateUserRequest) (*response.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.CreateUser(ctx, db.CreateUserParams{
		Email:          req.Email,
		PasswordHash:   string(hashedPassword),
		FirstName:      sql.NullString{String: req.FirstName, Valid: true},
		LastName:       sql.NullString{String: req.LastName, Valid: true},
		OrganizationID: ptrToNullInt64(req.OrganizationID),
		TenantID:       ptrToNullInt64(req.TenantID),
	})
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName.String,
		LastName:       user.LastName.String,
		OrganizationID: nullInt64ToPtr(user.OrganizationID),
		TenantID:       nullInt64ToPtr(user.TenantID),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, req requests.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := u.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		ID:        id,
		Email:     req.Email,
		FirstName: sql.NullString{String: req.FirstName, Valid: true},
		LastName:  sql.NullString{String: req.LastName, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName.String,
		LastName:       user.LastName.String,
		OrganizationID: nullInt64ToPtr(user.OrganizationID),
		TenantID:       nullInt64ToPtr(user.TenantID),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.DeleteUser(ctx, id)
}

func nullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func ptrToNullInt64(p *int64) sql.NullInt64 {
	if p != nil {
		return sql.NullInt64{Int64: *p, Valid: true}
	}
	return sql.NullInt64{}
}
