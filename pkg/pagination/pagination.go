package pagination

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Kachyr/findyourpet/findyourpet-backend/pkg/constants"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, pageSize, err := GetPageQueryParams(c)
		if err != nil {
			c.Status(http.StatusBadRequest)
			log.Err(err).Msg("Error")
		}

		if page < 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		dbClone := db.Session(&gorm.Session{})
		var total int64
		dbClone.Count(&total)

		totalPages := CalculateTotalPages(int(total), pageSize)

		c.Set("page", page)
		c.Set("pageSize", pageSize)
		c.Set("totalPages", totalPages)

		return db.Offset(offset).Limit(pageSize)
	}
}

func GetPageQueryParams(c *gin.Context) (int, int, error) {
	pageStr := c.DefaultQuery(constants.Page, "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return 0, 0, err
	}

	pageSizeStr := c.DefaultQuery(constants.PageSize, "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return 0, 0, err
	}

	return page, pageSize, nil
}

func CalculateTotalPages(totalElements, size int) int {
	if totalElements < 0 {
		// Handle negative total elements
		return 0
	}

	if size <= 0 {
		// Handle invalid page size
		return 0
	}

	return int(math.Ceil(float64(totalElements) / float64(size)))
}
