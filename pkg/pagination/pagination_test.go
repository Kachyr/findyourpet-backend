package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPageQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid query parameters", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/?page=2&page_size=20", nil)

		page, pageSize, err := pagination.GetPageQueryParams(c)

		assert.NoError(t, err)
		assert.Equal(t, 2, page)
		assert.Equal(t, 20, pageSize)
	})

	t.Run("invalid page parameter", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/?page=invalid&page_size=20", nil)

		page, pageSize, err := pagination.GetPageQueryParams(c)

		assert.Error(t, err)
		assert.Equal(t, 0, page)
		assert.Equal(t, 0, pageSize)
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("invalid pageSize parameter", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/?page=2&page_size=invalid", nil)

		page, pageSize, err := pagination.GetPageQueryParams(c)

		assert.Error(t, err)
		assert.Equal(t, 0, page)
		assert.Equal(t, 0, pageSize)
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	})

	t.Run("default query parameters", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest("GET", "/", nil)

		page, pageSize, err := pagination.GetPageQueryParams(c)

		assert.NoError(t, err)
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, pageSize)
	})
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name          string
		totalElements int
		size          int
		expectedPages int
	}{
		{
			name:          "Test with totalElements divisible by size",
			totalElements: 100,
			size:          10,
			expectedPages: 10,
		},
		{
			name:          "Test with totalElements not divisible by size",
			totalElements: 101,
			size:          10,
			expectedPages: 11,
		},
		{
			name:          "Test with zero totalElements",
			totalElements: 0,
			size:          10,
			expectedPages: 0,
		},
		{
			name:          "Test with zero size",
			totalElements: 100,
			size:          0,
			expectedPages: 0,
		},
		{
			name:          "Test with negative totalElements",
			totalElements: -100,
			size:          10,
			expectedPages: 0,
		},
		{
			name:          "Test with negative size",
			totalElements: 100,
			size:          -10,
			expectedPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := pagination.CalculateTotalPages(tt.totalElements, tt.size)
			if pages != tt.expectedPages {
				t.Errorf("CalculateTotalPages(%d, %d) = %d; want %d", tt.totalElements, tt.size, pages, tt.expectedPages)
			}
		})
	}
}
