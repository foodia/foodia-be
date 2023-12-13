package services

import (
	"context"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/entities"
	"foodia-be/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MerchantProductService struct {
	DB  *gorm.DB
	Log *zerolog.Logger
}

func NewMerchantProductService(ctx context.Context, db *gorm.DB) *MerchantProductService {
	logger := ctx.Value(enums.LoggerCtxKey).(*zerolog.Logger)

	return &MerchantProductService{
		DB:  db,
		Log: logger,
	}
}

func (service MerchantProductService) Create(input dto.MerchantProductRequest) (*dto.MerchantProductRequest, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	merchantProduct := entities.MerchantProduct{
		MerchantID:  input.MerchantID,
		Name:        input.Name,
		Description: input.Description,
		Price:       decimal.NewFromFloat(input.Price),
		QTY:         input.QTY,
	}

	if err := tx.Create(&merchantProduct).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	var productImages []entities.MerchantProductImage

	for _, image := range input.Images {
		productImages = append(productImages, entities.MerchantProductImage{
			MerchantProductID: merchantProduct.ID,
			ImageURL:          image.ImageURL,
		})
	}

	if err := tx.Create(&productImages).Error; err != nil {
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

	return &input, nil
}

func (service MerchantProductService) GetByMerchant(merchantId string, pagination *common.Pagination) ([]entities.MerchantProduct, *dto.ApiError) {
	var merchantProducts []entities.MerchantProduct

	query := service.DB.
		Preload("MerchantProductImage").
		Order("created_at desc").
		Where("merchant_id", merchantId).
		Find(&merchantProducts)

	if err := query.Scopes(common.Paginate(query, entities.MerchantProduct{}, pagination)).Find(&merchantProducts); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return merchantProducts, nil
}

func (service MerchantProductService) GetByID(id string) (*entities.MerchantProduct, *dto.ApiError) {
	var merchantProduct entities.MerchantProduct

	if err := service.DB.
		Preload("MerchantProductImage").
		Where("id", id).
		First(&merchantProduct); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return &merchantProduct, nil
}

func (service MerchantProductService) Update(id string, input dto.MerchantProductRequest) (*dto.MerchantProductRequest, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var merchantProduct entities.MerchantProduct

	if err := tx.First(&merchantProduct, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	update := entities.MerchantProduct{
		MerchantID:  input.MerchantID,
		Name:        input.Name,
		Description: input.Description,
		Price:       decimal.NewFromFloat(input.Price),
		QTY:         input.QTY,
	}

	if err := tx.Model(&merchantProduct).Updates(&update).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	var productImages []entities.MerchantProductImage

	// check length images
	if len(input.Images) > 0 {
		// delete images
		if err := tx.Where("merchant_product_id = ?", merchantProduct.ID).Delete(&entities.MerchantProductImage{}).Error; err != nil {
			service.Log.Error().Msg(err.Error())
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrInternalServerError,
				Message:    err.Error(),
			}
		}

		// insert images
		for _, image := range input.Images {
			productImages = append(productImages, entities.MerchantProductImage{
				MerchantProductID: merchantProduct.ID,
				ImageURL:          image.ImageURL,
			})
		}

		if err := tx.Create(&productImages).Error; err != nil {
			service.Log.Error().Msg(err.Error())
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrInternalServerError,
				Message:    err.Error(),
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	return &input, nil
}
