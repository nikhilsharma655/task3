package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"task3-go/internal/database"
	"task3-go/internal/models"
)

type ImportResult struct {
	TotalRows   int
	SuccessRows int
	FailedRows  int
	Errors      []string
}

func ProcessCSV(filePath string, jobID string) {

	result := ImportResult{}

	file, err := os.Open(filePath)
	if err != nil {
		updateJobFailed(jobID, err.Error())
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		updateJobFailed(jobID, "Unable to read CSV header")
		return
	}

	fmt.Println("CSV Header:", header)

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			result.FailedRows++
			result.Errors = append(
				result.Errors,
				fmt.Sprintf("Failed to read row: %v", err),
			)
			continue
		}

		result.TotalRows++

		product, validationError := validateProductRow(row)

		if validationError != nil {
			result.FailedRows++
			result.Errors = append(
				result.Errors,
				fmt.Sprintf("Row %d: %s", result.TotalRows, validationError.Error()),
			)
			continue
		}

		err = database.DB.Create(&product).Error

		if err != nil {
			result.FailedRows++
			result.Errors = append(
				result.Errors,
				fmt.Sprintf("Row %d: database error: %v", result.TotalRows, err),
			)
			continue
		}

		result.SuccessRows++
	}

	now := time.Now()

	status := "completed"

	if result.FailedRows > 0 && result.SuccessRows == 0 {
		status = "failed"
	}

	database.DB.Model(&models.ImportJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":        status,
			"total_rows":    result.TotalRows,
			"success_rows":  result.SuccessRows,
			"failed_rows":   result.FailedRows,
			"error_details": strings.Join(result.Errors, "\n"),
			"completed_at":  now,
		})

	fmt.Printf(
		"Import completed: Total=%d Success=%d Failed=%d\n",
		result.TotalRows,
		result.SuccessRows,
		result.FailedRows,
	)
}
func ProcessExcel(filePath string, jobID string) {

	result := ImportResult{}

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		updateJobFailed(jobID, err.Error())
		return
	}
	defer file.Close()

	sheetName := file.GetSheetName(0)

	rows, err := file.GetRows(sheetName)
	if err != nil {
		updateJobFailed(jobID, "Unable to read Excel rows")
		return
	}

	if len(rows) == 0 {
		updateJobFailed(jobID, "Excel file is empty")
		return
	}

	fmt.Println("Excel Header:", rows[0])

	for i, row := range rows[1:] {

		result.TotalRows++

		product, validationError := validateProductRow(row)

		if validationError != nil {
			result.FailedRows++

			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"Row %d: %s",
					i+2,
					validationError.Error(),
				),
			)

			continue
		}

		err = database.DB.Create(&product).Error

		if err != nil {
			result.FailedRows++

			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"Row %d: database error: %v",
					i+2,
					err,
				),
			)

			continue
		}

		result.SuccessRows++
	}

	now := time.Now()

	status := "completed"

	if result.FailedRows > 0 && result.SuccessRows == 0 {
		status = "failed"
	}

	database.DB.Model(&models.ImportJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":        status,
			"total_rows":    result.TotalRows,
			"success_rows":  result.SuccessRows,
			"failed_rows":   result.FailedRows,
			"error_details": strings.Join(result.Errors, "\n"),
			"completed_at":  now,
		})

	fmt.Printf(
		"Excel import completed: Total=%d Success=%d Failed=%d\n",
		result.TotalRows,
		result.SuccessRows,
		result.FailedRows,
	)
}

func validateProductRow(row []string) (models.Product, error) {

	if len(row) < 5 {
		return models.Product{}, fmt.Errorf(
			"expected 5 columns, got %d",
			len(row),
		)
	}

	name := strings.TrimSpace(row[0])
	category := strings.TrimSpace(row[1])
	priceString := strings.TrimSpace(row[2])
	quantityString := strings.TrimSpace(row[3])
	status := strings.TrimSpace(row[4])

	if name == "" {
		return models.Product{}, fmt.Errorf("product name is required")
	}

	if category == "" {
		return models.Product{}, fmt.Errorf("category is required")
	}

	price, err := strconv.ParseFloat(priceString, 64)

	if err != nil || price < 0 {
		return models.Product{}, fmt.Errorf("invalid price")
	}

	quantity, err := strconv.Atoi(quantityString)

	if err != nil || quantity < 0 {
		return models.Product{}, fmt.Errorf("invalid quantity")
	}

	if status != "active" && status != "inactive" {
		return models.Product{}, fmt.Errorf(
			"status must be active or inactive",
		)
	}

	product := models.Product{
		Name:     name,
		Category: category,
		Price:    price,
		Quantity: quantity,
		Status:   status,
	}

	return product, nil
}

func updateJobFailed(jobID string, errorMessage string) {

	now := time.Now()

	database.DB.Model(&models.ImportJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":        "failed",
			"error_details": errorMessage,
			"completed_at":  now,
		})
}
