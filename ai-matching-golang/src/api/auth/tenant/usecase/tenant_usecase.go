package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/tenant/requests"
	"ai-matching/src/api/auth/tenant/response"
	"ai-matching/src/domain/interface/repository"
	"context"

	"github.com/google/uuid"
)

type TenantUsecase struct {
	tenantRepo repository.TenantRepository
}

func NewTenantUsecase(tenantRepo repository.TenantRepository) *TenantUsecase {
	return &TenantUsecase{
		tenantRepo: tenantRepo,
	}
}

func (u *TenantUsecase) GetTenant(ctx context.Context, id uuid.UUID) (*response.TenantResponse, error) {
	tenant, err := u.tenantRepo.GetTenant(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.TenantResponse{
		ID:             tenant.ID,
		OrganizationID: tenant.OrganizationID,
		Name:           tenant.Name,
		Subdomain:      tenant.Subdomain,
		IsActive:       tenant.IsActive,
		CreatedAt:      tenant.CreatedAt,
		UpdatedAt:      tenant.UpdatedAt,
	}, nil
}

func (u *TenantUsecase) GetTenantBySubdomain(ctx context.Context, subdomain string) (*response.TenantResponse, error) {
	tenant, err := u.tenantRepo.GetTenantBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, err
	}

	return &response.TenantResponse{
		ID:             tenant.ID,
		OrganizationID: tenant.OrganizationID,
		Name:           tenant.Name,
		Subdomain:      tenant.Subdomain,
		IsActive:       tenant.IsActive,
		CreatedAt:      tenant.CreatedAt,
		UpdatedAt:      tenant.UpdatedAt,
	}, nil
}

func (u *TenantUsecase) ListTenantsByOrganization(ctx context.Context, organizationID uuid.UUID, page, pageSize int) (*response.TenantListResponse, error) {
	offset := (page - 1) * pageSize
	tenants, err := u.tenantRepo.ListTenantsByOrganization(ctx, organizationID, int32(pageSize), int32(offset))
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := u.tenantRepo.CountTenantsByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	tenantResponses := make([]response.TenantResponse, len(tenants))
	for i, tenant := range tenants {
		tenantResponses[i] = response.TenantResponse{
			ID:             tenant.ID,
			OrganizationID: tenant.OrganizationID,
			Name:           tenant.Name,
			Subdomain:      tenant.Subdomain,
			IsActive:       tenant.IsActive,
			CreatedAt:      tenant.CreatedAt,
			UpdatedAt:      tenant.UpdatedAt,
		}
	}

	return &response.TenantListResponse{
		Tenants:  tenantResponses,
		Total:    int(totalCount),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (u *TenantUsecase) CreateTenant(ctx context.Context, req requests.CreateTenantRequest) (*response.TenantResponse, error) {
	tenant, err := u.tenantRepo.CreateTenant(ctx, db.CreateTenantParams{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Subdomain:      req.Subdomain,
		IsActive:       req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &response.TenantResponse{
		ID:             tenant.ID,
		OrganizationID: tenant.OrganizationID,
		Name:           tenant.Name,
		Subdomain:      tenant.Subdomain,
		IsActive:       tenant.IsActive,
		CreatedAt:      tenant.CreatedAt,
		UpdatedAt:      tenant.UpdatedAt,
	}, nil
}

func (u *TenantUsecase) UpdateTenant(ctx context.Context, id uuid.UUID, req requests.UpdateTenantRequest) (*response.TenantResponse, error) {
	tenant, err := u.tenantRepo.UpdateTenant(ctx, db.UpdateTenantParams{
		ID:        id,
		Name:      req.Name,
		Subdomain: req.Subdomain,
		IsActive:  req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &response.TenantResponse{
		ID:             tenant.ID,
		OrganizationID: tenant.OrganizationID,
		Name:           tenant.Name,
		Subdomain:      tenant.Subdomain,
		IsActive:       tenant.IsActive,
		CreatedAt:      tenant.CreatedAt,
		UpdatedAt:      tenant.UpdatedAt,
	}, nil
}

func (u *TenantUsecase) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return u.tenantRepo.DeleteTenant(ctx, id)
}
