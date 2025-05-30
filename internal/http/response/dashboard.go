package response

import "github.com/gelozr/go-dash/internal/dashboard"

type Overview struct {
	InvoiceCount  int64 `json:"invoice_count"`
	CustomerCount int64 `json:"customer_count"`
	InvoiceStatus struct {
		Paid    float64 `json:"paid"`
		Pending float64 `json:"pending"`
	} `json:"invoice_status"`
}

func ToOverview(o *dashboard.Overview) Overview {
	var r Overview
	r.CustomerCount = o.CustomerCount
	r.InvoiceCount = o.InvoiceCount
	r.InvoiceStatus.Paid = o.InvoiceStatus.Paid
	r.InvoiceStatus.Pending = o.InvoiceStatus.Pending

	return r
}

type MonthlyRevenue struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

func ToMonthlyRevenue(r dashboard.MonthlyRevenue) MonthlyRevenue {
	return MonthlyRevenue{
		Month:  r.Month,
		Amount: r.Amount,
	}
}

func ToMonthlyRevenueList(data []dashboard.MonthlyRevenue) []MonthlyRevenue {
	return ToList(data, ToMonthlyRevenue)
}
