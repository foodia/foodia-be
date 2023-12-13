package dto

import "mime/multipart"

type MerchantRegistration struct {
	Fullname    string                `json:"fullname" form:"fullname" validate:"required"`
	Phone       string                `json:"phone" form:"phone" validate:"required"`
	Email       string                `json:"email" form:"email" validate:"required"`
	Password    string                `json:"password" form:"password" validate:"required"`
	Province    string                `json:"province" form:"province" validate:"required"`
	City        string                `json:"city" form:"city" validate:"required"`
	SubDistrict string                `json:"sub_district" form:"sub_district" validate:"required"`
	PostalCode  string                `json:"postal_code" form:"postal_code" validate:"required"`
	Address     string                `json:"address" form:"address" validate:"required"`
	Latitude    string                `json:"latitude" form:"latitude" validate:"required"`
	Longitude   string                `json:"longitude" form:"longitude" validate:"required"`
	NoLinkAja   string                `json:"no_link_aja" form:"no_link_aja" validate:"required"`
	KTPNumber   string                `json:"ktp_number" form:"ktp_number" validate:"required"`
	SelfPhoto   *multipart.FileHeader `json:"self_photo" form:"self_photo" validate:"required"`
	KTPPhoto    *multipart.FileHeader `json:"ktp_photo" form:"ktp_photo" validate:"required"`
}

type MerchantApproval struct {
	Status string `json:"status" validate:"required,oneof=waiting approved rejected"`
	Note   string `json:"note"`
}

type MerchantUpdate struct {
	Fullname    string                `json:"fullname" form:"fullname"`
	Phone       string                `json:"phone" form:"phone"`
	Email       string                `json:"email" form:"email"`
	Password    string                `json:"password" form:"password"`
	Province    string                `json:"province" form:"province"`
	City        string                `json:"city" form:"city"`
	SubDistrict string                `json:"sub_district" form:"sub_district"`
	PostalCode  string                `json:"postal_code" form:"postal_code"`
	Address     string                `json:"address" form:"address"`
	Latitude    string                `json:"latitude" form:"latitude"`
	Longitude   string                `json:"longitude" form:"longitude"`
	NoLinkAja   string                `json:"no_link_aja" form:"no_link_aja"`
	KTPNumber   string                `json:"ktp_number" form:"ktp_number"`
	SelfPhoto   *multipart.FileHeader `json:"self_photo" form:"self_photo"`
	KTPPhoto    *multipart.FileHeader `json:"ktp_photo" form:"ktp_photo"`
}
