package wrapper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data", omitempty`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

type PaginationResponse struct {
	Response
	Pagination *PaginationData `json:"pagination, omitempty"`
}

type PaginationData struct {
	CurrentPage  int64 `json:"current_page"`
	LastPage     int64 `json:"last_page"`
	PerPage      int64 `json:"per_page"`
	TotalRecords int64 `json:"total_records"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  true,
		Message: "Success",
		Code:    http.StatusOK,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Code:    code,
	})
}

func SuccessWithPagination(c *gin.Context, data interface{}, pagination *PaginationData) {
	c.JSON(http.StatusOK, PaginationResponse{
		Response: Response{
			Status:  true,
			Message: "Success",
			Code:    http.StatusOK,
			Data:    data,
		},
		Pagination: pagination,
	})
}
