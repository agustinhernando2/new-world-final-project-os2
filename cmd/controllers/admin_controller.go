// internal/controllers/admin_controller.go
package controllers

import (
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	AdminService services.AdminService
}

func NewAdminController(adminService services.AdminService) *AdminController {
	return &AdminController{
		AdminService: adminService,
	}
}

// GetDashboard retrieves the dashboard information.
//
//	@Summary		Retrieve the dashboard information
//	@Description	Retrieve the dashboard information about orders and offers. The user must be authenticated as an admin.
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Dashboard information"
//	@Failure		500	{object}	fiber.Map	"Failed to get dashboard"
//	@Router			/admin/dashboard [get]
//	@Security		ApiKeyAuth
func (c *AdminController) GetDashboard(ctx *fiber.Ctx) error {
	// Call service
	offers, orders, err := c.AdminService.GetDashboard()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dashboard.",
		})
	}

	// success 200
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orders,
		"offers": offers,
	})
}

// UpdateOrderStatus updates the status of an order.
//
//	@Summary		Update the status of an order
//	@Description	Update the status of an order. The new status must be provided in the request body. The user must be authenticated as an admin.
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Order ID"
//	@Param			status	body		string	true	"New status"
//	@Success		200		{object}	string	"Order status updated successfully"
//	@Failure		400		{object}	string	"Cannot parse input data"
//	@Failure		500		{object}	string	"Failed to update order status"
//	@Router			/admin/order/{id} [patch]
//	@Security		ApiKeyAuth
func (c *AdminController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	// Receive a JSON with the order ID and the new status
	// Get orderID from path
	orderID, err := ctx.ParamsInt("id")
	if err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data.",
			"error":   err.Error(),
		})
	}
	//get status from json body
	var status struct {
		Status string `json:"status"`
	}
	if err := ctx.BodyParser(&status); err != nil {
		// error 400 Bad Request
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse input data.",
			"error":   err.Error(),
		})
	}

	// Call service
	err = c.AdminService.UpdateOrderStatus(uint(orderID), status.Status)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status.",
		})
	}

	// success 200
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}

// UpdateItemsAvailables updates the available items in the storage.
//
//	@Summary		Update the available items
//	@Description	Update the available items in the Cpp Server Storage. The user must be authenticated as an admin.
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string	"Items updated successfully"
//	@Failure		500	{object}	string	"Failed to update"
//	@Router			/admin/updatesupplies [post]
//	@Security		ApiKeyAuth
func (c *AdminController) UpdateItemsAvailables(ctx *fiber.Ctx) error {
	res, err := c.AdminService.UpdateStorage()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update.",
		})
	}

	// success 200
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Items updated successfully",
		"items":   res,
	})
}
