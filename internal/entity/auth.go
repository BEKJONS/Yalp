package entity

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"` // consider using the Platform constants for type safety
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmail struct {
	Email    string `json:"email"`
	Otp      string `json:"otp"`
	Platform string `json:"platform"` // consider using the Platform constants for type safety
}
