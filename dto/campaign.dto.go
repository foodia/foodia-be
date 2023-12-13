package dto

type CampaignRequest struct {
	DetonatorID    int     `json:"detonator_id" validate:"required"`
	EventName      string  `json:"event_name" validate:"required"`
	EventType      string  `json:"event_type" validate:"required"`
	EventDate      string  `json:"event_date" validate:"required"`
	EventTime      string  `json:"event_time" validate:"required"`
	Description    string  `json:"description" validate:"required"`
	DonationTarget float64 `json:"donation_target" validate:"required"`
	Province       string  `json:"province" validate:"required"`
	City           string  `json:"city" validate:"required"`
	SubDistrict    string  `json:"sub_district" validate:"required"`
	PostalCode     string  `json:"postal_code" validate:"required"`
	Address        string  `json:"address" validate:"required"`
	Latitude       string  `json:"latitude" validate:"required"`
	Longitude      string  `json:"longitude" validate:"required"`
	ImageURL       string  `json:"image_url" validate:"required"`
	Products       []struct {
		MerchantProductID int `json:"merchant_product_id"`
	} `json:"products"`
}
