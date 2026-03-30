package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logic before the request is handled
		fmt.Println("Request received at:", time.Now())

		// Pass control to the next middleware/handler
		c.Next()

		// Logic after the request is handled (response time, etc.)
		fmt.Println("Request finished at:", time.Now())
	}
}
