package routers

import (
	"foodia-be/configs"
	"foodia-be/controllers"
	"foodia-be/enums"
	"foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseMediaRouter(ctx context.Context, r fiber.Router) {
	config := ctx.Value(enums.ConfigCtxKey).(*configs.EnvConfig)
	auth := middlewares.NewRBACMiddleware(config.JWTSecret, config.JWTExpirationDuration)
	ctrl := controllers.NewMediaController(ctx)

	mediaGroup := r.Group("/media")
	mediaGroup.Post("/upload", auth.AllowAll(), ctrl.MediaUpload)
}
