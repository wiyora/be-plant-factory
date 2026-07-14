package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Pagination
type Pagination struct {
	Page     uint64
	PageSize uint64
}

func (p Pagination) Offset() uint64 {
	if p.Page == 0 {
		return 0
	}

	return (p.Page - 1) * p.PageSize
}

func (p Pagination) ToResult(total uint64) PaginationResult {
	return PaginationResult{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: (total + p.PageSize - 1) / p.PageSize,
	}
}

// Search
type Search string

func (s Search) HasSearch() bool {
	return len([]rune(s)) >= 3
}

func (s Search) IsEmpty() bool {
	return s == ""
}

func (s Search) String() string {
	return string(s)
}

// Order
type Order struct {
	OrderBy string
	SortBy  SortDirection
}

func (o Order) IsEmpty() bool {
	return o.OrderBy == "" || o.SortBy == ""
}

func (o Order) String() string {
	return fmt.Sprintf("%s %s", o.OrderBy, o.SortBy)
}

// Date Range
type DateRange struct {
	Start time.Time
	End   time.Time
}

func (d DateRange) IsEmpty() bool {
	return d.Start.IsZero() && d.End.IsZero()
}

// DropdownFilter
type DropdownFilter struct {
	Search     Search
	Pagination Pagination
	ActiveIDs  []uuid.UUID
}

// DropdownItem
type DropdownItem struct {
	ID   uuid.UUID
	Name string
}
