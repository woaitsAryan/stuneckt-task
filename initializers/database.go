package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		CONFIG.DB_HOST, CONFIG.DB_USER, CONFIG.DB_PASSWORD, CONFIG.DB_NAME, CONFIG.DB_PORT)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to Connect to the database")
	} else {
		log.Println("Connected to database!")
	}
}
