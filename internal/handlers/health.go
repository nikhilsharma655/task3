package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task3-go/internal/database"
	"task3-go/internal/models"
	"task3-go/internal/services"
)

func RunHealthCheck(c *gin.Context) {
	services.RunHealthChecks()

	var checks []models.HealthCheck

	if err := database.DB.
		Order("checked_at DESC").
		Limit(20).
		Find(&checks).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch health check results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "health checks completed",
		"checks":  checks,
	})
}
