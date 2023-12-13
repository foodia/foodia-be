package routers

import (
	"foodia-be/configs"
	"foodia-be/controllers"
	"foodia-be/enums"
	"foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseCampaignRouter(ctx context.Context, r fiber.Router) {
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	auth := middlewares.NewRBACMiddleware(config.JWTSecret, config.JWTExpirationDuration)
	ctrl := controllers.NewCampaignController(ctx)

	campaignGroup := r.Group("/campaign")
	campaignGroup.Post("/create", auth.AllowAll(), ctrl.CampaignCreate)
	campaignGroup.Get("/filter", ctrl.GetAll)
	campaignGroup.Put("/update/:id", auth.AllowAll(), ctrl.CampaignUpdate)
	campaignGroup.Get("/fetch/:id", auth.AllowAll(), ctrl.GetByID)
}
