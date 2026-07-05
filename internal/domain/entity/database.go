package entity

// Pagination Sort Direction
type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

func (d SortDirection) Valid() bool {
	switch d {
	case Asc, Desc:
		return true
	default:
		return false
	}
}

func (d SortDirection) String() string {
	return string(d)
}

// Soft Delete Filter
type SoftDeleteFilter string

const (
	WithoutDeleted SoftDeleteFilter = "without"
	OnlyDeleted    SoftDeleteFilter = "only"
	WithDeleted    SoftDeleteFilter = "with"
)

func (f SoftDeleteFilter) Valid() bool {
	switch f {
	case WithoutDeleted, OnlyDeleted, WithDeleted:
		return true
	default:
		return false
	}
}

func (f SoftDeleteFilter) String() string {
	return string(f)
}
