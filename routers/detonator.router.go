package routers

import (
	"foodia-be/configs"
	"foodia-be/controllers"
	"foodia-be/enums"
	"foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseDetonatorRouter(ctx context.Context, r fiber.Router) {
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	auth := middlewares.NewRBACMiddleware(config.JWTSecret, config.JWTExpirationDuration)
	ctrl := controllers.NewDetonatorController(ctx)

	detonatorGroup := r.Group("/detonator")
	detonatorGroup.Post("/registration", ctrl.DetonatorRegistration)
	detonatorGroup.Get("/filter", auth.AllowAll(), ctrl.GetAllDetonator)
	detonatorGroup.Get("/fetch/:id", auth.AllowAll(), ctrl.GetByID)
	detonatorGroup.Put("/approval/:id", auth.AllowAll(), ctrl.DetonatorApproval)
	detonatorGroup.Put("/update/:id", auth.AllowAll(), ctrl.DetonatorUpdate)
}
