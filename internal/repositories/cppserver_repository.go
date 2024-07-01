package repositories

import (
    "encoding/json"
    "fmt"
	"io"
    "net/http"
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
)

// Config cppserver configuration
type Config struct {
    Host string
    Port string
}

// cppRepository maneja la configuración y las funciones del servidor
type cppRepository struct {
    Config *Config
}

// CppRepository define la interfaz para obtener suministros
type CppRepository interface {
    GetSupplies() ([]models.Item, error)
}

// NewCppNewCppRepository crea una nueva instancia de cppRepository
func NewCppRepository(config *Config) CppRepository {
    return &cppRepository{
        Config: config,
    }
}

// GetSupplies realiza una solicitud GET para obtener los suministros
func (handler *cppRepository) GetSupplies() ([]models.Item, error) {
    url := fmt.Sprintf("http://%s:%s/supplies", handler.Config.Host, handler.Config.Port)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error making GET request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    var data map[string]map[string]int
    // Deserializar la respuesta JSON
    if err := json.Unmarshal(body, &data); err != nil {
        return nil, fmt.Errorf("error deserializing the JSON: %v", err)
    }

    items := normalizeData(data)
    return items, nil
}

// normalizeData convierte los datos anidados en una lista de modelos Item
func normalizeData(data map[string]map[string]int) []models.Item {
    var items []models.Item
    for category, itemsMap := range data {
        for name, quantity := range itemsMap {
            item := models.Item{
                Name:     name,
                Category: category,
                Price:    float64(1.0),
                Quantity: int(quantity),
                Status:   "Available",
            }
            items = append(items, item)
        }
    }
    return items
}