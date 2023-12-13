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

type DetonatorController struct {
	DetonatorService *services.DetonatorService
}

func NewDetonatorController(ctx context.Context) *DetonatorController {
	db := ctx.Value(enums.GormCtxKey).(*gorm.DB)

	return &DetonatorController{
		DetonatorService: services.NewDetonatorService(ctx, db),
	}
}

func (ctrl DetonatorController) DetonatorRegistration(c *fiber.Ctx) error {
	var req dto.DetonatorRegistration
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

	detonator, fail := ctrl.DetonatorService.DetonatorRegistration(c, req)
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
		Body:    detonator,
	})
}

func (ctrl DetonatorController) GetAllDetonator(c *fiber.Ctx) error {

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

	detonators, fail := ctrl.DetonatorService.GetAllDetonator(&pagination)
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
		Body:    detonators,
		Meta:    pagination,
	})
}

func (ctrl DetonatorController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	detonator, fail := ctrl.DetonatorService.GetByID(id)
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
		Body:    detonator,
	})
}

func (ctrl DetonatorController) DetonatorApproval(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.DetonatorApproval
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

	_, fail := ctrl.DetonatorService.DetonatorApproval(id, req)
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

func (ctrl DetonatorController) DetonatorUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.DetonatorUpdate
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

	detonator, fail := ctrl.DetonatorService.DetonatorUpdate(c, id, req)
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
		Body:    detonator,
	})
}
