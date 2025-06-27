package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/user/requests"
	"ai-matching/src/api/auth/user/response"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
	"fmt"
)

type UserUsecase struct {
	userRepo       repository.UserRepository
	tenantUserRepo repository.TenantUserRepository
}

func NewUserUsecase(userRepo repository.UserRepository, tenantUserRepo repository.TenantUserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		tenantUserRepo: tenantUserRepo,
	}
}

func (u *UserUsecase) GetUser(ctx context.Context, id int64) (*response.UserResponse, error) {
	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get user's tenants
	tenants, err := u.tenantUserRepo.GetTenantsByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	tenantInfos := make([]response.TenantInfo, len(tenants))
	for i, tenant := range tenants {
		// Get user's role in this tenant
		tenantUser, err := u.tenantUserRepo.GetTenantUser(ctx, tenant.ID, id)
		if err != nil {
			return nil, err
		}

		tenantInfos[i] = response.TenantInfo{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Subdomain: tenant.Subdomain,
			Role:      tenantUser.Role.String,
		}
	}

	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Tenants:   tenantInfos,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
		// Get user's tenants
		tenants, err := u.tenantUserRepo.GetTenantsByUser(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		tenantInfos := make([]response.TenantInfo, len(tenants))
		for j, tenant := range tenants {
			// Get user's role in this tenant
			tenantUser, err := u.tenantUserRepo.GetTenantUser(ctx, tenant.ID, user.ID)
			if err != nil {
				return nil, err
			}

			tenantInfos[j] = response.TenantInfo{
				ID:        tenant.ID,
				Name:      tenant.Name,
				Subdomain: tenant.Subdomain,
				Role:      tenantUser.Role.String,
			}
		}

		userResponses[i] = response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
			Tenants:   tenantInfos,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
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
	// Create user with cognito ID
	user, err := u.userRepo.CreateUser(ctx, db.CreateUserParams{
		CognitoID: req.CognitoID,
		Email:     req.Email,
		FirstName: sql.NullString{String: req.FirstName, Valid: true},
		LastName:  sql.NullString{String: req.LastName, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	var tenantInfos []response.TenantInfo

	// If a tenant ID is provided, associate the user with that tenant
	if req.TenantID != nil {
		role := "member" // Default role
		if req.TenantRole != nil {
			role = *req.TenantRole
		}

		_, err = u.tenantUserRepo.AddUserToTenant(ctx, db.AddUserToTenantParams{
			TenantID: *req.TenantID,
			UserID:   user.ID,
			Role:     sql.NullString{String: role, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to associate user with tenant: %w", err)
		}

		// Get tenant info for response
		tenant, err := u.tenantUserRepo.GetTenantsByUser(ctx, user.ID)
		if err == nil && len(tenant) > 0 {
			tenantInfos = []response.TenantInfo{
				{
					ID:        tenant[0].ID,
					Name:      tenant[0].Name,
					Subdomain: tenant[0].Subdomain,
					Role:      role,
				},
			}
		}
	}

	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Tenants:   tenantInfos,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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

	// Get user's tenants
	tenants, err := u.tenantUserRepo.GetTenantsByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	tenantInfos := make([]response.TenantInfo, len(tenants))
	for i, tenant := range tenants {
		// Get user's role in this tenant
		tenantUser, err := u.tenantUserRepo.GetTenantUser(ctx, tenant.ID, id)
		if err != nil {
			return nil, err
		}

		tenantInfos[i] = response.TenantInfo{
			ID:        tenant.ID,
			Name:      tenant.Name,
			Subdomain: tenant.Subdomain,
			Role:      tenantUser.Role.String,
		}
	}

	return &response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Tenants:   tenantInfos,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.DeleteUser(ctx, id)
}
