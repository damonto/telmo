package auth

type VerifyOTPRequest struct {
	Code string `json:"code" validate:"required,len=6,numeric"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}

type OTPRequirementResponse struct {
	OTPRequired bool `json:"otpRequired"`
}
