package http

import "go-nextjs-dashboard/internal/dashboard"

type OverviewResponse struct {
	InvoiceCount  int64 `json:"invoice_count"`
	CustomerCount int64 `json:"customer_count"`
	InvoiceStatus struct {
		Paid    float64 `json:"paid"`
		Pending float64 `json:"pending"`
	} `json:"invoice_status"`
}

func newOverviewResponse(o *dashboard.Overview) OverviewResponse {
	var r OverviewResponse
	r.CustomerCount = o.CustomerCount
	r.InvoiceCount = o.InvoiceCount
	r.InvoiceStatus.Paid = o.InvoiceStatus.Paid
	r.InvoiceStatus.Pending = o.InvoiceStatus.Pending

	return r
}
