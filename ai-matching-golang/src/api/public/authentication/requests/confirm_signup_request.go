package requests

type ConfirmSignUpRequest struct {
	Email            string `json:"email" validate:"required,email" doc:"User email"`
	ConfirmationCode string `json:"confirmationCode" validate:"required" doc:"Confirmation code from email"`
}