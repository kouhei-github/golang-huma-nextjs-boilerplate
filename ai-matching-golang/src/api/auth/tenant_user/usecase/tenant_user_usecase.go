package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TenantUserUsecase struct {
	tenantUserRepo repository.TenantUserRepository
	tenantRepo     repository.TenantRepository
	userRepo       repository.UserRepository
}

func NewTenantUserUsecase(tenantUserRepo repository.TenantUserRepository, tenantRepo repository.TenantRepository, userRepo repository.UserRepository) *TenantUserUsecase {
	return &TenantUserUsecase{
		tenantUserRepo: tenantUserRepo,
		tenantRepo:     tenantRepo,
		userRepo:       userRepo,
	}
}

// AddUserToTenant adds a user to a tenant with a specified role
func (u *TenantUserUsecase) AddUserToTenant(ctx context.Context, tenantID, userID int64, role string) error {
	// Verify tenant exists
	_, err := u.tenantRepo.GetTenant(ctx, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("tenant not found")
		}
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Verify user exists
	_, err = u.userRepo.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user already belongs to tenant
	exists, err := u.tenantRepo.CheckUserBelongsToTenant(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to check user tenant membership: %w", err)
	}
	if exists {
		return errors.New("user already belongs to this tenant")
	}

	// Add user to tenant
	_, err = u.tenantUserRepo.AddUserToTenant(ctx, db.AddUserToTenantParams{
		TenantID: tenantID,
		UserID:   userID,
		Role:     sql.NullString{String: role, Valid: role != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}

	return nil
}

// RemoveUserFromTenant removes a user from a tenant
func (u *TenantUserUsecase) RemoveUserFromTenant(ctx context.Context, tenantID, userID int64) error {
	// Check if user belongs to tenant
	exists, err := u.tenantRepo.CheckUserBelongsToTenant(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to check user tenant membership: %w", err)
	}
	if !exists {
		return errors.New("user does not belong to this tenant")
	}

	err = u.tenantUserRepo.RemoveUserFromTenant(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from tenant: %w", err)
	}

	return nil
}

// GetUsersByTenant gets all users in a tenant
func (u *TenantUserUsecase) GetUsersByTenant(ctx context.Context, tenantID int64) ([]db.User, error) {
	// Verify tenant exists
	_, err := u.tenantRepo.GetTenant(ctx, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	users, err := u.tenantUserRepo.GetUsersByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by tenant: %w", err)
	}

	return users, nil
}

// GetTenantsByUser gets all tenants for a user
func (u *TenantUserUsecase) GetTenantsByUser(ctx context.Context, userID int64) ([]db.Tenant, error) {
	// Verify user exists
	_, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	tenants, err := u.tenantUserRepo.GetTenantsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants by user: %w", err)
	}

	return tenants, nil
}

// UpdateUserRoleInTenant updates a user's role in a tenant
func (u *TenantUserUsecase) UpdateUserRoleInTenant(ctx context.Context, tenantID, userID int64, role string) error {
	// Check if user belongs to tenant
	exists, err := u.tenantRepo.CheckUserBelongsToTenant(ctx, tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to check user tenant membership: %w", err)
	}
	if !exists {
		return errors.New("user does not belong to this tenant")
	}

	_, err = u.tenantUserRepo.UpdateUserRoleInTenant(ctx, tenantID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// ListTenantUsers lists all users in a tenant with their details
func (u *TenantUserUsecase) ListTenantUsers(ctx context.Context, tenantID int64) ([]db.ListTenantUsersRow, error) {
	// Verify tenant exists
	_, err := u.tenantRepo.GetTenant(ctx, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	users, err := u.tenantUserRepo.ListTenantUsers(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenant users: %w", err)
	}

	return users, nil
}

// GetUsersNotInTenant gets users that are not in a specific tenant
func (u *TenantUserUsecase) GetUsersNotInTenant(ctx context.Context, tenantID int64, limit, offset int32) ([]db.User, error) {
	// Verify tenant exists
	_, err := u.tenantRepo.GetTenant(ctx, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	users, err := u.userRepo.GetUsersNotInTenant(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users not in tenant: %w", err)
	}

	return users, nil
}