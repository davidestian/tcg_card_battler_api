package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success helper for standard 200/201 responses
func Success(c *gin.Context, message string, data interface{}) {
	if message == "" {
		message = "Success"
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
func Error(c *gin.Context, statusCode int, message string) {
	if message == "" {
		message = "Error"
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}
