package controller

import (
	"ai-matching/src/api/auth/user/requests"
	"ai-matching/src/api/auth/user/response"
	"ai-matching/src/api/auth/user/usecase"
	"ai-matching/src/infrastructure/middleware"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserController struct {
	usecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{
		usecase: userUsecase,
	}
}

type GetUserInput struct {
	ID uuid.UUID `path:"userId" doc:"User ID"`
}

type GetUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) GetUser(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
	resp, err := c.usecase.GetUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{Body: *resp}, nil
}

type ListUsersInput struct {
	Page     int `query:"page" default:"1" doc:"Page number"`
	PageSize int `query:"pageSize" default:"10" doc:"Page size"`
}

type ListUsersOutput struct {
	Body response.UserListResponse
}

func (c *UserController) ListUsers(ctx context.Context, input *ListUsersInput) (*ListUsersOutput, error) {
	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(user.OrganizationID)

	resp, err := c.usecase.ListUsers(ctx, input.Page, input.PageSize)
	if err != nil {
		return nil, err
	}

	return &ListUsersOutput{Body: *resp}, nil
}

type CreateUserInput struct {
	Body requests.CreateUserRequest
}

type CreateUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	resp, err := c.usecase.CreateUser(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{Body: *resp}, nil
}

type UpdateUserInput struct {
	ID   uuid.UUID `path:"userId" doc:"User ID"`
	Body requests.UpdateUserRequest
}

type UpdateUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) UpdateUser(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
	resp, err := c.usecase.UpdateUser(ctx, input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{Body: *resp}, nil
}

type DeleteUserInput struct {
	ID uuid.UUID `path:"userId" doc:"User ID"`
}

type DeleteUserOutput struct {
	Success bool `json:"success"`
}

func (c *UserController) DeleteUser(ctx context.Context, input *DeleteUserInput) (*DeleteUserOutput, error) {
	err := c.usecase.DeleteUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteUserOutput{Success: true}, nil
}

// Organization-scoped user endpoints

type GetOrganizationUserInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type GetOrganizationUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) GetOrganizationUser(ctx context.Context, input *GetOrganizationUserInput) (*GetOrganizationUserOutput, error) {
	// Get the user
	resp, err := c.usecase.GetUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetOrganizationUserOutput{Body: *resp}, nil
}

type ListOrganizationUsersInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	Page           int       `query:"page" default:"1" doc:"Page number"`
	PageSize       int       `query:"pageSize" default:"10" doc:"Page size"`
}

type ListOrganizationUsersOutput struct {
	Body response.UserListResponse
}

func (c *UserController) ListOrganizationUsers(ctx context.Context, input *ListOrganizationUsersInput) (*ListOrganizationUsersOutput, error) {
	resp, err := c.usecase.ListUsers(ctx, input.Page, input.PageSize)
	if err != nil {
		return nil, err
	}

	return &ListOrganizationUsersOutput{Body: *resp}, nil
}

type CreateOrganizationUserInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	Body           requests.CreateUserRequest
}

type CreateOrganizationUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) CreateOrganizationUser(ctx context.Context, input *CreateOrganizationUserInput) (*CreateOrganizationUserOutput, error) {
	resp, err := c.usecase.CreateUser(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &CreateOrganizationUserOutput{Body: *resp}, nil
}

type UpdateOrganizationUserInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
	Body           requests.UpdateUserRequest
}

type UpdateOrganizationUserOutput struct {
	Body response.UserResponse
}

func (c *UserController) UpdateOrganizationUser(ctx context.Context, input *UpdateOrganizationUserInput) (*UpdateOrganizationUserOutput, error) {
	resp, err := c.usecase.UpdateUser(ctx, input.UserID, input.Body)
	if err != nil {
		return nil, err
	}

	return &UpdateOrganizationUserOutput{Body: *resp}, nil
}

type DeleteOrganizationUserInput struct {
	OrganizationID uuid.UUID `path:"organizationId" doc:"Organization ID"`
	UserID         uuid.UUID `path:"userId" doc:"User ID"`
}

type DeleteOrganizationUserOutput struct {
	Success bool `json:"success"`
}

func (c *UserController) DeleteOrganizationUser(ctx context.Context, input *DeleteOrganizationUserInput) (*DeleteOrganizationUserOutput, error) {
	err := c.usecase.DeleteUser(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &DeleteOrganizationUserOutput{Success: true}, nil
}
