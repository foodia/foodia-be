package controllers

import (
	"context"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/enums"
	"foodia-be/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(ctx context.Context) *AuthController {
	db := ctx.Value(enums.GormCtxKey).(*gorm.DB)

	return &AuthController{
		AuthService: services.NewAuthService(ctx, db),
	}
}

func (ctrl AuthController) BasicAuthentication(c *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   err.Error(),
		})
	}

	if err := common.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	oauth, fail := ctrl.AuthService.BasicAuthentication(req)
	if fail != nil {
		return c.Status(fail.StatusCode.Code).JSON(dto.ApiResponse{
			Code:    fail.StatusCode.Code,
			Message: fail.StatusCode.Message,
			Error:   fail.Message,
		})
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Successfuly",
		Body:    oauth,
	})
}

func (ctrl AuthController) ValidateOTP(c *fiber.Ctx) error {
	var req dto.ValidateOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   err.Error(),
		})
	}

	if err := common.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	oauth, fail := ctrl.AuthService.ValidateOTP(req)
	if fail != nil {
		return c.Status(fail.StatusCode.Code).JSON(dto.ApiResponse{
			Code:    fail.StatusCode.Code,
			Message: fail.StatusCode.Message,
			Error:   fail.Message,
		})
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Successfuly",
		Body:    oauth,
	})
}
