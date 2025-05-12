package dashboard

type Overview struct {
	InvoiceCount  int64
	CustomerCount int64
	InvoiceStatus InvoiceStatus
}

type InvoiceStatus struct {
	Paid, Pending float64
}
