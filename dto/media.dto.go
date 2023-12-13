package dto

import "mime/multipart"

type MediaRequest struct {
	Destination string                `json:"destination" form:"destination" validate:"required,oneof=campaign merchant detonator product"`
	File        *multipart.FileHeader `json:"file" form:"file" validate:"required"`
}

type MediaResponse struct {
	Destination string `json:"destination"`
	FileUrl     string `json:"file_url"`
}
