package controllers

import (
	"context"
	"fmt"
	"path/filepath"

	"foodia-be/common"
	"foodia-be/dto"

	"github.com/gofiber/fiber/v2"
)

type MediaController struct {
}

func NewMediaController(ctx context.Context) *MediaController {
	return &MediaController{}
}

func (ctrl MediaController) MediaUpload(c *fiber.Ctx) error {
	var req dto.MediaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   err.Error(),
		})
	}

	// Populate files from form data to request variable
	paramFileKey := "file"

	file, err := c.FormFile(paramFileKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnprocessableEntity.Code,
			Message: fiber.ErrUnprocessableEntity.Message,
			Error:   fmt.Sprintf("error retrieving file from %s parameter: %v", paramFileKey, err),
		})
	}

	req.File = file

	if err := common.ValidateRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	// store self photo file to the storage
	fileUrl := req.Destination + "/" + common.GenerateUUID() + filepath.Ext(req.File.Filename)
	if err := c.SaveFile(req.File, fmt.Sprintf("./storage/%s", fileUrl)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ApiResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: fiber.ErrBadRequest.Message,
			Error:   err,
		})
	}

	return c.JSON(dto.ApiResponse{
		Code:    fiber.StatusOK,
		Message: "Successfuly",
		Body: dto.MediaResponse{
			Destination: req.Destination,
			FileUrl:     fileUrl,
		},
	})
}
