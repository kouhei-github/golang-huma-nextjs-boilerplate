package response

import "time"

type OrganizationResponse struct {
	ID          int64     `json:"id" doc:"Organization ID"`
	Name        string    `json:"name" doc:"Organization name"`
	Description string    `json:"description" doc:"Organization description"`
	IsActive    bool      `json:"isActive" doc:"Is organization active"`
	CreatedAt   time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt   time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations" doc:"List of organizations"`
	Total         int                    `json:"total" doc:"Total count"`
	Page          int                    `json:"page" doc:"Current page"`
	PageSize      int                    `json:"pageSize" doc:"Page size"`
}