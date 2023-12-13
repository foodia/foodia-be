package routers

import (
	// "time"

	"foodia-be/controllers"
	// "foodia-be/utils/enums"
	// "foodia-be/middlewares"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func UseAuthRouter(ctx context.Context, r fiber.Router) {
	// secret := ctx.Value(enums.JWTSecretCtxKey).(string)
	// expired := ctx.Value(enums.JWTExpiredCtxKey).(time.Duration)
	// auth := middlewares.NewRBACMiddleware(secret, expired)
	ctrl := controllers.NewAuthController(ctx)

	authGroup := r.Group("/auth")
	authGroup.Post("/login", ctrl.BasicAuthentication)
	authGroup.Post("/verify-otp", ctrl.ValidateOTP)
}
