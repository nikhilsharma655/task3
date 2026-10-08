package services

import (
	"fmt"
	"log"
	"time"

	"task3-go/internal/database"
	"task3-go/internal/models"
)

func RunHealthChecks() {
	log.Println("Starting health checks...")

	checkDatabaseConnection()
	checkInactiveProducts()
	checkOrphanProducts()

	log.Println("Health checks completed")
}

func saveHealthCheck(checkType, status, message, details string) {
	check := models.HealthCheck{
		CheckType: checkType,
		Status:    status,
		Message:   message,
		Details:   details,
		CheckedAt: time.Now(),
	}

	if err := database.DB.Create(&check).Error; err != nil {
		log.Println("Failed to save health check:", err)
	}
}

func checkDatabaseConnection() {
	sqlDB, err := database.DB.DB()
	if err != nil {
		saveHealthCheck(
			"database_connection",
			"failed",
			"failed to access database connection",
			err.Error(),
		)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		saveHealthCheck(
			"database_connection",
			"failed",
			"database ping failed",
			err.Error(),
		)
		return
	}

	saveHealthCheck(
		"database_connection",
		"healthy",
		"database connection is healthy",
		"",
	)
}

func checkInactiveProducts() {
	var count int64

	if err := database.DB.
		Model(&models.Product{}).
		Where("status = ?", "inactive").
		Count(&count).Error; err != nil {

		saveHealthCheck(
			"inactive_products",
			"failed",
			"failed to check inactive products",
			err.Error(),
		)
		return
	}

	if count > 0 {
		message := fmt.Sprintf("%d inactive product(s) found", count)

		saveHealthCheck(
			"inactive_products",
			"warning",
			message,
			"Review inactive products",
		)
		return
	}

	saveHealthCheck(
		"inactive_products",
		"healthy",
		"no inactive products found",
		"",
	)
}

func checkOrphanProducts() {
	var products []models.Product

	if err := database.DB.Find(&products).Error; err != nil {
		saveHealthCheck(
			"orphan_products",
			"failed",
			"failed to check product records",
			err.Error(),
		)
		return
	}

	// Products are currently standalone records,
	// so an orphan check means checking for invalid records.
	var invalidCount int

	for _, product := range products {
		if product.Name == "" || product.Category == "" {
			invalidCount++
		}
	}

	if invalidCount > 0 {
		message := fmt.Sprintf("%d potentially orphan/invalid product record(s) found", invalidCount)

		saveHealthCheck(
			"orphan_products",
			"warning",
			message,
			"Product records with missing required fields were detected",
		)
		return
	}

	saveHealthCheck(
		"orphan_products",
		"healthy",
		"no orphan or invalid product records found",
		"",
	)
}
