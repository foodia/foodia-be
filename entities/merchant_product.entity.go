package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type MerchantProduct struct {
	ID          int             `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	MerchantID  int             `json:"merchant_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	QTY         int             `json:"qty"`
	Status      string          `gorm:"default:'waiting'" json:"status"`
	IsActive    bool            `gorm:"default:false" json:"is_active"`
	CreatedAt   time.Time       `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt   time.Time       `gorm:"default:current_timestamp()" json:"updated_at"`
	DeletedAt   time.Time       `json:"deleted_at"`

	MerchantProductImage []MerchantProductImage `gorm:"foreignKey:MerchantProductID;references:ID" json:"images"`
}
