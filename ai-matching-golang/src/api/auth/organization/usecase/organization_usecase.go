package usecase

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/organization/requests"
	"ai-matching/src/api/auth/organization/response"
	"ai-matching/src/domain/interface/repository"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type OrganizationUsecase struct {
	orgRepo repository.OrganizationRepository
}

func NewOrganizationUsecase(orgRepo repository.OrganizationRepository) *OrganizationUsecase {
	return &OrganizationUsecase{
		orgRepo: orgRepo,
	}
}

func (u *OrganizationUsecase) GetOrganization(ctx context.Context, id uuid.UUID) (*response.OrganizationResponse, error) {
	org, err := u.orgRepo.GetOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description.String,
		IsActive:    org.IsActive,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

func (u *OrganizationUsecase) ListOrganizations(ctx context.Context, page, pageSize int) (*response.OrganizationListResponse, error) {
	offset := (page - 1) * pageSize
	orgs, err := u.orgRepo.ListOrganizations(ctx, int32(pageSize), int32(offset))
	if err != nil {
		return nil, err
	}

	// Get total count
	totalCount, err := u.orgRepo.CountOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	organizations := make([]response.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		organizations[i] = response.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description.String,
			IsActive:    org.IsActive,
			CreatedAt:   org.CreatedAt,
			UpdatedAt:   org.UpdatedAt,
		}
	}

	return &response.OrganizationListResponse{
		Organizations: organizations,
		Total:         int(totalCount),
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

func (u *OrganizationUsecase) CreateOrganization(ctx context.Context, req requests.CreateOrganizationRequest) (*response.OrganizationResponse, error) {
	org, err := u.orgRepo.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		IsActive:    req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &response.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description.String,
		IsActive:    org.IsActive,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

func (u *OrganizationUsecase) UpdateOrganization(ctx context.Context, id uuid.UUID, req requests.UpdateOrganizationRequest) (*response.OrganizationResponse, error) {
	org, err := u.orgRepo.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:          id,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		IsActive:    req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &response.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description.String,
		IsActive:    org.IsActive,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}, nil
}

func (u *OrganizationUsecase) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return u.orgRepo.DeleteOrganization(ctx, id)
}
