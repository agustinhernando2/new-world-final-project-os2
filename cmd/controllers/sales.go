package controllers

import (
	"strconv"

	// "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/sales"
	"github.com/gofiber/fiber/v2"
)

type SaleController struct {
	service sales.SaleService
}

// NewSaleController creates a new SaleController
func NewSaleController(service sales.SaleService) *SaleController {
	return &SaleController{service}
}

// GetAllSales handles GET /v1/sales
// @Summary Get all sales
// @Description Get all sales
// @Tags sales
// @Accept json
// @Produce json
// @Success 200 {array} models.Sale
// @Router /v1/sales [get]
// func (c *SaleController) GetAllSales(ctx *fiber.Ctx) error {
// 	sales, err := c.service.GetSales()
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return ctx.JSON(sales)
// }

// GetSaleById handles GET /v1/sales/:id
// @Summary Get a sale by ID
// @Description Get a sale by ID
// @Tags sales
// @Accept json
// @Produce json
// @Param id path int true "Sale ID"
// @Success 200 {object} models.Sale
// @Failure 404 {object} fiber.Map
// @Router /sales/{id} [get]
func (c *SaleController) GetSaleById(ctx *fiber.Ctx) error {

	m := ctx.Queries()
	//parse int
	idStr := m["id"]
	// Verificar si el ID está presente
	if idStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID query parameter is required",
		})
	}
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sale ID format",
		})
	}

	sale, err := c.service.GetSaleID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sale not found"})
	}
	return ctx.Status(200).JSON(sale)
}

// CreateSale handles POST /v1/sales
// @Summary Create a new sale
// @Description Create a new sale
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body models.Sale true "Sale"
// @Success 201
// @Failure 400 {object} fiber.Map
// @Router /v1/sales [post]
func (c *SaleController) CreateSale(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Sale created successfully"})
}
