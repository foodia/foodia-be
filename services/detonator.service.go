package services

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/entities"
	"foodia-be/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DetonatorService struct {
	DB          *gorm.DB
	Log         *zerolog.Logger
	AuthService *AuthService
}

func NewDetonatorService(ctx context.Context, db *gorm.DB) *DetonatorService {
	logger := ctx.Value(enums.LoggerCtxKey).(*zerolog.Logger)

	return &DetonatorService{
		DB:          db,
		Log:         logger,
		AuthService: NewAuthService(ctx, db),
	}
}

func (service DetonatorService) DetonatorRegistration(c *fiber.Ctx, input dto.DetonatorRegistration) (*entities.Detonator, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	// store self photo file to the storage
	selfPhoto := "detonator/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
	if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", selfPhoto)); err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// store ktp photo file to the storage
	ktpPhoto := "detonator/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
	if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", ktpPhoto)); err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// insert ouath
	ouath := entities.Oauth{
		Fullname: input.Fullname,
		Email:    input.Email,
		Phone:    input.Phone,
		UserId:   common.GenerateUUID(),
		Password: string(password),
		Role:     "detonator",
	}

	if err := tx.Create(&ouath).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// insert detonator
	detonator := entities.Detonator{
		UserId:    ouath.ID,
		KTPNumber: input.KTPNumber,
		KTPPhoto:  ktpPhoto,
		SelfPhoto: selfPhoto,
		Status:    "waiting",
	}

	if err := tx.Create(&detonator).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// send OTP
	OTP := dto.OTPRequest{
		Email: ouath.Email,
	}

	if err := service.AuthService.SendOTP(OTP); err != nil {
		service.Log.Error().Msg(err.StatusCode.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.StatusCode.Error(),
		}
	}

	if err := tx.Commit().Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	return &detonator, nil
}

func (service DetonatorService) GetAllDetonator(pagination *common.Pagination) ([]entities.Detonator, *dto.ApiError) {
	var detonators []entities.Detonator

	query := service.DB.
		Preload("Oauth").
		Order("created_at desc").
		Find(&detonators)

	if err := query.Scopes(common.Paginate(query, entities.Detonator{}, pagination)).Find(&detonators); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return detonators, nil
}

func (service DetonatorService) GetByID(id string) (*entities.Detonator, *dto.ApiError) {
	var detonator entities.Detonator

	if err := service.DB.
		Preload("Oauth").
		Where("id", id).
		First(&detonator); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return &detonator, nil
}

func (service DetonatorService) DetonatorApproval(id string, input dto.DetonatorApproval) (*entities.Detonator, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var detonator entities.Detonator
	if err := tx.First(&detonator, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	update := entities.Detonator{
		Status: input.Status,
		Note:   input.Note,
	}

	if err := tx.Model(&detonator).Updates(&update).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	if err := tx.Commit().Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	return &update, nil
}

func (service DetonatorService) DetonatorUpdate(c *fiber.Ctx, id string, input dto.DetonatorUpdate) (*entities.Detonator, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var selfPhoto string
	// store self photo file to the storage
	if input.SelfPhoto != nil {
		selfPhoto = "detonator/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
		if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", selfPhoto)); err != nil {
			service.Log.Error().Msg(err.Error())
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrInternalServerError,
				Message:    err.Error(),
			}
		}
	}

	var ktpPhoto string
	// store ktp photo file to the storage
	if input.KTPPhoto != nil {
		ktpPhoto = "detonator/" + common.GenerateUUID() + filepath.Ext(input.KTPPhoto.Filename)
		if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", ktpPhoto)); err != nil {
			service.Log.Error().Msg(err.Error())
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrInternalServerError,
				Message:    err.Error(),
			}
		}
	}

	var detonator entities.Detonator
	if err := tx.First(&detonator, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    errors.New("detonator not found").Error(),
		}
	}

	detonatorUpdate := entities.Detonator{
		KTPNumber: input.KTPNumber,
		KTPPhoto:  ktpPhoto,
		SelfPhoto: selfPhoto,
	}

	if err := tx.Model(&detonator).Updates(&detonatorUpdate).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	var ouath entities.Oauth
	if err := tx.First(&ouath, "id", detonator.UserId).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    errors.New("ouath not found").Error(),
		}
	}

	ouathUpdate := entities.Oauth{
		Fullname: input.Fullname,
		Email:    input.Email,
		Phone:    input.Phone,
	}

	if err := tx.Model(&ouath).Updates(&ouathUpdate).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	if err := tx.Commit().Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	detonator.Oauth = &ouath

	return &detonator, nil
}
