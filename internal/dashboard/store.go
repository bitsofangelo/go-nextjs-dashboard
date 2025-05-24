package dashboard

import "context"

type Store interface {
	GetOverview(ctx context.Context) (*Overview, error)
	ListMonthlyRevenues(ctx context.Context) ([]MonthlyRevenue, error)
}

type Overview struct {
	InvoiceCount  int64
	CustomerCount int64
	InvoiceStatus InvoiceStatus
}

type InvoiceStatus struct {
	Paid, Pending float64
}

type MonthlyRevenue struct {
	Month  string
	Amount float64
}

var monthOrder = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}
