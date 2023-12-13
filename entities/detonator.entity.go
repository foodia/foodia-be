package entities

import (
	"time"
)

type Detonator struct {
	ID        int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	UserId    int       `gorm:"type:int(11);not null;unique" json:"user_id"`
	SelfPhoto string    `gorm:"type:varchar(100);not null" json:"self_photo"`
	KTPPhoto  string    `gorm:"type:varchar(100);not null" json:"ktp_photo"`
	KTPNumber string    `gorm:"type:varchar(16);not null" json:"ktp_number"`
	Status    string    `gorm:"default:'waiting'" json:"status"`
	Note      string    `gorm:"type:text" json:"note"`
	CreatedAt time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`

	Oauth *Oauth `gorm:"foreignKey:ID;references:UserId" json:"oauth"`
}
