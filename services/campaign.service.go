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

type CampaignService struct {
	DB  *gorm.DB
	Log *zerolog.Logger
}

func NewCampaignService(ctx context.Context, db *gorm.DB) *CampaignService {
	logger := ctx.Value(enums.LoggerCtxKey).(*zerolog.Logger)

	return &CampaignService{
		DB:  db,
		Log: logger,
	}
}

func (service CampaignService) Create(input dto.CampaignRequest) (*dto.CampaignRequest, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	campaign := entities.Campaign{
		DetonatorID:    input.DetonatorID,
		EventName:      input.EventName,
		EventType:      input.EventType,
		EventDate:      input.EventDate,
		EventTime:      input.EventTime,
		Description:    input.Description,
		DonationTarget: decimal.NewFromFloat(input.DonationTarget),
		Province:       input.Province,
		City:           input.City,
		SubDistrict:    input.SubDistrict,
		PostalCode:     input.PostalCode,
		Address:        input.Address,
		Latitude:       input.Latitude,
		Longitude:      input.Longitude,
		ImageURL:       input.ImageURL,
	}

	if err := tx.Create(&campaign).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	var orders []entities.Order

	for _, product := range input.Products {
		orders = append(orders, entities.Order{
			MerchantProductID: product.MerchantProductID,
			CampaignID:        campaign.ID,
		})
	}

	if err := tx.Create(&orders).Error; err != nil {
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

func (service CampaignService) GetAll(c *fiber.Ctx, pagination *common.Pagination) ([]entities.Campaign, *dto.ApiError) {
	var campaigns []entities.Campaign

	var where string
	detonatorId := c.Query("detonator_id", "")

	where = `1 = 1`

	if detonatorId != "" {
		where += ` AND detonator_id = '` + detonatorId + `' `
	}

	query := service.DB.
		Preload("Detonator").
		Preload("Detonator.Oauth").
		Order("created_at desc").
		Where(where).
		Find(&campaigns)

	if err := query.Scopes(common.Paginate(query, entities.Campaign{}, pagination)).Find(&campaigns); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return campaigns, nil
}

func (service CampaignService) GetByID(id string) (*entities.Campaign, *dto.ApiError) {
	var campaign entities.Campaign

	if err := service.DB.
		Preload("Detonator").
		Preload("Detonator.Oauth").
		Where("id", id).
		First(&campaign); err.Error != nil {
		service.Log.Error().Msg(err.Error.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error.Error(),
		}
	}

	return &campaign, nil
}

func (service CampaignService) Update(id string, input dto.CampaignRequest) (*dto.CampaignRequest, *dto.ApiError) {
	tx := service.DB.Begin()
	defer tx.Rollback()

	var campaign entities.Campaign
	if err := tx.First(&campaign, "id", id).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	update := entities.Campaign{
		DetonatorID:    input.DetonatorID,
		EventName:      input.EventName,
		EventType:      input.EventType,
		EventDate:      input.EventDate,
		EventTime:      input.EventTime,
		Description:    input.Description,
		DonationTarget: decimal.NewFromFloat(input.DonationTarget),
		Province:       input.Province,
		City:           input.City,
		SubDistrict:    input.SubDistrict,
		PostalCode:     input.PostalCode,
		Address:        input.Address,
		Latitude:       input.Latitude,
		Longitude:      input.Longitude,
		ImageURL:       input.ImageURL,
	}

	if err := tx.Model(&campaign).Updates(&update).Error; err != nil {
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
