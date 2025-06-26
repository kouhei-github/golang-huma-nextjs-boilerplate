package controller

import (
	"ai-matching/src/api/public/health/response"
	"context"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

type HealthOutput struct {
	Body response.HealthResponse
}

func (c *HealthController) GetHealth(ctx context.Context, input *struct{}) (*HealthOutput, error) {
	return &HealthOutput{
		Body: response.HealthResponse{
			Status:  "ok",
			Version: "1.0.0",
		},
	}, nil
}
