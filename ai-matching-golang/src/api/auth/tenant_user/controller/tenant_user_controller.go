package controller

import (
	"ai-matching/src/api/auth/tenant_user/requests"
	"ai-matching/src/api/auth/tenant_user/response"
	"ai-matching/src/api/auth/tenant_user/usecase"
	"context"

	"github.com/gofiber/fiber/v2"
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
	OrganizationID uuid.UUID                       `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID                       `path:"tenantId" doc:"Tenant ID"`
	Body           requests.AddUserToTenantRequest `doc:"Add user to tenant request"`
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
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
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

// Renamed from GetTenantUsersInput to ListTenantUsersInput for consistency
type ListTenantUsersInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

type ListTenantUsersOutput struct {
	Body response.TenantUsersResponse
}

func (c *TenantUserController) ListTenantUsers(ctx context.Context, input *ListTenantUsersInput) (*ListTenantUsersOutput, error) {
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

	return &ListTenantUsersOutput{
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
	OrganizationID uuid.UUID                      `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID                      `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID                      `path:"userId" doc:"User ID"`
	Body           requests.UpdateUserRoleRequest `doc:"Update user role request"`
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
	Body response.UsersListResponse
}

// New methods for organization-scoped endpoints

type GetTenantUserInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type GetTenantUserOutput struct {
	Body response.TenantUserInfo
}

func (c *TenantUserController) GetTenantUser(ctx context.Context, input *GetTenantUserInput) (*GetTenantUserOutput, error) {
	// Get all users in the tenant and find the specific one
	users, err := c.usecase.ListTenantUsers(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.UserID == input.UserID {
			return &GetTenantUserOutput{
				Body: response.TenantUserInfo{
					UserID:    u.UserID,
					Email:     u.Email,
					FirstName: u.FirstName.String,
					LastName:  u.LastName.String,
					Role:      u.Role.String,
				},
			}, nil
		}
	}

	return nil, fiber.NewError(fiber.StatusNotFound, "User not found in tenant")
}

type GetUserTenantsInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type GetUserTenantsInOrganizationOutput struct {
	Body response.UserTenantsResponse
}

func (c *TenantUserController) GetUserTenantsInOrganization(ctx context.Context, input *GetUserTenantsInOrganizationInput) (*GetUserTenantsInOrganizationOutput, error) {
	// Get all tenants for the user
	tenants, err := c.usecase.GetTenantsByUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// Filter by organization
	var tenantList []response.TenantDetails
	for _, t := range tenants {
		// TODO: Need to check if tenant belongs to organization
		tenantList = append(tenantList, response.TenantDetails{
			ID:        t.ID,
			Name:      t.Name,
			Subdomain: t.Subdomain,
		})
	}

	return &GetUserTenantsInOrganizationOutput{
		Body: response.UserTenantsResponse{
			Tenants: tenantList,
		},
	}, nil
}

// Organization-scoped methods

type ListTenantUsersInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
}

type ListTenantUsersInOrganizationOutput struct {
	Body response.TenantUsersResponse
}

func (c *TenantUserController) ListTenantUsersInOrganization(ctx context.Context, input *ListTenantUsersInOrganizationInput) (*ListTenantUsersInOrganizationOutput, error) {
	// TODO: Verify tenant belongs to organization
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

	return &ListTenantUsersInOrganizationOutput{
		Body: response.TenantUsersResponse{
			Users: userList,
		},
	}, nil
}

type GetTenantUserInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type GetTenantUserInOrganizationOutput struct {
	Body response.TenantUserInfo
}

func (c *TenantUserController) GetTenantUserInOrganization(ctx context.Context, input *GetTenantUserInOrganizationInput) (*GetTenantUserInOrganizationOutput, error) {
	// TODO: Verify tenant belongs to organization
	// Get the specific user in the tenant
	users, err := c.usecase.ListTenantUsers(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.UserID == input.UserID {
			return &GetTenantUserInOrganizationOutput{
				Body: response.TenantUserInfo{
					UserID:    u.UserID,
					Email:     u.Email,
					FirstName: u.FirstName.String,
					LastName:  u.LastName.String,
					Role:      u.Role.String,
				},
			}, nil
		}
	}

	return nil, fiber.NewError(fiber.StatusNotFound, "User not found in tenant")
}

type AddUserToTenantInOrganizationInput struct {
	OrganizationID uuid.UUID                       `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID                       `path:"tenantId" doc:"Tenant ID"`
	Body           requests.AddUserToTenantRequest `doc:"Add user to tenant request"`
}

type AddUserToTenantInOrganizationOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) AddUserToTenantInOrganization(ctx context.Context, input *AddUserToTenantInOrganizationInput) (*AddUserToTenantInOrganizationOutput, error) {
	// TODO: Verify tenant belongs to organization
	err := c.usecase.AddUserToTenant(ctx, input.TenantID, input.Body.UserID, input.Body.Role)
	if err != nil {
		return nil, err
	}

	return &AddUserToTenantInOrganizationOutput{
		Body: response.MessageResponse{
			Message: "User added to tenant successfully",
		},
	}, nil
}

type UpdateUserRoleInOrganizationInput struct {
	OrganizationID uuid.UUID                      `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID                      `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID                      `path:"userId" doc:"User ID"`
	Body           requests.UpdateUserRoleRequest `doc:"Update user role request"`
}

type UpdateUserRoleInOrganizationOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) UpdateUserRoleInOrganization(ctx context.Context, input *UpdateUserRoleInOrganizationInput) (*UpdateUserRoleInOrganizationOutput, error) {
	err := c.usecase.UpdateUserRoleInTenant(ctx, input.TenantID, input.UserID, input.Body.Role)
	if err != nil {
		return nil, err
	}

	return &UpdateUserRoleInOrganizationOutput{
		Body: response.MessageResponse{
			Message: "User role updated successfully",
		},
	}, nil
}

type RemoveUserFromTenantInOrganizationInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	TenantID       uuid.UUID `path:"tenantId" doc:"Tenant ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type RemoveUserFromTenantInOrganizationOutput struct {
	Body response.MessageResponse
}

func (c *TenantUserController) RemoveUserFromTenantInOrganization(ctx context.Context, input *RemoveUserFromTenantInOrganizationInput) (*RemoveUserFromTenantInOrganizationOutput, error) {
	err := c.usecase.RemoveUserFromTenant(ctx, input.TenantID, input.UserID)
	if err != nil {
		return nil, err
	}

	return &RemoveUserFromTenantInOrganizationOutput{
		Body: response.MessageResponse{
			Message: "User removed from tenant successfully",
		},
	}, nil
}
