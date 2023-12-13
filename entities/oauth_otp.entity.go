package entities

import (
	"time"
)

type OauthOTP struct {
	ID        int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:int(11)" json:"email"`
	OTPCode   string    `gorm:"type:varchar(6);not null" json:"otp_code"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}
