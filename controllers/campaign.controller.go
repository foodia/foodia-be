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

type CampaignController struct {
	CampaignService *services.CampaignService
}

func NewCampaignController(ctx context.Context) *CampaignController {
	db := ctx.Value(enums.GormCtxKey).(*gorm.DB)

	return &CampaignController{
		CampaignService: services.NewCampaignService(ctx, db),
	}
}

func (ctrl CampaignController) CampaignCreate(c *fiber.Ctx) error {
	var req dto.CampaignRequest
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

	campaign, fail := ctrl.CampaignService.Create(req)
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
		Body:    campaign,
	})
}

func (ctrl CampaignController) GetAll(c *fiber.Ctx) error {
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

	campaigns, fail := ctrl.CampaignService.GetAll(c, &pagination)
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
		Body:    campaigns,
		Meta:    pagination,
	})
}

func (ctrl CampaignController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	campaign, fail := ctrl.CampaignService.GetByID(id)
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
		Body:    campaign,
	})
}

func (ctrl CampaignController) CampaignUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.CampaignRequest
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

	campaign, fail := ctrl.CampaignService.Update(id, req)
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
		Body:    campaign,
	})
}
