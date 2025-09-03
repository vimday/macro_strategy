package main

import (
	"log"
	"macro_strategy/internal/api"
)

func main() {
	log.Println("Starting Macro Strategy Backend Server...")

	router := api.SetupRouter()

	log.Println("Server listening on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
