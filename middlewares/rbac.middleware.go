package middlewares

import (
	"strings"
	"time"

	"foodia-be/common"
	"foodia-be/dto"
	"foodia-be/enums"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RBACMiddleware struct {
	Secret   string
	Duration time.Duration
}

func NewRBACMiddleware(secret string, duration time.Duration) *RBACMiddleware {
	return &RBACMiddleware{
		Secret:   secret,
		Duration: duration,
	}
}

// allowRole is a middleware function that validates and refreshes JWT tokens.
// It checks the "Authorization" header in the request, validates the JWT token,
// and refreshes the token if it is about to expire.
// It returns a Fiber handler function that can be used as middleware.
func (m RBACMiddleware) allowRole(allowed []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the "Authorization" header from the request
		authorization := c.GetReqHeaders()[fiber.HeaderAuthorization]

		// Split the authorization header into fields
		authFields := strings.Fields(authorization)
		if len(authFields) < 2 {
			// Return unauthorized response if the authorization header is missing or incomplete
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ApiResponse{
				Code:    fiber.ErrUnauthorized.Code,
				Message: fiber.ErrUnauthorized.Message,
				Error:   jwt.ErrTokenSignatureInvalid.Error(),
			})
		}

		if authFields[0] != "Bearer" {
			// Return unauthorized response if the authorization type is not "Bearer"
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ApiResponse{
				Code:    fiber.ErrUnauthorized.Code,
				Message: fiber.ErrUnauthorized.Message,
				Error:   jwt.ErrTokenSignatureInvalid.Error(),
			})
		}

		// Unmarshal and validate the JWT claims using the provided secret
		claims, err := common.UnmarshalClaims(m.Secret, authFields[1])
		if err != nil {
			// Return unauthorized response if the token is invalid or expired
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ApiResponse{
				Code:    fiber.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		// Refresh the token if it is about to expire
		if time.Until(claims.ExpiresAt.Time) <= 10*time.Minute {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(m.Duration))
			token, err := common.MarshalClaims(m.Secret, claims)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(dto.ApiResponse{
					Code:    fiber.StatusUnauthorized,
					Message: err.Error(),
				})
			}
			authorization = token.GetBearer()
		}

		// Set the updated authorization header to include the refreshed token
		c.Set(fiber.HeaderAuthorization, authorization)

		// Store the validated JWT claims in the context locals for future use
		c.Locals("session", claims)
		// Return forbidden response if the user's role does not match the required role
		for _, allow := range allowed {
			if claims.Session == common.GenerateSHA256(m.Secret, allow) {
				// Set the updated authorization header and store the claims in the context locals
				return c.Next()
			}
		}

		// Return StatusUnauthorized if role user not in list allowd
		return c.Status(fiber.StatusForbidden).JSON(dto.ApiResponse{
			Code:    fiber.ErrUnauthorized.Code,
			Message: fiber.ErrUnauthorized.Message,
			Error:   enums.ErrAccessForbidden.Error(),
		})
	}
}

func (m RBACMiddleware) AllowSuperAdmin() fiber.Handler {
	return m.allowRole([]string{"superadmin"})
}

func (m RBACMiddleware) AllowMerchant() fiber.Handler {
	return m.allowRole([]string{"merchant"})
}

func (m RBACMiddleware) AllowDetonator() fiber.Handler {
	return m.allowRole([]string{"detonator"})
}

func (m RBACMiddleware) AllowAll() fiber.Handler {
	return m.allowRole([]string{"superadmin", "detonator", "merchant"})
}
