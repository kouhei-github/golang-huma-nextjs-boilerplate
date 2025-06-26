package controller

import (
	"ai-matching/db/sqlc"
	"ai-matching/src/api/auth/tenant/requests"
	"ai-matching/src/api/auth/tenant/response"
	"ai-matching/src/api/auth/tenant/usecase"
	"ai-matching/src/infrastructure/repository"
	"context"
)

type TenantController struct {
	usecase *usecase.TenantUsecase
}

func NewTenantController(queries db.Querier) *TenantController {
	tenantRepo := repository.NewTenantRepository(queries)
	tenantUsecase := usecase.NewTenantUsecase(tenantRepo)

	return &TenantController{
		usecase: tenantUsecase,
	}
}

type GetTenantInput struct {
	ID int64 `path:"id" doc:"Tenant ID"`
}

type GetTenantOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) GetTenant(ctx context.Context, input *GetTenantInput) (*GetTenantOutput, error) {
	resp, err := c.usecase.GetTenant(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetTenantOutput{Body: *resp}, nil
}

type GetTenantBySubdomainInput struct {
	Subdomain string `path:"subdomain" doc:"Tenant subdomain"`
}

type GetTenantBySubdomainOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) GetTenantBySubdomain(ctx context.Context, input *GetTenantBySubdomainInput) (*GetTenantBySubdomainOutput, error) {
	resp, err := c.usecase.GetTenantBySubdomain(ctx, input.Subdomain)
	if err != nil {
		return nil, err
	}

	return &GetTenantBySubdomainOutput{Body: *resp}, nil
}

type ListTenantsByOrganizationInput struct {
	OrganizationID int64 `path:"organizationId" doc:"Organization ID"`
	Page           int   `query:"page" default:"1" doc:"Page number"`
	PageSize       int   `query:"pageSize" default:"10" doc:"Page size"`
}

type ListTenantsByOrganizationOutput struct {
	Body response.TenantListResponse
}

func (c *TenantController) ListTenantsByOrganization(ctx context.Context, input *ListTenantsByOrganizationInput) (*ListTenantsByOrganizationOutput, error) {
	resp, err := c.usecase.ListTenantsByOrganization(ctx, input.OrganizationID, input.Page, input.PageSize)
	if err != nil {
		return nil, err
	}

	return &ListTenantsByOrganizationOutput{Body: *resp}, nil
}

type CreateTenantInput struct {
	Body requests.CreateTenantRequest
}

type CreateTenantOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) CreateTenant(ctx context.Context, input *CreateTenantInput) (*CreateTenantOutput, error) {
	resp, err := c.usecase.CreateTenant(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateTenantOutput{Body: *resp}, nil
}

type UpdateTenantInput struct {
	ID   int64 `path:"id" doc:"Tenant ID"`
	Body requests.UpdateTenantRequest
}

type UpdateTenantOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) UpdateTenant(ctx context.Context, input *UpdateTenantInput) (*UpdateTenantOutput, error) {
	resp, err := c.usecase.UpdateTenant(ctx, input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &UpdateTenantOutput{Body: *resp}, nil
}

type DeleteTenantInput struct {
	ID int64 `path:"id" doc:"Tenant ID"`
}

type DeleteTenantOutput struct {
	Success bool `json:"success"`
}

func (c *TenantController) DeleteTenant(ctx context.Context, input *DeleteTenantInput) (*DeleteTenantOutput, error) {
	err := c.usecase.DeleteTenant(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteTenantOutput{Success: true}, nil
}
