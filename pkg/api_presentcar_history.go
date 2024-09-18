package pkg

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetHistory godoc
//
//	@Summary		Get all history records
//	@Description	Get a list of all history records
//	@Tags			History
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{array}		PresentCarHistory	"List of history records"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No history records found"
//	@Router			/fyc/history [get]
func GetHistoryAPI(c *gin.Context) {
	extraReq := strings.ToLower(c.DefaultQuery("extra", "false"))

	if extraReq == "true" || extraReq == "1" || extraReq == "yes" {
		log.Info().Msg("Fetching history with extra data")

		ctx := context.Background()
		hist, err := GetAllPresentHistoryExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error fetching history with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all history with extra data",
				"code":    10,
			})
			return
		}

		if len(hist) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No history records found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, hist)
		return
	}

	ctx := context.Background()
	hist, err := GetAllPresentHistory(ctx)
	if err != nil {
		log.Err(err).Msg("Error fetching all history records")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all history",
			"code":    10,
		})
		return
	}

	if len(hist) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No history records found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, hist)
}

// GetHistoryByID godoc
//
//	@Summary		Get history record by ID
//	@Description	Get a specific history record by ID
//	@Tags			History
//	@Produce		json
//	@Param			id	path		int	true	"History record ID"
//	@Success		200	{object}	PresentCarHistory
//	@Failure		400	{object}	map[string]interface{}	"Invalid ID format"
//	@Failure		404	{object}	map[string]interface{}	"History record not found"
//	@Router			/fyc/history/{id} [get]
func GetHistoryByIDAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "History ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	history, err := GetPresentHistoryByID(ctx, id)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error retrieving history by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "History record not found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetHistoryByLPN godoc
//
//	@Summary		Get history record by LPN
//	@Description	Get a specific history record by LPN
//	@Tags			History
//	@Produce		json
//	@Param			lpn	path		string	true	"History record LPN"
//	@Success		200	{object}	PresentCarHistory
//	@Failure		400	{object}	map[string]interface{}	"Invalid LPN format"
//	@Failure		404	{object}	map[string]interface{}	"History record not found"
//	@Router			/fyc/history/{lpn} [get]
func GetHistoryByLPNAPI(c *gin.Context) {
	lpn := c.Param("lpn")
	extraReq := c.Query("extra")

	ctx := context.Background()
	hist, err := GetPresentHistoryByLPN(ctx, lpn)
	if err != nil {
		log.Err(err).Str("lpn", lpn).Msg("Error retrieving history by LPN")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "History record not found",
			"code":    9,
		})
		return
	}

	if extraReq == "yes" {
		c.JSON(http.StatusOK, hist)
	} else {
		response := ResponsePCH{
			ID:         hist.ID,
			LPN:        hist.LPN,
			CurrZoneID: hist.CurrZoneID,
			LastZoneID: hist.LastZoneID,
			CamID:      hist.CamID,
			Image:      hist.Image,
			Confidence: hist.Confidence,
			CameraBody: hist.CameraBody,
		}
		c.JSON(http.StatusOK, response)
	}
}

// CreateHistory godoc
//
//	@Summary		Add a new history record
//	@Description	Add a new history record to the database
//	@Tags			History
//	@Accept			json
//	@Produce		json
//	@Param			history	body		PresentCarHistory	true	"History record data"
//	@Success		201		{object}	PresentCarHistory	"History record created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create a new history record"
//	@Router			/fyc/history [post]
func CreateHistoryAPI(c *gin.Context) {
	var hist PresentCarHistory

	if err := c.ShouldBindJSON(&hist); err != nil {
		log.Err(err).Msg("Invalid request payload for history creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := CreatePresentHistory(ctx, &hist); err != nil {
		log.Err(err).Msg("Error creating new history")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create history",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	c.JSON(http.StatusCreated, hist)
}

// UpdateHistory godoc
//
//	@Summary		Update a history record by ID
//	@Description	Update an existing history record by ID
//	@Tags			History
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"History record ID"
//	@Param			history	body		PresentCarHistory	true	"Updated history record data"
//	@Success		200		{object}	map[string]interface{}	"History record updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload or ID mismatch"
//	@Failure		404		{object}	map[string]interface{}	"History record not found"
//	@Failure		500		{object}	map[string]interface{}	"Failed to update history record"
//	@Router			/fyc/history/{id} [put]
func UpdateHistoryAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format for history update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates PresentCarHistory
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for history update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the param ID",
			"code":    13,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := UpdatePresentHistory(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating history by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update history",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No history found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "History modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteHistory godoc
//
//	@Summary		Delete a history record by ID
//	@Description	Delete a history record by ID
//	@Tags			History
//	@Param			id	path		int		true	"History record ID"
//	@Success		200	{object}	map[string]interface{}	"History record deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid ID format"
//	@Failure		404	{object}	map[string]interface{}	"History record not found"
//	@Failure		500	{object}	map[string]interface{}	"Failed to delete history record"
//	@Router			/fyc/history/{id} [delete]
func DeleteHistoryAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format for history deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := DeletePresentHistory(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting history")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete history",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No history found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      "History deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
