package controllers

import (
	"context"
	"strconv"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/enums"
	"foodia-be/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MerchantProductController struct {
	MerchantProductService *services.MerchantProductService
}

func NewMerchantProductController(ctx context.Context) *MerchantProductController {
	db := ctx.Value(enums.GormCtxKey).(*gorm.DB)

	return &MerchantProductController{
		MerchantProductService: services.NewMerchantProductService(ctx, db),
	}
}

func (ctrl MerchantProductController) MerchantProductCreate(c *fiber.Ctx) error {
	var req dto.MerchantProductRequest
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

	merchantProduct, fail := ctrl.MerchantProductService.Create(req)
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
		Body:    merchantProduct,
	})
}

func (ctrl MerchantProductController) GetByMerchant(c *fiber.Ctx) error {
	merchantId := c.Query("merchant_id")

	if merchantId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.StatusBadRequest,
			Message: fiber.ErrBadRequest.Message,
			Error:   "merchantId cannot be null!",
		})
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = common.DefaultPage
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil {
		perPage = common.DefaultPerPage
	}

	pagination := common.Pagination{
		Page:    page,
		PerPage: perPage,
	}

	merchantProducts, fail := ctrl.MerchantProductService.GetByMerchant(merchantId, &pagination)
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
		Body:    merchantProducts,
		Meta:    pagination,
	})
}

func (ctrl MerchantProductController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	merchantProduct, fail := ctrl.MerchantProductService.GetByID(id)
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
		Body:    merchantProduct,
	})
}

func (ctrl MerchantProductController) MerchantProductUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.MerchantProductRequest
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

	merchantProduct, fail := ctrl.MerchantProductService.Update(id, req)
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
		Body:    merchantProduct,
	})
}
