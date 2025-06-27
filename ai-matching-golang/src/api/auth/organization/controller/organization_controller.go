package controller

import (
	"ai-matching/src/api/auth/organization/requests"
	"ai-matching/src/api/auth/organization/response"
	"ai-matching/src/api/auth/organization/usecase"
	"context"

	"github.com/google/uuid"
)

type OrganizationController struct {
	usecase *usecase.OrganizationUsecase
}

func NewOrganizationController(orgUsecase *usecase.OrganizationUsecase) *OrganizationController {
	return &OrganizationController{
		usecase: orgUsecase,
	}
}

type GetOrganizationInput struct {
	ID uuid.UUID `path:"id" doc:"Organization ID"`
}

type GetOrganizationOutput struct {
	Body response.OrganizationResponse
}

func (c *OrganizationController) GetOrganization(ctx context.Context, input *GetOrganizationInput) (*GetOrganizationOutput, error) {
	resp, err := c.usecase.GetOrganization(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetOrganizationOutput{Body: *resp}, nil
}

type ListOrganizationsInput struct {
	Page     int `query:"page" default:"1" doc:"Page number"`
	PageSize int `query:"pageSize" default:"10" doc:"Page size"`
}

type ListOrganizationsOutput struct {
	Body response.OrganizationListResponse
}

func (c *OrganizationController) ListOrganizations(ctx context.Context, input *ListOrganizationsInput) (*ListOrganizationsOutput, error) {
	resp, err := c.usecase.ListOrganizations(ctx, input.Page, input.PageSize)
	if err != nil {
		return nil, err
	}

	return &ListOrganizationsOutput{Body: *resp}, nil
}

type CreateOrganizationInput struct {
	Body requests.CreateOrganizationRequest
}

type CreateOrganizationOutput struct {
	Body response.OrganizationResponse
}

func (c *OrganizationController) CreateOrganization(ctx context.Context, input *CreateOrganizationInput) (*CreateOrganizationOutput, error) {
	resp, err := c.usecase.CreateOrganization(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateOrganizationOutput{Body: *resp}, nil
}

type UpdateOrganizationInput struct {
	ID   uuid.UUID `path:"id" doc:"Organization ID"`
	Body requests.UpdateOrganizationRequest
}

type UpdateOrganizationOutput struct {
	Body response.OrganizationResponse
}

func (c *OrganizationController) UpdateOrganization(ctx context.Context, input *UpdateOrganizationInput) (*UpdateOrganizationOutput, error) {
	resp, err := c.usecase.UpdateOrganization(ctx, input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &UpdateOrganizationOutput{Body: *resp}, nil
}

type DeleteOrganizationInput struct {
	ID uuid.UUID `path:"id" doc:"Organization ID"`
}

type DeleteOrganizationOutput struct {
	Success bool `json:"success"`
}

func (c *OrganizationController) DeleteOrganization(ctx context.Context, input *DeleteOrganizationInput) (*DeleteOrganizationOutput, error) {
	err := c.usecase.DeleteOrganization(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteOrganizationOutput{Success: true}, nil
}
