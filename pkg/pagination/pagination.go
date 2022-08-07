package pagination

import (
	"errors"
	"math"
)

type Pagination struct {
	CurrentPage  int64 `json:"current_page,omitempty"`
	PageSize     int64 `json:"page_size,omitempty"`
	FirstPage    int64 `json:"first_page,omitempty"`
	LastPage     int64 `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

func New(totalRecords, page, pageSize int64) (*Pagination, error) {
	if totalRecords == 0 {
		return nil, errors.New("total records should not be 0")
	}

	return &Pagination{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int64(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}, nil
}
