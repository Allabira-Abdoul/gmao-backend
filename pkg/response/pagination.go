package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Pagination parameters.
type Pagination struct {
	Page    int
	PerPage int
}

// GetPagination parses pagination parameters from query string,
// falling back to default values if not provided.
func GetPagination(c *gin.Context, defaultPage, defaultPerPage int) Pagination {
	p := Pagination{
		Page:    defaultPage,
		PerPage: defaultPerPage,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &p.Page)
	}
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		fmt.Sscanf(perPageStr, "%d", &p.PerPage)
	}

	return p
}

// NewMeta calculates pagination metadata.
func NewMeta(page, perPage int, total int64) *Meta {
	totalPages := total / int64(perPage)
	if total%int64(perPage) != 0 {
		totalPages++
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
