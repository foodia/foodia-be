package dto

type MerchantProductRequest struct {
	MerchantID  int     `json:"merchant_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	QTY         int     `json:"qty" validate:"required"`
	Images      []struct {
		ImageURL string `json:"image_url"`
	} `json:"images"`
}
