package services

import (
	"context"
	"embed"
	"strings"
	"time"

	"foodia-be/common"
	"foodia-be/configs"
	"foodia-be/dto"
	"foodia-be/entities"
	"foodia-be/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/wneessen/go-mail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB     *gorm.DB
	Log    *zerolog.Logger
	Mail   *MailService
	Config *configs.EnvConfig
}

func NewAuthService(ctx context.Context, db *gorm.DB) *AuthService {
	logger := ctx.Value(enums.LoggerCtxKey).(*zerolog.Logger)
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	smtp := ctx.Value(enums.SmtpCtxKey).(*mail.Client)
	template := ctx.Value(enums.TemplateCtxKey).(embed.FS)

	return &AuthService{
		DB:     db,
		Log:    logger,
		Config: config,
		Mail:   NewMailService(config.SmtpSender, smtp, template),
	}
}

func (service AuthService) BasicAuthentication(input dto.AuthRequest) (*dto.AuthResponse, *dto.ApiError) {
	var oauth entities.Oauth
	if err := service.DB.First(&oauth, "email", input.Email).Error; err != nil {
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oauth.Password), []byte(input.Password)); err != nil {
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrBadRequest,
			Message:    err.Error(),
		}
	}

	claims := &dto.JWTClaims{
		UserId:  oauth.ID,
		Session: common.GenerateSHA256(service.Config.JWTSecret, oauth.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(service.Config.JWTExpirationDuration),
			},
		},
	}

	token, err := common.MarshalClaims(service.Config.JWTSecret, claims)
	if err != nil {
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrBadRequest,
			Message:    err.Error(),
		}
	}

	var status string
	var note string
	var userId int

	if oauth.Role == "merchant" {
		var merchant entities.Merchant
		if err := service.DB.Find(&merchant, "user_id", oauth.ID).Error; err != nil {
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrNotFound,
				Message:    err.Error(),
			}
		}

		status = merchant.Status
		note = merchant.Note
		userId = merchant.ID
	}

	if oauth.Role == "detonator" {
		var detonator entities.Detonator
		if err := service.DB.Find(&detonator, "user_id", oauth.ID).Error; err != nil {
			return nil, &dto.ApiError{
				StatusCode: fiber.ErrNotFound,
				Message:    err.Error(),
			}
		}

		status = detonator.Status
		note = detonator.Note
		userId = detonator.ID
	}

	oauthResponse := dto.AuthResponse{
		Fullname: oauth.Fullname,
		Phone:    oauth.Phone,
		Email:    oauth.Email,
		Role:     oauth.Role,
		Token:    token.TokenString,
		User: &dto.User{
			Status: status,
			Note:   note,
			ID:     userId,
		},
	}

	return &oauthResponse, nil
}

func (service AuthService) SendOTP(input dto.OTPRequest) *dto.ApiError {
	tx := service.DB.Begin()
	defer tx.Rollback()

	OTP := entities.OauthOTP{
		Email:     input.Email,
		OTPCode:   common.GenerateOTP(),
		ExpiredAt: input.ExpiredAt,
	}

	if err := tx.Create(&OTP).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	go func() {
		if err := service.Mail.SendOTP(input.Email, OTP.OTPCode); err != nil {
			service.Log.Error().Msg(err.Error())
		}
	}()

	if err := tx.Commit().Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return &dto.ApiError{
			StatusCode: fiber.ErrInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}

func (service AuthService) ValidateOTP(input dto.ValidateOTPRequest) (*dto.AuthResponse, *dto.ApiError) {
	var oauth entities.Oauth
	if err := service.DB.First(&oauth, "email", input.Email).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	var otp entities.OauthOTP
	if err := service.DB.First(&otp, "email", input.Email).Error; err != nil {
		service.Log.Error().Msg(err.Error())
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrNotFound,
			Message:    err.Error(),
		}
	}

	if strings.Compare(input.Code, otp.OTPCode) != 0 {
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrBadRequest,
			Message:    "OTP doesn't match, please recheck your OTP code",
		}
	}

	claims := &dto.JWTClaims{
		UserId:  oauth.ID,
		Session: common.GenerateSHA256(service.Config.JWTSecret, oauth.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(service.Config.JWTExpirationDuration),
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	token, err := common.MarshalClaims(service.Config.JWTSecret, claims)
	if err != nil {
		return nil, &dto.ApiError{
			StatusCode: fiber.ErrBadRequest,
			Message:    err.Error(),
		}
	}

	oauthResponse := dto.AuthResponse{
		Fullname: oauth.Fullname,
		Phone:    oauth.Phone,
		Email:    oauth.Email,
		Role:     oauth.Role,
		Token:    token.TokenString,
	}

	return &oauthResponse, nil
}
