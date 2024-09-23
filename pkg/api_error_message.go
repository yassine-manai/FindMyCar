package pkg

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetAllErrorCode godoc
//
//	@Summary		Get all error messages or specific one
//	@Description	Get a list of all error messages or a specific one by code
//	@Tags			Errors
//	@Produce		json
//	@Param			code	query		string	false	"Error code to fetch specific error message"
//	@Success		200	{object}	[]ErrorMessage
//	@Router			/fyc/errors [get]
func GetAllErrorCode(c *gin.Context) {
	ctx := context.Background()
	codeReq := c.Query("code")

	if codeReq != "" {
		errorMessage, err := GetErrorMessageByCode(ctx, codeReq)
		if err != nil {
			log.Err(err).Msg("Error getting error message by code")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting error message",
				"code":    10,
			})
			return
		}

		if errorMessage == nil {
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

	errors, err := GetAllErrors(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all error messages")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all error messages",
			"code":    10,
		})
		return
	}

	if len(errors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No error messages found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, errors)
}

// CreateErrorMessageAPI godoc
//
//	@Summary		Create a new error message
//	@Description	Create a new error message
//	@Tags			Errors
//	@Accept		json
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

	c.JSON(http.StatusCreated, errMsg)
}

// UpdateErrorMessageAPI godoc
//
//	@Summary		Update an existing error message
//	@Description	Update an existing error message by code
//	@Tags			Errors
//	@Accept		json
//	@Produce		json
//	@Param			code	query		string	true	"Error message code"
//	@Param			errMsg	body		ErrorMessage	true	"Updated error message object"
//	@Success		200	{object}	ErrorMessage
//	@Router			/fyc/errors/{code} [put]
func UpdateErrorMessageAPI(c *gin.Context) {
	code := c.Query("code")
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
//	@Summary		Delete an error message
//	@Description	Delete an error message by code
//	@Tags			Errors
//	@Param			code	path		string	true	"Error message code"
//	@Success		204	{object}	nil
//	@Router			/fyc/errors/{code} [delete]
func DeleteErrorMessageAPI(c *gin.Context) {
	code := c.Param("code")

	ctx := context.Background()
	if err := DeleteErrorMessage(ctx, code); err != nil {
		log.Error().Err(err).Msg("Failed to delete error message")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Failed to delete error message",
			"code":    500,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
