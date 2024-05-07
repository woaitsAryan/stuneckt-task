package initializers

import (
	"fmt"

	"github.com/woaitsAryan/stuneckt-task/models"
)

func AutoMigrate() {
	fmt.Println("\nStarting Migrations...")
	DB.AutoMigrate(
		&models.User{},
		&models.Post{},

	)
	fmt.Println("Migrations Finished!")
}
