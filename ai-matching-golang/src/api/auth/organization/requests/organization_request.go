package requests

type CreateOrganizationRequest struct {
	Name        string `json:"name" validate:"required" doc:"Organization name"`
	Description string `json:"description,omitempty" doc:"Organization description"`
	IsActive    bool   `json:"isActive" doc:"Is organization active"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name" validate:"required" doc:"Organization name"`
	Description string `json:"description,omitempty" doc:"Organization description"`
	IsActive    bool   `json:"isActive" doc:"Is organization active"`
}