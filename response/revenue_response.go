package response

import "go-nextjs-dashboard/model"

type RevenueResponse struct {
	Month   string  `json:"month"`
	Revenue float32 `json:"revenue"`
}

func NewRevenueResponse(revenues []model.Revenue) map[string]any {
	revenuResponse := make([]RevenueResponse, len(revenues))

	for i, revenue := range revenues {
		revenuResponse[i] = RevenueResponse{
			Month:   string(revenue.Month),
			Revenue: revenue.Revenue,
		}
	}

	response := map[string]any{"data": revenuResponse}

	return response
}
