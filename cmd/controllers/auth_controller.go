package controllers

import (
	"fmt"
	"time"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/cmd/tools"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// GetOffers retrieves a list of available offers.
//
//	@Summary		Retrieve a list of available offers
//	@Description	Retrieve a list of available offers.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.Item	"List of available offers"
//	@Failure		500	{object}	fiber.Map		"Failed to retrieve offers"
//	@Failure		401	{object}	fiber.Map		"Unauthorized"
//	@Router			/auth/offers [get]
//	@Security		ApiKeyAuth
func (c *AuthController) GetOffers(ctx *fiber.Ctx) error {
	items, err := c.AuthService.GetOffers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve offers.",
		})
	}
	return ctx.JSON(items)
}

// CheckoutOrders Buy a list of orders.
func (c *AuthController) CheckoutOrders(ctx *fiber.Ctx) error {
	var items []models.OrderItem
	// Check and parse input data
	if err := ctx.BodyParser(&items); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data",
			"error":   err.Error(),
		})
	}
	// Get userID from giber scope
	user := ctx.Locals("user")
	details, ok := user.(*models.User)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data",
			"error":   "user not found",
		})
	}
	// Call service
	order, err := c.AuthService.CheckoutOrders(details.ID, items)
	if err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error checking out orders",
			"error":   err.Error(),
		})
	}

	// success 201 Created
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Orders checked out successfully.",
		"order":   order,
	})
}

func (c *AuthController) GetOrderStatus(ctx *fiber.Ctx) error {
	// Implementar lógica para obtener el estado de una orden específica
	return ctx.SendStatus(fiber.StatusNotImplemented)
}

// RegisterUser Register a new user with the data provided in the request body.
//
//	@Summary		Register user
//	@Description	Register a new user with the data provided in the request body.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.User	true	"Data of the user to register"
//	@Success		201		{object}	models.User	"User registered successfully"
//	@Failure		400		{object}	string		"Cannot parse input data or registering user"
//	@Router			/auth/register [post]
func (c *AuthController) RegisterUser(ctx *fiber.Ctx) error {
	var newUser models.User

	// Check and parse input data
	if err := ctx.BodyParser(&newUser); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data",
			"error":   err.Error(),
		})
	}

	// Call service
	if err := c.AuthService.RegisterUser(&newUser); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error registering user",
			"error":   err.Error(),
		})
	}

	// Create and sign token
	var secret string
	if newUser.IsAdmin {
		secret = tools.GetEnvValue("ADMIN_SECRET_KEY", "asecret_key")
	} else {
		secret = tools.GetEnvValue("USER_SECRET", "usecret_key")
	}
	tokenString, err := createAndSignToken(newUser.ID, secret)
	if err != nil {
		// error 500 Internal Server Error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating token",
			"error":   err.Error(),
		})
	}

	// Send the token in a cookie
	setCookie(ctx, tokenString)

	// success 201 Created
	// return ctx.Redirect("/auth/offers")
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully.",
		"user":    newUser,
	})
}

// LoginUser Authenticate a user with the data provided in the request body.
//
//	@Summary		Authenticate user
//	@Description	Authenticate a user with the data provided in the request body.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.User	true	"Data of the user to authenticate"
//	@Success		201		{object}	models.User	"User logged successfully"
//	@Failure		400		{object}	string		"Cannot parse input data or login user"
//	@Failure		500		{object}	string		"Error creating token"
//	@Router			/auth/login [post]
func (c *AuthController) LoginUser(ctx *fiber.Ctx) error {
	var loggedUser models.User

	// Check and parse input data
	if err := ctx.BodyParser(&loggedUser); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data",
			"error":   err.Error(),
		})
	}

	// Call service
	if err := c.AuthService.AuthenticateUser(&loggedUser); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error login user",
			"error":   err.Error(),
		})
	}

	// Create and sign token
	var secret string
	if loggedUser.IsAdmin {
		secret = tools.GetEnvValue("ADMIN_SECRET_KEY", "asecret_key")
	} else {
		secret = tools.GetEnvValue("USER_SECRET_KEY", "usecret_key")
	}

	tokenString, err := createAndSignToken(loggedUser.ID, secret)
	if err != nil {
		// error 500 Internal Server Error
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating token",
			"error":   err.Error(),
		})
	}
	// Send the token in a cookie
	setCookie(ctx, tokenString)

	// success 201 Created
	// return ctx.Redirect("/auth/offers")
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User logged successfully.",
		"user":    loggedUser,
	})
}

func createAndSignToken(id uint, secret string) (string, error) {
	// secret = "xIvBy5CyQ0HDeelAlHmNhAMGBlvuHfITXOdOftjaiEg="
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(), // 100 days
	})

	hmacSampleSecret := []byte(secret)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func setCookie(ctx *fiber.Ctx, tokenString string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString, // signed token
		Expires:  time.Now().Add(time.Hour * 24 * 100),
		HTTPOnly: true,
	})
}

// GetSignupPage retrieves the signup page.
//
//	@Summary		Retrieve the signup page
//	@Description	Retrieve the signup page using HTML template.
//	@Tags			Auth
//	@Accept			html
//	@Produce		html
//	@Success		200	{object}	string	"Signup page"
//	@Router			/signup [get]
func (c *AuthController) GetSignupPage(ctx *fiber.Ctx) error {
	return ctx.Render("sessions/signup", fiber.Map{})
}

// GetLoginPage retrieves the login page.
//
//	@Summary		Retrieve the login page
//	@Description	Retrieve the login page using HTML template.
//	@Tags			Auth
//	@Accept			html
//	@Produce		html
//	@Success		200	{object}	string	"Login page"
//	@Router			/login [get]
func (c *AuthController) GetLoginPage(ctx *fiber.Ctx) error {
	// Render index
	return ctx.Render("sessions/login", fiber.Map{
		"Title": "Hello, World!",
	})
}
