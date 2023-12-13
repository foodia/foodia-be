package routers

import (
	"foodia-be/configs"
	"foodia-be/controllers"
	"foodia-be/enums"
	"foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseMerchantProductRouter(ctx context.Context, r fiber.Router) {
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	auth := middlewares.NewRBACMiddleware(config.JWTSecret, config.JWTExpirationDuration)
	ctrl := controllers.NewMerchantProductController(ctx)

	merchantGroup := r.Group("/merchant-product")
	merchantGroup.Post("/create", auth.AllowAll(), ctrl.MerchantProductCreate)
	merchantGroup.Get("/filter", auth.AllowAll(), ctrl.GetByMerchant)
	merchantGroup.Put("/update/:id", auth.AllowAll(), ctrl.MerchantProductUpdate)
	merchantGroup.Get("/fetch/:id", auth.AllowAll(), ctrl.GetByID)
}
