package service

import (
	"database/sql"
	"go-nextjs-dashboard/config"
	"go-nextjs-dashboard/model"
	"go-nextjs-dashboard/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FilteredCustomerDTO struct {
	ID            uuid.UUID
	Name          string
	Email         string
	ImageURL      string
	TotalInvoices uint32
	TotalPending  float64
	TotalPaid     float64
}

type CustomerService struct {
	DB *gorm.DB
}

func NewCustomerService() *CustomerService {
	return &CustomerService{DB: config.DB}
}

func (s *CustomerService) GetCustomers() ([]model.Customer, error) {
	var customers []model.Customer

	err := s.DB.Select("id, name").Order("name ASC").Find(&customers).Error

	return customers, err
}

func (s *CustomerService) GetPaginatedCustomers(page int, size int) (*utils.PaginatedResult[model.Customer], error) {
	paginate, err := utils.Paginate(s.DB, model.Customer{}, 1, 5)

	if err != nil {
		return nil, err
	}

	return &paginate, nil
}

func (s *CustomerService) GetFilteredCustomers(search string) ([]FilteredCustomerDTO, error) {
	// type filteredCustomer struct {
	// 	ID            uuid.UUID `json:"id"`
	// 	Name          string    `json:"name"`
	// 	Email         string    `json:"email"`
	// 	ImageURL      string    `json:"image_url"`
	// 	TotalInvoices uint32    `json:"total_invoices"`
	// 	TotalPending  float64   `json:"total_pending"`
	// 	TotalPaid     float64   `json:"total_paid"`
	// }

	var result []FilteredCustomerDTO

	err := s.DB.Model(&model.Customer{}).
		Select(`
			customers.id,
			customers.name,
			customers.email,
			customers.image_url,
			COUNT(invoices.id) AS total_invoices,
			SUM(CASE WHEN invoices.status = 'pending' THEN invoices.amount ELSE 0 END) AS total_pending,
			SUM(CASE WHEN invoices.status = 'paid' THEN invoices.amount ELSE 0 END) AS total_paid
		`).
		Joins("LEFT JOIN invoices ON customers.id = invoices.customer_id").
		Where(`
			customers.name LIKE @search OR
        	customers.email LIKE @search
		`, sql.Named("search", "%"+search+"%")).
		Group("customers.id, customers.name, customers.email, customers.image_url").
		Order("customers.name").
		Scan(&result).
		Error

	return result, err
}
