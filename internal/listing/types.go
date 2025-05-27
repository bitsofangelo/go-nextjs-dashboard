package listing

import "gorm.io/gorm"

type SortOrder int

const (
	SortOldest SortOrder = iota
	SortLatest
)

// Page carries the page request details.
type Page struct {
	Page int
	Size int
}

func NewPage(page int, size int) Page {
	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 10
	} else if size > 100 {
		size = 100
	}

	return Page{
		Page: page,
		Size: size,
	}
}

func (p Page) Scope() func(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Size

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(p.Size)
	}
}

// Result is a generic container returned by store methods.
type Result[T any] struct {
	Items   []T
	Total   int64
	Page    int
	PerPage int
	HasNext bool
	HasPrev bool
	// NextCursor string
}

func NewResult[T any](items []T, p Page, total int64) Result[T] {
	return Result[T]{
		Items:   items,
		Total:   total,
		Page:    p.Page,
		PerPage: p.Size,
		HasNext: int64(p.Page*p.Size) < total,
		HasPrev: p.Page > 1,
	}
}
