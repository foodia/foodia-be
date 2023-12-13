package entities

import (
	"time"
)

type Oauth struct {
	ID        int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	UserId    string    `gorm:"type:varchar(100);not null;unique" json:"user_id"`
	Fullname  string    `gorm:"type:varchar(100);not null" json:"fullname"`
	Email     string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	Phone     string    `gorm:"type:varchar(15);not null;unique" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Role      string    `gorm:"type:varchar(100);not null;default:'superadmin'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	IsLocked  bool      `gorm:"default:false" json:"is_locked"`
	CreatedAt time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}
