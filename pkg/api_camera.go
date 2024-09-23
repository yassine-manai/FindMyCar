package pkg

import (
	"context"
	"fmc/functions"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetCameraAPI godoc
//
//	@Summary		Get cameras or specific camera by ID
//	@Description	Get a list of cameras or a specific camera by ID with optional extra data
//	@Tags			Cameras
//	@Produce		json
//	@Param			id		query		string	false	"Camera ID"
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{object}	Camera		"List of cameras or a single camera"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No cameras found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid camera ID"
//	@Router			/fyc/cameras [get]
func GetCameraAPI(c *gin.Context) {
	log.Debug().Msg("GetCameraAPI request")
	ctx := context.Background()
	idStr := c.Query("id")
	extraReq := c.DefaultQuery("extra", "false")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		log.Info().Str("camera_id", idStr).Msg("Fetching camera by ID")

		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid camera ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		// Fetch camera by ID
		camera, err := GetCameraByID(ctx, id)
		if err != nil {
			log.Err(err).Str("camera_id", idStr).Msg("Error retrieving camera by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Camera not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("camera_id", idStr).Msg("Camera fetched successfully")
		c.JSON(http.StatusOK, camera)
		return
	}

	log.Info().Str("extra", extraReq).Msg("Fetching all cameras")

	if strings.ToLower(extraReq) == "true" || strings.ToLower(extraReq) == "1" || strings.ToLower(extraReq) == "yes" {

		log.Debug().Msg("Fetching cameras with extra data")
		cameras, err := GetAllCameraExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all cameras with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all cameras with extra data",
				"code":    10,
			})
			return
		}

		if len(cameras) == 0 {
			log.Info().Interface("camera_list", cameras).Msg("No cameras found with extra data")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No cameras found",
				"code":    9,
			})
			return
		}

		log.Info().Int("camera_count", len(cameras)).Msg("Cameras fetched successfully")
		c.JSON(http.StatusOK, cameras)
		return
	}

	cameras, err := GetAllCamera(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all cameras")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all cameras",
			"code":    10,
		})
		return
	}

	if len(cameras) == 0 {
		log.Info().Msg("No cameras found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No cameras found",
			"code":    9,
		})
		return
	}

	log.Info().Int("camera_count", len(cameras)).Msg("Cameras fetched successfully")
	c.JSON(http.StatusOK, cameras)
}

// CreateCamera godoc
//
//	@Summary		Add a new camera
//	@Description	Add a new camera to the database
//	@Tags			Cameras
//	@Accept			json
//	@Produce		json
//	@Param			camera	body		Camera	true	"Camera data"
//	@Success		201		{object}	Camera	"Camera created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload"
//	@Failure		500		{object}	map[string]interface{}	"Failed to create a new camera"
//	@Router			/fyc/cameras [post]
func CreateCameraAPI(c *gin.Context) {
	var camera Camera

	if err := c.ShouldBindJSON(&camera); err != nil {
		log.Err(err).Msg("Invalid request payload for camera creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Msg("Creating new camera")
	if !functions.Contains(Zonelist, *camera.ZoneIdIn) {
		*camera.ZoneIdIn = 0
	}

	if !functions.Contains(Zonelist, *camera.ZoneIdOut) {
		*camera.ZoneIdOut = 0
	}

	ctx := context.Background()
	if err := CreateCamera(ctx, &camera); err != nil {
		log.Err(err).Msg("Error creating new camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("camera_id", camera.ID).Msg("Camera created successfully")
	LoadCameralist()
	c.JSON(http.StatusCreated, camera)
}

// UpdateCamera godoc
//
//	@Summary		Update a camera by ID
//	@Description	Update an existing camera by ID
//	@Tags			Cameras
//	@Accept			json
//	@Produce		json
//	@Param			id		query		int			true	"Camera ID"
//	@Param			camera	body		Camera		true	"Updated camera data"
//	@Success		200		{object}	map[string]interface{}	"Camera updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request payload or ID mismatch"
//	@Failure		404		{object}	map[string]interface{}	"Camera not found"
//	@Failure		500		{object}	map[string]interface{}	"Failed to update camera"
//	@Router			/fyc/cameras [put]
func UpdateCameraAPI(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	log.Info().Str("camera_id", idStr).Msg("Updating camera")

	if err != nil {
		log.Error().Str("camera_id", idStr).Msg("Invalid ID format for camera update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates Camera

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for camera update")
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

	ctx := context.Background()
	rowsAffected, err := UpdateCamera(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Int64("Rows Affected", rowsAffected).Msg("No camera found with the specified ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No camera found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("camera_id", idStr).Msg("Camera updated successfully")
	LoadCameralist()
	c.JSON(http.StatusOK, gin.H{
		"message":       "Camera updated successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteCamera godoc
//
//	@Summary		Delete a camera by ID
//	@Description	Delete a camera by ID
//	@Tags			Cameras
//	@Param			id	query		int		true	"Camera ID"
//	@Success		200	{object}	map[string]interface{}	"Camera deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Invalid ID format"
//	@Failure		404	{object}	map[string]interface{}	"Camera not found"
//	@Failure		500	{object}	map[string]interface{}	"Failed to delete camera"
//	@Router			/fyc/cameras [delete]
func DeleteCameraAPI(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)

	log.Info().Str("camera_id", idStr).Msg("Deleting camera")

	if err != nil {
		log.Error().Str("camera_id", idStr).Msg("Invalid ID format for camera deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := DeleteCamera(ctx, id)
	if err != nil {
		log.Err(err).Str("camera_id", idStr).Msg("Error deleting camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Msg("No camera found with the specified ID")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No camera found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("camera_id", idStr).Msg("Camera deleted successfully")
	LoadCameralist()
	c.JSON(http.StatusOK, gin.H{
		"success":      "Camera deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
