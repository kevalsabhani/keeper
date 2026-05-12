package models

// Pagination holds domain-level pagination metadata returned by service list methods.
// It is framework-agnostic and contains no HTTP or transport concerns.
type Pagination struct {
	CurrentPage int
	TotalPages  int
	TotalCount  int
}

// NewPagination calculates and returns a Pagination value for a given page, limit, and total record count.
func NewPagination(page, limit, total int) *Pagination {
	return &Pagination{
		CurrentPage: page,
		TotalPages:  (total + limit - 1) / limit,
		TotalCount:  total,
	}
}
