package filters

import (
	"strings"
)

const (
	ASC  int = 1
	DESC     = -1
)

type Filters struct {
	Page         int64
	PageSize     int64
	Sort         string
	SortSafeList []string
}

func New(page int64, pageSize int64, sort, defaultSort string, sortSafeList []string) *Filters {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 30
	}

	if sort == "" {
		sort = defaultSort
	}

	return &Filters{Page: page, PageSize: pageSize, Sort: sort, SortSafeList: sortSafeList}
}

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ErrorResponse

func (e ValidationErrors) Error() string {
	return "invalid filter params"
}

func ValidateFilters(f *Filters) error {
	var messages ValidationErrors

	if f.Page < 0 {
		messages = append(messages, ErrorResponse{"page", "must be greater than zero"})
	}

	if f.Page >= 10_000_000 {
		messages = append(messages, ErrorResponse{"page", "must be a maximum of 10 million"})
	}

	if f.PageSize < 0 {
		messages = append(messages, ErrorResponse{"page_size", "must be greater than zero"})
	}

	if f.PageSize > 100 {
		messages = append(messages, ErrorResponse{"page_size", "must be a maximum of 100"})
	}

	if f.Sort != "" && f.SortColumn() == "" {
		messages = append(messages, ErrorResponse{"sort", "invalid sort value"})
	}

	if len(messages) != 0 {
		return messages
	}

	return nil
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	return ""
}

func (f Filters) SortDirection() int {
	if strings.HasPrefix(f.Sort, "-") {
		return DESC
	}

	return ASC
}

func (f Filters) Limit() int64 {
	return f.PageSize
}

func (f Filters) Offset() int64 {
	return (f.Page - 1) * f.PageSize
}
