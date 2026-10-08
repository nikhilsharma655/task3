package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task3-go/internal/database"
	"task3-go/internal/models"
	"task3-go/internal/services"
)

func UploadFile(c *gin.Context) {

	// Get uploaded file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}

	// Check file extension
	extension := filepath.Ext(file.Filename)

	if extension != ".csv" && extension != ".xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only CSV and Excel (.xlsx) files are supported",
		})
		return
	}

	// Create unique job ID
	jobID := uuid.New().String()

	// Create uploads directory if it doesn't exist
	err = os.MkdirAll("uploads", 0755)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create upload directory",
		})
		return
	}

	// Create unique filename
	fileName := fmt.Sprintf(
		"%s_%d%s",
		jobID,
		time.Now().Unix(),
		extension,
	)

	filePath := filepath.Join("uploads", fileName)

	// Save uploaded file
	err = c.SaveUploadedFile(file, filePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save uploaded file",
		})
		return
	}

	// Create import job
	job := models.ImportJob{
		ID:       jobID,
		FileName: file.Filename,
		Status:   "processing",
	}

	err = database.DB.Create(&job).Error

	if err != nil {
		os.Remove(filePath)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create import job",
		})
		return
	}

	// Process file in background
	go func() {

		go func() {
			switch extension {
			case ".csv":
				services.ProcessCSV(filePath, jobID)

			case ".xlsx":
				services.ProcessExcel(filePath, jobID)
			}
		}()

	}()

	// Return immediately
	c.JSON(http.StatusAccepted, gin.H{
		"message":   "file uploaded successfully",
		"job_id":    jobID,
		"status":    "processing",
		"file_name": file.Filename,
	})
}
func GetImportStatus(c *gin.Context) {

	jobID := c.Param("job_id")

	var job models.ImportJob

	err := database.DB.
		Where("id = ?", jobID).
		First(&job).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "import job not found",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}
