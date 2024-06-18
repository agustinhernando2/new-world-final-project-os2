package sales

import (
	"fmt"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repository"
	"github.com/go-playground/validator/v10"
	// "github.com/gofiber/fiber/v2"
)

type SaleService interface {
	GetSaleID(uint) (models.Sale, error)
	// GetSale(id int) (models.Sale, error)
	CreateSale(*models.Sale) error
}

type saleService struct {
	repo repository.DbRepository
}

// NewSaleService creates a new instance of saleService
func NewSaleService(repo repository.DbRepository) SaleService {
	return &saleService{repo}
}

func (s *saleService) GetSaleID(id uint) (models.Sale, error) {
	sale, err := s.repo.GetSaleID(id)
	if err != nil {
		return models.Sale{}, err
	}
	return *sale, nil
}

func (s *saleService) CreateSale(sale *models.Sale) error {

	validate := validator.New()
	if err := validate.Struct(*sale); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Printf("Field: %s, Error: %s\n", err.StructField(), err.Tag())
			}
		}
		return err
	}

	return s.repo.CreateSale(sale)
}
