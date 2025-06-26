package response

type HealthResponse struct {
	Status  string `json:"status" example:"ok" doc:"Health status"`
	Version string `json:"version" example:"1.0.0" doc:"API version"`
}