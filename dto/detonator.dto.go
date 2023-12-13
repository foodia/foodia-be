package dto

import "mime/multipart"

type DetonatorRegistration struct {
	Fullname  string                `json:"fullname" form:"fullname" validate:"required"`
	Phone     string                `json:"phone" form:"phone" validate:"required"`
	Email     string                `json:"email" form:"email" validate:"required"`
	Password  string                `json:"password" form:"password" validate:"required"`
	KTPNumber string                `json:"ktp_number" form:"ktp_number" validate:"required"`
	SelfPhoto *multipart.FileHeader `json:"self_photo" form:"self_photo" validate:"required"`
	KTPPhoto  *multipart.FileHeader `json:"ktp_photo" form:"ktp_photo" validate:"required"`
}

type DetonatorApproval struct {
	Status string `json:"status" validate:"required,oneof=waiting approved rejected"`
	Note   string `json:"note"`
}

type DetonatorUpdate struct {
	Fullname  string                `json:"fullname" form:"fullname"`
	Phone     string                `json:"phone" form:"phone"`
	Email     string                `json:"email" form:"email"`
	Password  string                `json:"password" form:"password"`
	KTPNumber string                `json:"ktp_number" form:"ktp_number"`
	SelfPhoto *multipart.FileHeader `json:"self_photo" form:"self_photo"`
	KTPPhoto  *multipart.FileHeader `json:"ktp_photo" form:"ktp_photo"`
}
