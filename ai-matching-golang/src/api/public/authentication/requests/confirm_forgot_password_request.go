package requests

type ConfirmForgotPasswordRequest struct {
	Email            string `json:"email" validate:"required,email" doc:"User email"`
	Password         string `json:"password" validate:"required,min=8" doc:"New password"`
	ConfirmationCode string `json:"confirmationCode" validate:"required" doc:"Confirmation code from email"`
}
