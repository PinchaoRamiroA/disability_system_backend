package main

import (
	"log"
	"net/http"

	"disability_system_backend/internal/shared/database"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.NewConnection()

	if err != nil {
		log.Fatal(err)
	}

	err = database.HealthCheck(db)

	if err != nil {
		log.Fatal("database health check failed")
	}

	log.Println("database connected successfully")
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}
