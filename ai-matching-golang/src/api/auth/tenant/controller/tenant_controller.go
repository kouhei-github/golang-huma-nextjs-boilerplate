package controller

import (
	"ai-matching/src/api/auth/tenant/requests"
	"ai-matching/src/api/auth/tenant/response"
	"ai-matching/src/api/auth/tenant/usecase"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TenantController struct {
	usecase *usecase.TenantUsecase
}

func NewTenantController(tenantUsecase *usecase.TenantUsecase) *TenantController {
	return &TenantController{
		usecase: tenantUsecase,
	}
}

type GetTenantInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	ID             uuid.UUID `path:"tenantId" doc:"Tenant ID"`
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
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	Page           int       `query:"page" default:"1" doc:"Page number"`
	PageSize       int       `query:"pageSize" default:"10" doc:"Page size"`
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
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	Body           requests.CreateTenantRequest
}

type CreateTenantOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) CreateTenant(ctx context.Context, input *CreateTenantInput) (*CreateTenantOutput, error) {
	// Set organization ID from path parameter
	input.Body.OrganizationID = input.OrganizationID
	resp, err := c.usecase.CreateTenant(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateTenantOutput{Body: *resp}, nil
}

type UpdateTenantInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	ID             uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	Body           requests.UpdateTenantRequest
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
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	ID             uuid.UUID `path:"tenantId" doc:"Tenant ID"`
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

// Organization-scoped tenant endpoints

type GetTenantInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

type GetTenantInOrganizationOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) GetTenantInOrganization(ctx context.Context, input *GetTenantInOrganizationInput) (*GetTenantInOrganizationOutput, error) {
	// First get the tenant
	resp, err := c.usecase.GetTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	// Verify it belongs to the organization
	if resp.OrganizationID != input.OrganizationID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Tenant not found in organization")
	}

	return &GetTenantInOrganizationOutput{Body: *resp}, nil
}

type CreateTenantInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	Body           requests.CreateTenantRequest
}

type CreateTenantInOrganizationOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) CreateTenantInOrganization(ctx context.Context, input *CreateTenantInOrganizationInput) (*CreateTenantInOrganizationOutput, error) {
	// Set the organization ID from the path
	input.Body.OrganizationID = input.OrganizationID

	resp, err := c.usecase.CreateTenant(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateTenantInOrganizationOutput{Body: *resp}, nil
}

type UpdateTenantInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	Body           requests.UpdateTenantRequest
}

type UpdateTenantInOrganizationOutput struct {
	Body response.TenantResponse
}

func (c *TenantController) UpdateTenantInOrganization(ctx context.Context, input *UpdateTenantInOrganizationInput) (*UpdateTenantInOrganizationOutput, error) {
	// First get the tenant to verify organization
	tenant, err := c.usecase.GetTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	// Verify it belongs to the organization
	if tenant.OrganizationID != input.OrganizationID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Tenant not found in organization")
	}

	resp, err := c.usecase.UpdateTenant(ctx, input.TenantID, input.Body)
	if err != nil {
		return nil, err
	}

	return &UpdateTenantInOrganizationOutput{Body: *resp}, nil
}

type DeleteTenantInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

type DeleteTenantInOrganizationOutput struct {
	Success bool `json:"success"`
}

func (c *TenantController) DeleteTenantInOrganization(ctx context.Context, input *DeleteTenantInOrganizationInput) (*DeleteTenantInOrganizationOutput, error) {
	// First get the tenant to verify organization
	tenant, err := c.usecase.GetTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	// Verify it belongs to the organization
	if tenant.OrganizationID != input.OrganizationID {
		return nil, fiber.NewError(fiber.StatusNotFound, "Tenant not found in organization")
	}

	err = c.usecase.DeleteTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	return &DeleteTenantInOrganizationOutput{Success: true}, nil
}
