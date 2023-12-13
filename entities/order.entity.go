package entities

import (
	"time"
)

type Order struct {
	ID                int       `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	CampaignID        int       `json:"campaign_id"`
	MerchantProductID int       `json:"merchant_product_id"`
	OrderStatus       string    `gorm:"default:'waiting'" json:"order_status"`
	Note              string    `json:"note"`
	CreatedAt         time.Time `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
}
