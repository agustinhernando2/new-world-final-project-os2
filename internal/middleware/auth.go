package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/cmd/tools"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	AuthService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}

func (c *AuthMiddleware) UserMiddleware(ctx *fiber.Ctx) error {
	// Get token from cookie
	tokenStr := ctx.Cookies("Authorization")

	if tokenStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   "No token provided",
		})

	}

	// Extract the JWT token from the cookie

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		hmacSampleSecret := []byte(tools.GetEnvValue("USER_SECRET_KEY", "secret_key"))
		return hmacSampleSecret, nil
	})
	if err != nil {
		log.Fatal(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})

	}

	// Check expiry of the token
	if claims["ttl"].(int64) < time.Now().Unix() {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   "Token expired",
		})

	}

	// Extract the user from the token
	userID := claims["userID"].(uint)
	user, err := c.AuthService.GetUserFromId(userID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})

	}

	// Set the current user in the context
	ctx.Locals("user", user)

	// Continue
	return ctx.Next()
}

func (c *AuthMiddleware) AdminMiddleware(ctx *fiber.Ctx) error {
	// Get token from cookie
	tokenStr := ctx.Cookies("Authorization")

	if tokenStr == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   "No token provided",
		})

	}

	// Extract the JWT token from the cookie

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		hmacSampleSecret := []byte(tools.GetEnvValue("ADMIN_SECRET_KEY", "secret_key"))
		return hmacSampleSecret, nil
	})
	if err != nil {
		log.Fatal(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   "Invalid token",
		})

	}

	// Check expiry of the token
	if claims["ttl"].(int64) < time.Now().Unix() {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   "Token expired",
		})

	}

	// Extract the user from the token
	userID := claims["userID"].(uint)
	user, err := c.AuthService.GetUserFromId(userID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized access",
			"error":   err.Error(),
		})

	}

	// Set the current user in the context
	ctx.Locals("user", user)

	// Continue
	return ctx.Next()
}
