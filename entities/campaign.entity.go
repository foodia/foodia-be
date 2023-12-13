package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Campaign struct {
	ID             int             `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
	DetonatorID    int             `json:"detonator_id"`
	EventName      string          `json:"event_name"`
	EventType      string          `json:"event_type"`
	EventDate      string          `json:"event_date"`
	EventTime      string          `json:"event_time"`
	Description    string          `json:"description"`
	DonationTarget decimal.Decimal `json:"donation_target"`
	Province       string          `json:"province"`
	City           string          `json:"city"`
	SubDistrict    string          `json:"sub_district"`
	PostalCode     string          `json:"postal_code"`
	Address        string          `json:"address"`
	Latitude       string          `json:"latitude"`
	Longitude      string          `json:"longitude"`
	Status         string          `gorm:"default:'waiting'" json:"status"`
	IsActive       bool            `gorm:"default:false" json:"is_active"`
	ImageURL       string          `json:"image_url"`
	CreatedAt      time.Time       `gorm:"default:current_timestamp()"  json:"created_at"`
	UpdatedAt      time.Time       `gorm:"default:current_timestamp()" json:"updated_at"`

	Detonator *Detonator `gorm:"foreignKey:ID;references:DetonatorID" json:"detonator"`
}
