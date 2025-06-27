package controller

import (
	"ai-matching/src/api/auth/tenant_user/requests"
	"ai-matching/src/api/auth/tenant_user/response"
	"ai-matching/src/api/auth/tenant_user/usecase"
	"context"

	"github.com/google/uuid"
)

type TenantUserController struct {
	usecase *usecase.TenantUserUsecase
}

func NewTenantUserController(tenantUserUsecase *usecase.TenantUserUsecase) *TenantUserController {
	return &TenantUserController{
		usecase: tenantUserUsecase,
	}
}

type AddUserToTenantInput struct {
	TenantID uuid.UUID                      `path:"tenantId" doc:"Tenant ID"`
	Body     requests.AddUserToTenantRequest `doc:"Add user to tenant request"`
}

type AddUserToTenantOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) AddUserToTenant(ctx context.Context, input *AddUserToTenantInput) (*AddUserToTenantOutput, error) {
	err := c.usecase.AddUserToTenant(ctx, input.TenantID, input.Body.UserID, input.Body.Role)
	if err != nil {
		return nil, err
	}

	return &AddUserToTenantOutput{
		Body: response.MessageResponse{
			Message: "User added to tenant successfully",
		},
	}, nil
}

type RemoveUserFromTenantInput struct {
	TenantID uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	UserID   uuid.UUID `path:"userId" doc:"User ID"`
}

type RemoveUserFromTenantOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) RemoveUserFromTenant(ctx context.Context, input *RemoveUserFromTenantInput) (*RemoveUserFromTenantOutput, error) {
	err := c.usecase.RemoveUserFromTenant(ctx, input.TenantID, input.UserID)
	if err != nil {
		return nil, err
	}

	return &RemoveUserFromTenantOutput{
		Body: response.MessageResponse{
			Message: "User removed from tenant successfully",
		},
	}, nil
}

type GetTenantUsersInput struct {
	TenantID uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

type GetTenantUsersOutput struct {
	Body response.TenantUsersResponse
}

func (c *TenantUserController) GetTenantUsers(ctx context.Context, input *GetTenantUsersInput) (*GetTenantUsersOutput, error) {
	users, err := c.usecase.ListTenantUsers(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	var userList []response.TenantUserInfo
	for _, u := range users {
		userList = append(userList, response.TenantUserInfo{
			UserID:    u.UserID,
			Email:     u.Email,
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
			Role:      u.Role.String,
		})
	}

	return &GetTenantUsersOutput{
		Body: response.TenantUsersResponse{
			Users: userList,
		},
	}, nil
}

type GetUserTenantsInput struct {
	UserID uuid.UUID `path:"userId" doc:"User ID"`
}

type GetUserTenantsOutput struct {
	Body response.UserTenantsResponse
}

func (c *TenantUserController) GetUserTenants(ctx context.Context, input *GetUserTenantsInput) (*GetUserTenantsOutput, error) {
	tenants, err := c.usecase.GetTenantsByUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	var tenantList []response.TenantDetails
	for _, t := range tenants {
		tenantList = append(tenantList, response.TenantDetails{
			ID:        t.ID,
			Name:      t.Name,
			Subdomain: t.Subdomain,
		})
	}

	return &GetUserTenantsOutput{
		Body: response.UserTenantsResponse{
			Tenants: tenantList,
		},
	}, nil
}

type UpdateUserRoleInput struct {
	TenantID uuid.UUID                   `path:"tenantId" doc:"Tenant ID"`
	UserID   uuid.UUID                   `path:"userId" doc:"User ID"`
	Body     requests.UpdateUserRoleRequest `doc:"Update user role request"`
}

type UpdateUserRoleOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) UpdateUserRole(ctx context.Context, input *UpdateUserRoleInput) (*UpdateUserRoleOutput, error) {
	err := c.usecase.UpdateUserRoleInTenant(ctx, input.TenantID, input.UserID, input.Body.Role)
	if err != nil {
		return nil, err
	}

	return &UpdateUserRoleOutput{
		Body: response.MessageResponse{
			Message: "User role updated successfully",
		},
	}, nil
}

type GetUsersNotInTenantInput struct {
	TenantID uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	Page     int       `query:"page" default:"1" doc:"Page number"`
	PageSize int       `query:"pageSize" default:"20" doc:"Page size"`
}

type GetUsersNotInTenantOutput struct {
	Body response.UsersResponse
}

func (c *TenantUserController) GetUsersNotInTenant(ctx context.Context, input *GetUsersNotInTenantInput) (*GetUsersNotInTenantOutput, error) {
	limit := int32(input.PageSize)
	offset := int32((input.Page - 1) * input.PageSize)

	users, err := c.usecase.GetUsersNotInTenant(ctx, input.TenantID, limit, offset)
	if err != nil {
		return nil, err
	}

	var userList []response.UserDetails
	for _, u := range users {
		userList = append(userList, response.UserDetails{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
		})
	}

	return &GetUsersNotInTenantOutput{
		Body: response.UsersResponse{
			Users: userList,
		},
	}, nil
}