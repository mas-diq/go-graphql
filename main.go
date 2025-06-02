package main

import (
	"mas-diq/go-graphql/config"
	"mas-diq/go-graphql/models"
	"mas-diq/go-graphql/routes"
)

func main() {
	// Initialize database
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(&models.User{})

	// Setup routes
	r := routes.SetupRouter()

	// Start server
	r.Run(":8000")
}
