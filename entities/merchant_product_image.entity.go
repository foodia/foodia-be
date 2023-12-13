package entities

import (
	"time"
)

type MerchantProductImage struct {
	ID                int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	MerchantProductID int       `json:"merchant_product_id"`
	ImageURL          string    `json:"image_url"`
	CreatedAt         time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}
