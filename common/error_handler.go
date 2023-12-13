package common

import (
	"errors"

	"foodia-be/dto"
	"foodia-be/enums"

	"github.com/gofiber/fiber/v2"
)

func getStatusCode(err error) int {
	if err == nil {
		return fiber.StatusOK
	}
	switch err {
	case enums.ErrInternalServor:
		return fiber.StatusInternalServerError
	case enums.ErrNotFound:
		return fiber.StatusNotFound
	case enums.ErrBadParamInput, enums.ErrIncorrectCredential, enums.ErrInvalidRefreshToken:
		return fiber.StatusBadRequest
	case enums.ErrAccessForbidden:
		return fiber.StatusForbidden
	case enums.ErrUnauthorized, enums.ErrInvalidToken, enums.ErrExpiredToken:
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}

// ErrorHandler implements centralized error handling & returns JSON with status code and error message
// Source : https://docs.gofiber.io/guide/error-handling/
func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var code int
		var errorMsg string

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			errorMsg = fiberErr.Message
		} else {
			code = getStatusCode(err)
			errorMsg = err.Error()
		}

		return ctx.Status(code).JSON(dto.ApiResponse{
			Code:    code,
			Message: errorMsg,
		})
	}
}
