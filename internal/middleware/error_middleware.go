package middleware

import (
	web "tcg_card_battler/web-api/internal"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute the actual handler first

		// Check if there are any errors attached to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last() // Capture the most recent error
			c.JSON(c.Writer.Status(), web.APIResponse{
				Success: false,
				Message: "Operation failed",
				Error:   err.Error(),
			})
			// Abort to prevent multiple response writes
			c.Abort()
		}
	}
}
