package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CustomErrorHandler is a middleware that catches and formats errors
func CustomErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute the request

		// Check if there were any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			status := http.StatusNotFound

			if e, ok := err.(*gin.Error); ok {
				status = e.Meta.(int)
			}

			errorResponse := ErrorResponse{
				Status:  status,
				Code:    http.StatusText(status),
				Message: err.Error(),
			}

			c.JSON(status, errorResponse)
			c.Abort()
		}
	}
}

/*
func main() {
    r := gin.Default()

    // Add the custom error handler middleware
    r.Use(CustomErrorHandler())

    // Example route that might produce an error
    r.GET("/example", func(c *gin.Context) {
        // Simulate an error
        c.Error(gin.Error{Err: fmt.Errorf("something went wrong"), Meta: http.StatusBadRequest})
    })

    r.Run(":8080")
} */
