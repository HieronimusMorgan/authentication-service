package response

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Response struct {
	Status  int         `json:"status"`          // HTTP status code
	Message string      `json:"message"`         // Descriptive message
	Data    interface{} `json:"data,omitempty"`  // Any additional data
	Error   interface{} `json:"error,omitempty"` // Error details (if any)
}

type ListResponse struct {
	Status    int         `json:"status"`          // HTTP status code
	Message   string      `json:"message"`         // Response message
	Timestamp string      `json:"timestamp"`       // ISO timestamp
	Data      *PagedData  `json:"data,omitempty"`  // Paginated data (optional)
	Error     interface{} `json:"error,omitempty"` // Error details (optional)
}

type PagedData struct {
	Total     int64       `json:"total"`      // Total items in the database
	PageIndex int         `json:"page_index"` // Current page index
	PageSize  int         `json:"page_size"`  // Items per page
	Items     interface{} `json:"items"`      // Actual data list
}

func SendResponse(c *gin.Context, status int, message string, data interface{}, err interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
	})
}

func SendResponseList(c *gin.Context, status int, message string, data PagedData, err interface{}) {
	c.JSON(status, ListResponse{
		Status:    status,
		Timestamp: time.Now().In(time.FixedZone("GMT+7", 7*3600)).Format(time.DateTime),
		Message:   message,
		Data:      &data,
		Error:     err,
	})
}
