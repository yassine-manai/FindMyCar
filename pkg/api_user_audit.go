package pkg

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

)

// GetUserAuditAPI godoc
//
//	@Summary		Get UserAudit or specific UserAudit by ID
//	@Description	Get a list of UserAudit or a specific UserAudit by ID with optional extra data
//	@Tags			User Audit
//	@Produce		json
//	@Param			id		query		int	false	"UserAudit ID"
//	@Success		200		{object}	UserAudit		"List of UserAudit or a single UserAudit"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No UserAudit found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid UserAudit ID"
//	@Router			/fyc/UserAudit [get]
func GetUserAuditAPI(c *gin.Context) {
	log.Debug().Msg("GetUserAuditAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		log.Info().Str("UserAudit_id", idStr).Msg("Fetching UserAudit by ID")

		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid UserAudit ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		UserAudit, err := GetUserAuditById(ctx, id)
		if err != nil {
			log.Err(err).Str("UserAudit_id", idStr).Msg("Error retrieving UserAudit by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "UserAudit not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("UserAudit_id", idStr).Msg("UserAudit fetched successfully")
		c.JSON(http.StatusOK, UserAudit)
		return
	}

	UserAudit, err := GetAllUserAudits(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all UserAudit")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all UserAudit",
			"code":    10,
		})
		return
	}

	if len(UserAudit) == 0 {
		log.Info().Msg("No UserAudit found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No UserAudit found",
			"code":    9,
		})
		return
	}

	log.Info().Int("UserAudit_count", len(UserAudit)).Msg("UserAudit fetched successfully")
	c.JSON(http.StatusOK, UserAudit)
}

// CreateUserAudit godoc
//
//	@Summary		Add a new UserAudit
//	@Description	Add a new UserAudit to the database
//	@Tags			User Audit
//	@Accept			json
//	@Produce		json
//	@Param			UserAudit	body		UserAudit	true	"UserAudit data"
//	@Success		201		{object}	UserAudit	"UserAudit created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create a new UserAudit"
//	@Router			/fyc/UserAudit [post]
func CreateUserAuditAPI(c *gin.Context) {
	log.Debug().Msg("CreateUserAuditAPI request")
	ctx := context.Background()
	var newUserAudit UserAudit

	if err := c.ShouldBindJSON(&newUserAudit); err != nil {
		log.Err(err).Msg("Invalid input for new UserAudit")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if err := CreateUserAudit(ctx, &newUserAudit); err != nil {
		log.Err(err).Msg("Error creating new UserAudit")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new UserAudit",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("UserAudit_id", newUserAudit.ID).Msg("UserAudit created successfully")
	c.JSON(http.StatusCreated, newUserAudit)
}

// UpdateUserAudit godoc
//
//	@Summary		Update a UserAudit by ID
//	@Description	Update an existing UserAudit by ID
//	@Tags			User Audit
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int			true	"UserAudit ID"
//	@Param			UserAudit	body		UserAudit		true	"Updated UserAudit data"
//	@Success		200		{object}	map[string]interface{}	"UserAudit updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload or ID mismatch"
//	@Failure		404		{object}	map[string]interface{}	"UserAudit not found"
//	@Failure		500		{object}	map[string]interface{}	"Failed to update UserAudit"
//	@Router			/fyc/UserAudit [put]
func UpdateUserAuditAPI(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	var updates UserAudit
	ctx := context.Background()
	log.Info().Str("UserAudit_id", idStr).Msg("Updating UserAudit")

	if err != nil {
		log.Error().Str("UserAudit_id", idStr).Msg("Invalid ID format for UserAudit update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for UserAudit update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.ID != id {
		log.Info().Msg("The ID in the request body does not match the query ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the query ID",
			"code":    13,
		})
		return
	}

	rowsAffected, err := UpdateUserAudit(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating UserAudit")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update UserAudit",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("UserAudit_id", idStr).Int64("Rows Affected", rowsAffected).Msg("No UserAudit found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No UserAudit found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("UserAudit_id", idStr).Msg("UserAudit updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":       "UserAudit updated successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteUserAuditAPI godoc
//
//	@Summary		Soft delete a UserAudit
//	@Description	Soft delete a UserAudit by setting the is_deleted flag to true
//	@Tags			User Audit
//	@Param			id		query		int	true	"UserAudit ID"
//	@Success		200		{object}	map[string]interface{}	"UserAudit deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid UserAudit ID"
//	@Failure		500		{object}	map[string]interface{}	"Failed to delete UserAudit"
//	@Router			/fyc/UserAudit [delete]
func DeleteUserAuditAPI(c *gin.Context) {
	idStr := c.Query("id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Msg("No UserAudit ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "UserAudit ID must be provided",
			"code":    12,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid UserAudit ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	log.Info().Int("UserAudit_id", id).Msg("Attempting to soft delete UserAudit")

	rowsAffected, err := DeleteUserAudit(ctx, id)
	if err != nil {
		log.Err(err).Int("UserAudit_id", id).Msg("Failed to soft delete UserAudit")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete UserAudit",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No UserAudit found with the specified ID ------  affected rows 0 ",
			"code":    9,
		})
		return
	}
	log.Info().Int("UserAudit_id", id).Msg("UserAudit deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "UserAudit deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})

}
