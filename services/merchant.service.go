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

type MerchantService struct {
	DB          *gorm.DB
	Log         *zerolog.Logger
	AuthService *AuthService
}

func NewMerchantService(ctx context.Context, db *gorm.DB) *MerchantService {
	logger := ctx.Value(enums.LoggerCtxKey).(*zerolog.Logger)

	return &MerchantService{
		DB:          db,
		Log:         logger,
		AuthService: NewAuthService(ctx, db),
	}
}

func (service MerchantService) MerchantRegistration(c *fiber.Ctx, input dto.MerchantRegistration) (*entities.Merchant, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	// store self photo file to the storage
	selfPhoto := "merchant/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
	if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", selfPhoto)); err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// store ktp photo file to the storage
	ktpPhoto := "merchant/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
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
		Role:     "merchant",
	}

	if err := tx.Create(&ouath).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	// insert merchant
	merchant := entities.Merchant{
		UserId:      ouath.ID,
		KTPNumber:   input.KTPNumber,
		KTPPhoto:    ktpPhoto,
		SelfPhoto:   selfPhoto,
		Status:      "waiting",
		Province:    input.Province,
		City:        input.City,
		SubDistrict: input.SubDistrict,
		PostalCode:  input.PostalCode,
		Address:     input.Address,
		Latitude:    input.Latitude,
		Longitude:   input.Latitude,
		NoLinkAja:   input.NoLinkAja,
	}

	if err := tx.Create(&merchant).Error; err != nil {
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

	return &merchant, nil
}

func (service MerchantService) GetAllMerchant(pagination *common.Pagination) ([]entities.Merchant, *dto.ApiError) {
	var merchants []entities.Merchant

	query := service.DB.
		Preload("Oauth").
		Preload("MerchantProduct").
		Order("created_at desc").
		Find(&merchants)

	if err := query.Scopes(common.Paginate(query, entities.Merchant{}, pagination)).Find(&merchants); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return merchants, nil
}

func (service MerchantService) GetByID(id string) (*entities.Merchant, *dto.ApiError) {
	var merchant entities.Merchant

	if err := service.DB.
		Preload("Oauth").
		Preload("MerchantProduct").
		Where("id", id).
		First(&merchant); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return &merchant, nil
}

func (service MerchantService) MerchantApproval(id string, input dto.MerchantApproval) (*entities.Merchant, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var merchant entities.Merchant
	if err := tx.First(&merchant, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	update := entities.Merchant{
		Status: input.Status,
		Note:   input.Note,
	}

	if err := tx.Model(&merchant).Updates(&update).Error; err != nil {
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

func (service MerchantService) MerchantUpdate(c *fiber.Ctx, id string, input dto.MerchantUpdate) (*entities.Merchant, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var selfPhoto string
	// store self photo file to the storage
	if input.SelfPhoto != nil {
		selfPhoto = "merchant/" + common.GenerateUUID() + filepath.Ext(input.SelfPhoto.Filename)
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
		ktpPhoto = "merchant/" + common.GenerateUUID() + filepath.Ext(input.KTPPhoto.Filename)
		if err := c.SaveFile(input.SelfPhoto, fmt.Sprintf("./storage/%s", ktpPhoto)); err != nil {
			service.Log.Error().Msg(err.Error())
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrInternalServerError,
				Message:    err.Error(),
			}
		}
	}

	var merchant entities.Merchant
	if err := tx.First(&merchant, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    errors.New("merchant not found").Error(),
		}
	}

	merchantUpdate := entities.Merchant{
		KTPNumber:   input.KTPNumber,
		KTPPhoto:    ktpPhoto,
		SelfPhoto:   selfPhoto,
		Province:    input.Province,
		City:        input.City,
		SubDistrict: input.SubDistrict,
		PostalCode:  input.PostalCode,
		Address:     input.Address,
		Latitude:    input.Latitude,
		Longitude:   input.Latitude,
		NoLinkAja:   input.NoLinkAja,
	}

	if err := tx.Model(&merchant).Updates(&merchantUpdate).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	var ouath entities.Oauth
	if err := tx.First(&ouath, "id", merchant.UserId).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    errors.New("oauth not found").Error(),
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

	merchant.Oauth = &ouath

	return &merchant, nil
}
