package utils

import (
	"log"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the given password matches the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseDate(date string) time.Time {
	formattedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatalf("Failed to parse date %s: %v", date, err)
	}
	return formattedDate
}

type PaginatedResult[T any] struct {
	Records      []T
	TotalRecords int64
	TotalPages   int
	CurrentPage  int
	PageSize     int
}

func Paginate[T any](db *gorm.DB, model T, page, pageSize int) (PaginatedResult[T], error) {
	var result PaginatedResult[T]
	var totalRecords int64

	// Determine the table name from the model type
	modelValue := reflect.ValueOf(model)
	tableName := db.NamingStrategy.TableName(modelValue.Type().Name())

	// Get the total count of records
	db.Table(tableName).Count(&totalRecords)

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch the paginated records
	records := make([]T, 0)
	db.Table(tableName).Limit(pageSize).Offset(offset).Find(&records)

	// Populate the result
	result.Records = records
	result.TotalRecords = totalRecords
	result.TotalPages = int((totalRecords + int64(pageSize) - 1) / int64(pageSize))
	result.CurrentPage = page
	result.PageSize = pageSize

	return result, nil
}
