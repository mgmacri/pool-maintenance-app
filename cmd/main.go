package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mgmacri/pool-maintenance-app/internal/delivery"
)

func main() {
	r := gin.Default()
	healthHandler := delivery.NewHealthHandler()
	r.GET("/health", healthHandler.Check)
	r.Run(":8080")
}
