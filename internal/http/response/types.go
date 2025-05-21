package response

import (
	"go-nextjs-dashboard/internal/listing"
)

type Response[T any] struct {
	Data T `json:"data"`
}

func New[T any](data T) Response[T] {
	return Response[T]{
		Data: data,
	}
}

type PaginatedResponse[T any] struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
	Data        []T   `json:"data"`
}

// PaginateList converts a listing.Result[S] to PaginatedResponse[R]
func PaginateList[S, R any](res listing.Result[S], mapper func(S) R) *PaginatedResponse[R] {
	data := make([]R, len(res.Items))
	for i, v := range res.Items {
		data[i] = mapper(v)
	}

	return &PaginatedResponse[R]{
		Total:       res.Total,
		CurrentPage: res.Page,
		PerPage:     res.PerPage,
		HasNext:     res.HasNext,
		HasPrev:     res.HasPrev,
		Data:        data,
	}
}

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
