package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseFormInt(context *gin.Context, field string) int {
	val := context.PostForm(field)
	if val == "" {
		return 0
	}
	intVal, _ := strconv.Atoi(val)
	return intVal
}

func ParseFormUint(context *gin.Context, field string) uint {
	val := context.PostForm(field)
	if val == "" {
		return 0
	}
	uintVal, _ := strconv.ParseUint(val, 10, 32)
	return uint(uintVal)
}

func ParseFormFloat(context *gin.Context, field string) float64 {
	val := context.PostForm(field)
	if val == "" {
		return 0.0
	}
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func GetPageIndexPageSize(context *gin.Context) (int, int, error) {
	pageSize, err := strconv.Atoi(context.DefaultQuery(PageSize, "10"))
	if err != nil || pageSize <= 0 {
		return 0, 0, errors.New("page_size must be a positive integer")
	}

	pageIndex, err := strconv.Atoi(context.DefaultQuery(PageIndex, "1"))
	if err != nil || pageIndex <= 0 {
		return 0, 0, errors.New("page_index must be a positive integer")
	}
	return pageIndex, pageSize, nil
}
