package main

import (
	// "fmt"
	"fmt"
	"log"
	"os"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/cmd/tools"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/cmd/controllers"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/services"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"

	_ "github.com/ICOMP-UNC/newworld-agustinhernando2/docs" // Swag CLI, exec init()

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)



func ConnectDatabase(pass, user, name string) (db *gorm.DB) {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		user,
		pass,
		name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	dbMigrate(db)
	return db
}

func dbMigrate(db *gorm.DB) {
	// Migrate the schema
	log.Println("Running migrations")
	db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{}, &models.OrderItem{})
}

// @title			Order Api
// @version		1.0
// @description	This is an Order Api just for young people
// @termsOfService	http://swagger.io/terms/
func main() {
	// init envs
	tools.Init_env()
	port := tools.GetEnvValue("PORT", "3000")
	cppport := tools.GetEnvValue("CPPPORT", "8888")
	cppipv4 := tools.GetEnvValue("CPPIPV4", "192.168.100.148")

	// Configuración del repositorio de suministros
	config := &repositories.Config{
		Host: cppipv4,
		Port: cppport,
	}
	cppRepo := repositories.NewCppRepository(config)

	db_pass := tools.GetEnvValue("DB_PASSWORD", "1234")
	db_user := tools.GetEnvValue("DB_USER", "agustinhernando")
	db_name := tools.GetEnvValue("DB_NAME", "agustinhernando")

	db := ConnectDatabase(db_pass, db_user, db_name)

	// Start repositories
	userRepo := repositories.NewUserRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Start services
	authService := services.NewAuthService(userRepo, itemRepo, orderRepo)
	adminService := services.NewAdminService(userRepo, itemRepo, orderRepo, cppRepo)

	// Start controllers
	authController := controllers.NewAuthController(authService)
	adminController := controllers.NewAdminController(adminService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Configurar el motor de plantillas HTML
	engine := html.New("./templates", ".html")
	// Start Fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Routes to login and signup
	firstGroup := app.Group("/")
	firstGroup.Get("/signup", authController.GetSignupPage)
	firstGroup.Get("/login", authController.GetLoginPage)
	firstGroup.Post("/register", authController.RegisterUser)
	firstGroup.Post("/login", authController.LoginUser)

	// Authentication routes
	authGroup := app.Group("/auth", authMiddleware.UserMiddleware)
	authGroup.Get("/offers", authController.GetOffers)
	authGroup.Post("/checkout", authController.CheckoutOrders)
	authGroup.Get("/orders/:id", authController.GetOrderStatus)

	// Administration routes
	adminGroup := app.Group("/admin", authMiddleware.AdminMiddleware)
	adminGroup.Get("/dashboard", adminController.GetDashboard)
	adminGroup.Post("/updatesupplies", adminController.UpdateItemsAvailables)
	adminGroup.Patch("/orders/:id", adminController.UpdateOrderStatus)

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start server
	log.Printf("Starting server on port %s...", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
