package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mgmacri/pool-maintenance-app/internal/delivery"
)

func main() {
	r := gin.Default()
	healthHandler := delivery.NewHealthHandler()
	r.GET("/health", healthHandler.Check)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
