package controllers

import (
	"context"
	"fmt"
	"strconv"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/enums"
	"foodia-be/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MerchantController struct {
	MerchantService *services.MerchantService
}

func NewMerchantController(ctx context.Context) *MerchantController {
	db := ctx.Value(enums.GormCtxKey).(*gorm.DB)

	return &MerchantController{
		MerchantService: services.NewMerchantService(ctx, db),
	}
}

func (ctrl MerchantController) MerchantRegistration(c *fiber.Ctx) error {
	var req dto.MerchantRegistration
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   err.Error(),
		})
	}

	// Populate files from form data to request variable
	paramSelfPhotoKey := "self_photo"

	selfPhotoFile, err := c.FormFile(paramSelfPhotoKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   fmt.Sprintf("error retrieving file from %s parameter: %v", paramSelfPhotoKey, err),
		})
	}

	// Populate files from form data to request variable
	paramKTPPhotoKey := "self_photo"

	KTPPhotoFile, err := c.FormFile(paramKTPPhotoKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   fmt.Sprintf("error retrieving file from %s parameter: %v", paramKTPPhotoKey, err),
		})
	}

	req.SelfPhoto = selfPhotoFile
	req.KTPPhoto = KTPPhotoFile

	if err := common.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	merchant, fail := ctrl.MerchantService.MerchantRegistration(c, req)
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
		Body:    merchant,
	})
}

func (ctrl MerchantController) GetAllMerchant(c *fiber.Ctx) error {

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

	merchants, fail := ctrl.MerchantService.GetAllMerchant(&pagination)
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
		Body:    merchants,
		Meta:    pagination,
	})
}

func (ctrl MerchantController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	merchant, fail := ctrl.MerchantService.GetByID(id)
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
		Body:    merchant,
	})
}

func (ctrl MerchantController) MerchantApproval(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.MerchantApproval
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

	_, fail := ctrl.MerchantService.MerchantApproval(id, req)
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
	})
}

func (ctrl MerchantController) MerchantUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.MerchantUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   err.Error(),
		})
	}

	// Populate files from form data to request variable
	paramSelfPhotoKey := "self_photo"

	selfPhotoFile, err := c.FormFile(paramSelfPhotoKey)
	if selfPhotoFile != nil {
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
				Code:    fiber.ErrUnprocessableEntity.Code,
				Message: fiber.ErrUnprocessableEntity.Message,
				Error:   fmt.Sprintf("error retrieving file from %s parameter: %v", paramSelfPhotoKey, err),
			})
		}

		req.SelfPhoto = selfPhotoFile
	}

	// Populate files from form data to request variable
	paramKTPPhotoKey := "ktp_photo"

	KTPPhotoFile, err := c.FormFile(paramKTPPhotoKey)
	if KTPPhotoFile != nil {
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
				Code:    fiber.ErrUnprocessableEntity.Code,
				Message: fiber.ErrUnprocessableEntity.Message,
				Error:   fmt.Sprintf("error retrieving file from %s parameter: %v", paramKTPPhotoKey, err),
			})
		}

		req.KTPPhoto = KTPPhotoFile
	}

	if err := common.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	merchant, fail := ctrl.MerchantService.MerchantUpdate(c, id, req)
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
		Body:    merchant,
	})
}
