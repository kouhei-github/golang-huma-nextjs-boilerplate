package requests

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" doc:"User email"`
}
