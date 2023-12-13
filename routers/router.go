package routers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func UseRouter(ctx context.Context, r fiber.Router) {
	prefix := r.Group("/api/v1")

	UseAuthRouter(ctx, prefix)
	UseDetonatorRouter(ctx, prefix)
	UseMerchantRouter(ctx, prefix)
	UseMediaRouter(ctx, prefix)
	UseMerchantProductRouter(ctx, prefix)
	UseCampaignRouter(ctx, prefix)
}
