package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"task3-go/internal/database"
	"task3-go/internal/handlers"
	"task3-go/internal/middleware"
	"task3-go/internal/models"
	"task3-go/internal/services"
)

func main() {

	database.ConnectDatabase()

	
	if err := database.DB.AutoMigrate(
		&models.Product{},
		&models.ImportJob{},
		&models.User{},
		&models.HealthCheck{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	
	services.StartHealthCheckScheduler()

	
	router := gin.Default()

	
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task 3 Go API is running",
		})
	})

	

	router.POST("/api/auth/register", handlers.Register)
	router.POST("/api/auth/login", handlers.Login)

	

	router.POST("/api/import", handlers.UploadFile)
	router.GET("/api/import/:job_id", handlers.GetImportStatus)


	router.GET(
		"/api/health/check",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager"),
		handlers.RunHealthCheck,
	)

	

	router.GET(
		"/api/protected",
		middleware.AuthMiddleware(),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "You accessed a protected route",
				"user_id": c.MustGet("user_id"),
				"email":   c.MustGet("email"),
				"role":    c.MustGet("role"),
			})
		},
	)

	router.GET(
		"/api/admin",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome Admin",
			})
		},
	)

	
	router.GET(
		"/api/manager",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager"),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome Admin/Manager",
			})
		},
	)

	
	router.GET(
		"/api/viewer",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager", "viewer"),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome authenticated user",
			})
		},
	)

	router.POST(
		"/api/products",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager"),
		handlers.CreateProduct,
	)

	
	router.GET(
		"/api/products",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager", "viewer"),
		handlers.GetProducts,
	)

	router.GET(
		"/api/products/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager", "viewer"),
		handlers.GetProduct,
	)

	
	router.PUT(
		"/api/products/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin", "manager"),
		handlers.UpdateProduct,
	)

	
	router.DELETE(
		"/api/products/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRole("admin"),
		handlers.DeleteProduct,
	)


	log.Println("Server running on http://localhost:8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
