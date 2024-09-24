package pkg

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetAllErrorCode godoc
//
//	@Summary		Get all error messages or a specific one
//	@Description	Get a list of all error messages or a specific one by code and language
//	@Tags			Errors
//	@Produce		json
//	@Param			code	query		string	false	"Error code to fetch specific error message"
//	@Param			lang	query		string	false	"Language of the error message"
//	@Success		200	{object}	[]ErrorMessage
//	@Router			/fyc/errors [get]
func GetAllErrorCode(c *gin.Context) {
	ctx := context.Background()
	code := c.Query("code")
	langReq := c.Query("lang")
	codeReq, _ := strconv.Atoi(code)

	if code != "" {
		errorMessage, err := GetErrorMessageByFilter(ctx, codeReq, langReq)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving error message by code")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error retrieving error message",
				"code":    10,
			})
			return
		}

		if errorMessage.Code == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No error message found for the provided code",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, errorMessage)
		return
	} else {
		errorMessage, err := GetErrorMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving error message by code")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error retrieving error message",
				"code":    10,
			})
			return
		}

		if len(errorMessage) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No error message found for the provided code",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, errorMessage)
		return
	}
}

// CreateErrorMessageAPI godoc
//
//	@Summary		Create a new error message
//	@Description	Create a new error message
//	@Tags			Errors
//	@Accept			json
//	@Produce		json
//	@Param			errMsg	body		ErrorMessage	true	"Error message object"
//	@Success		201	{object}	ErrorMessage
//	@Router			/fyc/errors [post]
func CreateErrorMessageAPI(c *gin.Context) {
	var errMsg ErrorMessage
	if err := c.ShouldBindJSON(&errMsg); err != nil {
		log.Error().Err(err).Msg("Invalid input for error message")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	ctx := context.Background()
	if err := CreateErrorMessage(ctx, &errMsg); err != nil {
		log.Error().Err(err).Msg("Failed to create error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Failed to create error message",
			"code":    500,
		})
		return
	}

	log.Info().Int("code", errMsg.Code).Msg("Error Message created successfully")
	c.JSON(http.StatusCreated, errMsg)
}

// UpdateErrorMessageAPI godoc
//
//	@Summary		Update an existing error message
//	@Description	Update an existing error message by code
//	@Tags			Errors
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"Error message code"
//	@Param			errMsg	body		ErrorMessage	true	"Updated error message object"
//	@Success		200	{object}	ErrorMessage
//	@Router			/fyc/errors/{code} [put]
func UpdateErrorMessageAPI(c *gin.Context) {
	codeStr := c.Query("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid error code format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid code format",
			"message": "Code must be a valid integer",
			"code":    400,
		})
		return
	}

	var errMsg ErrorMessage
	if err := c.ShouldBindJSON(&errMsg); err != nil {
		log.Error().Err(err).Msg("Invalid input for error message")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	errMsg.Code = code

	ctx := context.Background()
	if err := UpdateErrorMessage(ctx, &errMsg); err != nil {
		log.Error().Err(err).Msg("Failed to update error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Failed to update error message",
			"code":    500,
		})
		return
	}

	c.JSON(http.StatusOK, errMsg)
}

// DeleteErrorMessageAPI godoc
//
//	@Summary		Delete a specific language from an error message
//	@Description	Delete a specific language entry from the messages field of an error message by code
//	@Tags			Errors
//	@Param			code	query		string	true	"Error message code"
//	@Param			lang	query		string	true	"Language of the error message"
//	@Success		204	{object}	nil
//	@Router			/fyc/errors [delete]
func DeleteErrorMessageAPI(c *gin.Context) {
	codeStr := c.Query("code")
	langQuery := c.Query("lang")

	// Validate the code parameter
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid error code format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid code format",
			"message": "Code must be a valid integer",
			"code":    400,
		})
		return
	}

	// Ensure language parameter is provided
	if langQuery == "" {
		log.Error().Msg("Missing language parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing language parameter",
			"message": "A valid language parameter must be provided",
			"code":    400,
		})
		return
	}

	// Call the DeleteErrorMessage function to remove the specific language entry
	ctx := context.Background()
	rowsAffected, err := DeleteErrorMessage(ctx, code, langQuery)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete error message",
			"message": err.Error(),
			"code":    500,
		})
		return
	}

	// Handle case where no rows were affected (i.e., no such code or language found)
	if rowsAffected == 0 {
		log.Info().Int("code", code).Str("lang", langQuery).Msg("No error message found with the specified code and language")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No error message found for the provided code and language",
			"code":    404,
		})
		return
	}

	// Success response
	log.Info().Int("code", code).Str("lang", langQuery).Msg("Error message language deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "Deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
