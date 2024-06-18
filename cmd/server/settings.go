package settings

import (
	"fmt"
	"log"
	"os"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

func Get(key, def string) string {
	value, ok := LookupEnv(key)
	if ok {
		return value
	}
	log.Println(key, ": default value used")
	return def
}

func connectDatabase(pass, user, name string) {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv(user),
		os.Getenv(pass),
		os.Getenv(name),
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
}

func dbMigrate() {
	// Migrate the schema
	log.Println("Running migrations")
	db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{}, &models.OrderItem{})
}
