package routers

import (
	"foodia-be/configs"
	"foodia-be/controllers"
	"foodia-be/enums"
	"foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseMerchantRouter(ctx context.Context, r fiber.Router) {
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	auth := middlewares.NewRBACMiddleware(config.JWTSecret, config.JWTExpirationDuration)
	ctrl := controllers.NewMerchantController(ctx)

	merchantGroup := r.Group("/merchant")
	merchantGroup.Post("/registration", ctrl.MerchantRegistration)
	merchantGroup.Get("/filter", auth.AllowAll(), ctrl.GetAllMerchant)
	merchantGroup.Get("/fetch/:id", auth.AllowAll(), ctrl.GetByID)
	merchantGroup.Put("/approval/:id", auth.AllowAll(), ctrl.MerchantApproval)
	merchantGroup.Put("/update/:id", auth.AllowAll(), ctrl.MerchantUpdate)
}
