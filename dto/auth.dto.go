package dto

import "time"

type AuthRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type AuthResponse struct {
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     string `json:"role,omitempty"`
	Token    string `json:"token,omitempty"`
	User     *User  `json:"user"`
}

type User struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type OTPRequest struct {
	OTPCode   string    `json:"otp_code"`
	ExpiredAt time.Time `json:"expired_at"`
	Email     string    `json:"email"`
}

type ValidateOTPRequest struct {
	Email string `validate:"required,email" json:"email,omitempty"`
	Code  string `validate:"required,numeric" json:"code,omitempty"`
}
