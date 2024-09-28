package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fmc/functions"
)

// GetsignAPI godoc
//
//	@Summary		Get sign or specific sign by ID
//	@Description	Get a list of sign or a specific sign by ID with optional extra data
//	@Tags			Sign
//	@Produce		json
//	@Param			id		query		int	false	"sign ID"
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{object}	Sign		"List of sign or a single sign"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No sign found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid sign ID"
//	@Router			/fyc/sign [get]
func GetSignAPI(c *gin.Context) {
	log.Debug().Msg("GetsignAPI request")
	ctx := context.Background()
	idStr := c.Query("id")
	extraReq := c.DefaultQuery("extra", "false")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		log.Info().Str("sign_id", idStr).Msg("Fetching sign by ID")

		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid sign ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		sign, err := GetSignById(ctx, id)
		if err != nil {
			log.Err(err).Str("sign_id", idStr).Msg("Error retrieving Sign by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Sign not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("sign_id", idStr).Msg("sign fetched successfully")
		c.JSON(http.StatusOK, sign)
		return
	}

	log.Info().Str("extra", extraReq).Msg("Fetching all sign")

	sign, err := GetAllSigns(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all sign",
			"code":    10,
		})
		return
	}

	if len(sign) == 0 {
		log.Info().Msg("No sign found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No sign found",
			"code":    9,
		})
		return
	}

	log.Info().Int("sign_count", len(sign)).Msg("Sign fetched successfully")
	c.JSON(http.StatusOK, sign)
}

// GetsignAPI godoc
//
//	@Summary		Get enabled sign or a specific sign by ID
//	@Description	Get a list of enabled sign or a specific sign by ID with optional extra data
//	@Tags			Sign
//	@Produce		json
//	@Param			id		query		int	false	"sign ID"
//	@Success		200		{object}	Sign		"List of enabled sign or a single sign"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No sign found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid sign ID"
//	@Router			/fyc/signEnabled [get]
func GetSignEnabledAPI(c *gin.Context) {
	log.Debug().Msg("GetsignEnabledAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid sign ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		sign, err := GetEnabledSignByID(ctx, id)
		if err != nil {
			log.Err(err).Str("sign_id", idStr).Msg("Error retrieving sign by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "sign not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("sign_id", idStr).Msg("Enabled sign fetched successfully")
		c.JSON(http.StatusOK, sign)
		return
	}

	// Fetch all enabled sign
	sign, err := GetSignListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled sign",
			"code":    10,
		})
		return
	}

	if len(sign) == 0 {
		log.Info().Msg("No enabled sign found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled sign found",
			"code":    9,
		})
		return
	}

	log.Info().Int("sign_count", len(sign)).Msg("Enabled sign fetched successfully")
	c.JSON(http.StatusOK, sign)
}

// GetsignAPI godoc
//
//	@Summary		Get deleted sign or a specific sign by ID
//	@Description	Get a list of deleted sign or a specific sign by ID with optional extra data
//	@Tags			Sign
//	@Produce		json
//	@Param			id		query		int	false	"sign ID"
//	@Success		200		{object}	Sign		"List of deleted sign or a single sign"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No sign found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid sign ID"
//	@Router			/fyc/signDeleted [get]
func GetSignDeletedAPI(c *gin.Context) {
	log.Debug().Msg("GetsignDeletedAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid sign ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		sign, err := GetDeletedSignByID(ctx, id)
		if err != nil {
			log.Err(err).Str("sign_id", idStr).Msg("Error retrieving sign by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "sign not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("sign_id", idStr).Msg("Deleted sign fetched successfully")
		c.JSON(http.StatusOK, sign)
		return
	}

	// Fetch all deleted sign
	sign, err := GetSignListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted sign",
			"code":    10,
		})
		return
	}

	if len(sign) == 0 {
		log.Info().Msg("No deleted sign found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted sign found",
			"code":    9,
		})
		return
	}

	log.Info().Int("sign_count", len(sign)).Msg("Deleted sign fetched successfully")
	c.JSON(http.StatusOK, sign)
}

// ChangeStateAPI godoc
//
//	@Summary		Change sign state or retrieve sign by ID
//	@Description	Change the state of a sign (e.g., enabled/disabled) or retrieve a sign by ID
//	@Tags			Sign
//	@Produce		json
//	@Param			state	query		bool	false	"sign State"
//	@Param			id		query		int 	false	"sign ID"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No sign found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid sign ID or state"
//	@Router			/fyc/signState [put]
func ChangeSigntateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid sign ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	stateStr := c.Query("state")
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		log.Err(err).Str("state", stateStr).Msg("Invalid state format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid state format",
			"message": "State must be a boolean value (true/false)",
			"code":    13,
		})
		return
	}

	rowsAffected, err := ChangeSignState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("sign with id %d is already enabled", id) {
			log.Info().Str("sign_id", idStr).Msg("sign is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("sign_id", idStr).Msg("Error changing sign state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error changing sign state",
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("sign_id", idStr).Msg("sign not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "sign not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("sign_id", idStr).Bool("state", state).Msg("sign state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "sign state changed successfully",
		"rowsAffected": rowsAffected,
	})
}

// Createsign godoc
//
//	@Summary		Add a new sign
//	@Description	Add a new sign to the database
//	@Tags			Sign
//	@Accept			json
//	@Produce		json
//	@Param			sign	body		Sign	true	"sign data"
//	@Success		201		{object}	Sign	"sign created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create a new sign"
//	@Router			/fyc/sign [post]
func CreateSignAPI(c *gin.Context) {
	log.Debug().Msg("CreatesignAPI request")
	ctx := context.Background()
	var newSign Sign

	if err := c.ShouldBindJSON(&newSign); err != nil {
		log.Err(err).Msg("Invalid input for new sign")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Msg("Creating new sign")
	if !functions.Contains(Zonelist, newSign.ZoneID) {
		newSign.ZoneID = 0
	}

	if err := CreateSign(ctx, &newSign); err != nil {
		log.Err(err).Msg("Error creating new sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new sign",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("sign_id", newSign.ID).Msg("Sign created successfully")
	c.JSON(http.StatusCreated, newSign)
}

// Updatesign godoc
//
//	@Summary		Update a sign by ID
//	@Description	Update an existing sign by ID
//	@Tags			Sign
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int			true	"sign ID"
//	@Param			sign	body		Sign		true	"Updated sign data"
//	@Success		200		{object}	map[string]interface{}	"sign updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload or ID mismatch"
//	@Failure		404		{object}	map[string]interface{}	"sign not found"
//	@Failure		500		{object}	map[string]interface{}	"Failed to update sign"
//	@Router			/fyc/sign [put]
func UpdateSignAPI(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	var updates Sign
	ctx := context.Background()
	log.Info().Str("sign_id", idStr).Msg("Updating sign")

	if err != nil {
		log.Error().Str("sign_id", idStr).Msg("Invalid ID format for sign update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for sign update")
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

	log.Info().Msg("Updating Sign")
	if !functions.Contains(Zonelist, updates.ZoneID) {
		updates.ZoneID = 0
	}

	rowsAffected, err := UpdateSign(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update sign",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("sign_id", idStr).Int64("Rows Affected", rowsAffected).Msg("No sign found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No sign found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("sign_id", idStr).Msg("sign updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":       "sign updated successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeletesignAPI godoc
//
//	@Summary		Soft delete a sign
//	@Description	Soft delete a sign by setting the is_deleted flag to true
//	@Tags			Sign
//	@Param			id		query		int	true	"sign ID"
//	@Success		200		{object}	map[string]interface{}	"sign deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid sign ID"
//	@Failure		500		{object}	map[string]interface{}	"Failed to delete sign"
//	@Router			/fyc/sign [delete]
func DeleteSignAPI(c *gin.Context) {
	idStr := c.Query("id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Msg("No sign ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "sign ID must be provided",
			"code":    12,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid sign ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	log.Info().Int("sign_id", id).Msg("Attempting to soft delete sign")

	rowsAffected, err := DeleteSign(ctx, id)
	if err != nil {
		log.Err(err).Int("sign_id", id).Msg("Failed to soft delete sign")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete sign",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No Sign found with the specified ID ------  affected rows 0 ",
			"code":    9,
		})
		return
	}
	log.Info().Int("sign_id", id).Msg("sign deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success":      "Sign deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})

}
