package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Status: "success",
		Data:   data,
	})
}

func ErrorResponse(c *gin.Context, error string, errorCode int) {
	c.JSON(errorCode, APIResponse{
		Status: "error",
		Error:  error,
	})
}
