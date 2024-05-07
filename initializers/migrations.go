package initializers

import (
	"fmt"

	"github.com/woaitsAryan/stuneckt-task/models"
)

func AutoMigrate() {
	fmt.Println("\nStarting Migrations...")
	DB.AutoMigrate(
		&models.User{},

	)
	fmt.Println("Migrations Finished!")
}
