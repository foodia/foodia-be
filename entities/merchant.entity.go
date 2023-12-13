package entities

import (
	"time"
)

type Merchant struct {
	ID          int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	Province    string    `gorm:"type:varchar(100)" json:"province"`
	City        string    `gorm:"type:varchar(100)" json:"city"`
	SubDistrict string    `gorm:"type:varchar(100)" json:"sub_district"`
	PostalCode  string    `gorm:"type:varchar(50)" json:"postal_code"`
	Address     string    `gorm:"type:text" json:"address"`
	Latitude    string    `gorm:"type:varchar(100)" json:"latitude"`
	Longitude   string    `gorm:"type:varchar(100)" json:"longitude"`
	NoLinkAja   string    `gorm:"type:varchar(15)" json:"no_link_aja"`
	UserId      int       `gorm:"type:int(11);not null;unique" json:"user_id"`
	SelfPhoto   string    `gorm:"type:varchar(100);not null" json:"self_photo"`
	KTPPhoto    string    `gorm:"type:varchar(100);not null" json:"ktp_photo"`
	KTPNumber   string    `gorm:"type:varchar(16);not null" json:"ktp_number"`
	Status      string    `gorm:"default:'waiting'" json:"status"`
	Note        string    `gorm:"type:text" json:"note"`
	CreatedAt   time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp()" json:"updated_at"`

	Oauth           *Oauth            `gorm:"foreignKey:ID;references:UserId" json:"oauth"`
	MerchantProduct []MerchantProduct `gorm:"foreignKey:MerchantID;references:ID" json:"products"`
}
