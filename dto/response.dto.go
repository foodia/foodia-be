package dto

import "github.com/gofiber/fiber/v2"

type ApiResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   any    `json:"error,omitempty"`
	Body    any    `json:"body"`
	Meta    any    `json:"meta,omitempty"`
}

type ApiFieldError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type ApiError struct {
	StatusCode *fiber.Error
	Message    string
}
