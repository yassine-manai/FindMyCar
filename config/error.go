package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

// Define a global variable to hold the error messages
var errorMessages map[string]map[string]ErrorMessage

func LoadErrorMessages() error {
	data, err := ioutil.ReadFile("errorCodelang.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &errorMessages)
	if err != nil {
		return err
	}
	return nil
}

func GetErrorMessage(lang string, key string) ErrorMessage {
	if messages, ok := errorMessages[lang]; ok {
		if errMessage, exists := messages[key]; exists {
			return errMessage
		}
	}
	return errorMessages["en"][key]
}
