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

func (c *AdminController) GetDashboard(ctx *fiber.Ctx) error {
	// Implementar lógica para obtener el estado de todas las órdenes
	return ctx.SendStatus(fiber.StatusNotImplemented)
}

func (c *AdminController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	// Implementar lógica para actualizar el estado de una orden específica
	return ctx.SendStatus(fiber.StatusNotImplemented)
}
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
