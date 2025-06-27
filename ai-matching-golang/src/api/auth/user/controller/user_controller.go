package controller

import (
	"ai-matching/src/api/auth/user/requests"
	"ai-matching/src/api/auth/user/response"
	"ai-matching/src/api/auth/user/usecase"
	"context"
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
	ID int64 `path:"id" doc:"User ID"`
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
	ID   int64 `path:"id" doc:"User ID"`
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
	ID int64 `path:"id" doc:"User ID"`
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
