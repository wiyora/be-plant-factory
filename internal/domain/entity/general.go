package entity

type PaginationResult struct {
	Page       uint64
	Limit      uint64
	Total      uint64
	TotalPages uint64
}
