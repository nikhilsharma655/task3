package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"task3-go/internal/database"
	"task3-go/internal/models"
)

type ProductRequest struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category" binding:"required"`
	Price    float64 `json:"price" binding:"gte=0"`
	Quantity int     `json:"quantity" binding:"gte=0"`
	Status   string  `json:"status" binding:"required"`
}

func CreateProduct(c *gin.Context) {

	var request ProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if request.Status != "active" && request.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status must be active or inactive",
		})
		return
	}

	product := models.Product{
		Name:     strings.TrimSpace(request.Name),
		Category: strings.TrimSpace(request.Category),
		Price:    request.Price,
		Quantity: request.Quantity,
		Status:   request.Status,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func GetProduct(c *gin.Context) {

	id := c.Param("id")

	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {

	id := c.Param("id")

	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	var request ProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Status = strings.ToLower(strings.TrimSpace(request.Status))

	if request.Status != "active" && request.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status must be active or inactive",
		})
		return
	}

	product.Name = strings.TrimSpace(request.Name)
	product.Category = strings.TrimSpace(request.Category)
	product.Price = request.Price
	product.Quantity = request.Quantity
	product.Status = request.Status

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {

	id := c.Param("id")

	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted successfully",
	})
}

func GetProducts(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	status := strings.TrimSpace(c.Query("status"))
	category := strings.TrimSpace(c.Query("category"))
	sortBy := strings.TrimSpace(c.DefaultQuery("sort", "created_at"))

	query := database.DB.Model(&models.Product{})

	// Filtering
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Count before pagination
	var total int64

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count products",
		})
		return
	}

	// Sorting
	allowedSorts := map[string]string{
		"name":       "name",
		"date":       "created_at",
		"created_at": "created_at",
		"price":      "price",
	}

	sortColumn, exists := allowedSorts[sortBy]

	if !exists {
		sortColumn = "created_at"
	}

	query = query.Order(sortColumn + " DESC")

	// Pagination
	offset := (page - 1) * limit

	var products []models.Product

	if err := query.
		Offset(offset).
		Limit(limit).
		Find(&products).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
