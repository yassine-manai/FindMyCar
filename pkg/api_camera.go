package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"fmc/functions"
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

// GetCameraAPI godoc
//
//	@Summary		Get enabled cameras or a specific camera by ID
//	@Description	Get a list of enabled cameras or a specific camera by ID with optional extra data
//	@Tags			Cameras
//	@Produce		json
//	@Param			id		query		string	false	"Camera ID"
//	@Success		200		{object}	Camera		"List of enabled cameras or a single camera"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No cameras found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid camera ID"
//	@Router			/fyc/camerasEnabled [get]
func GetCameraEnabledAPI(c *gin.Context) {
	log.Debug().Msg("GetCameraEnabledAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid camera ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		camera, err := GetCameraEnabledByID(ctx, id)
		if err != nil {
			log.Err(err).Str("camera_id", idStr).Msg("Error retrieving camera by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Camera not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("camera_id", idStr).Msg("Enabled camera fetched successfully")
		c.JSON(http.StatusOK, camera)
		return
	}

	// Fetch all enabled cameras
	cameras, err := GetCameraListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled cameras")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled cameras",
			"code":    10,
		})
		return
	}

	if len(cameras) == 0 {
		log.Info().Msg("No enabled cameras found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled cameras found",
			"code":    9,
		})
		return
	}

	log.Info().Int("camera_count", len(cameras)).Msg("Enabled cameras fetched successfully")
	c.JSON(http.StatusOK, cameras)
}

// GetCameraAPI godoc
//
//	@Summary		Get deleted cameras or a specific camera by ID
//	@Description	Get a list of deleted cameras or a specific camera by ID with optional extra data
//	@Tags			Cameras
//	@Produce		json
//	@Param			id		query		string	false	"Camera ID"
//	@Success		200		{object}	Camera		"List of deleted cameras or a single camera"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No cameras found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid camera ID"
//	@Router			/fyc/camerasDeleted [get]
func GetCameraDeletedAPI(c *gin.Context) {
	log.Debug().Msg("GetCameraDeletedAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid camera ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		camera, err := GetCameraDeletedByID(ctx, id)
		if err != nil {
			log.Err(err).Str("camera_id", idStr).Msg("Error retrieving camera by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Camera not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("camera_id", idStr).Msg("Deleted camera fetched successfully")
		c.JSON(http.StatusOK, camera)
		return
	}

	// Fetch all deleted cameras
	cameras, err := GetCameraListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted cameras")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted cameras",
			"code":    10,
		})
		return
	}

	if len(cameras) == 0 {
		log.Info().Msg("No deleted cameras found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted cameras found",
			"code":    9,
		})
		return
	}

	log.Info().Int("camera_count", len(cameras)).Msg("Deleted cameras fetched successfully")
	c.JSON(http.StatusOK, cameras)
}

// ChangeStateAPI godoc
//
//	@Summary		Change camera state or retrieve cameras by ID
//	@Description	Change the state of a camera (e.g., enabled/disabled) or retrieve a camera by ID
//	@Tags			Cameras
//	@Produce		json
//	@Param			state	query		bool	false	"Camera State"
//	@Param			id		query		int 	false	"Camera ID"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Failure		404		{object}	map[string]interface{}	"No cameras found"
//	@Failure		400		{object}	map[string]interface{}	"Bad request: Invalid camera ID or state"
//	@Router			/fyc/cameraState [put]
func ChangeCameraStateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid camera ID format")
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

	rowsAffected, err := ChangeCameraState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("camera with id %d is already enabled", id) {
			log.Info().Str("camera_id", idStr).Msg("Camera is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("camera_id", idStr).Msg("Error changing camera state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error changing camera state",
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("camera_id", idStr).Msg("Camera not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Camera not found or state unchanged",
			"code":    9,
		})
		return
	}

	LoadCameralist()
	log.Info().Str("camera_id", idStr).Bool("state", state).Msg("Camera state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Camera state changed successfully",
		"rowsAffected": rowsAffected,
	})
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
	log.Debug().Msg("CreateCameraAPI request")
	ctx := context.Background()
	var newCam Camera

	if err := c.ShouldBindJSON(&newCam); err != nil {
		log.Err(err).Msg("Invalid input for new camera")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	log.Info().Msg("Creating new camera")
	if !functions.Contains(Zonelist, *newCam.ZoneIdIn) {
		*newCam.ZoneIdIn = 0
	}

	if !functions.Contains(Zonelist, *newCam.ZoneIdOut) {
		*newCam.ZoneIdOut = 0
	}

	if err := CreateCamera(ctx, &newCam); err != nil {
		log.Err(err).Msg("Error creating new camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("camera_id", newCam.ID).Msg("Camera created successfully")
	LoadCameralist()
	c.JSON(http.StatusCreated, newCam)
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
	var updates Camera
	ctx := context.Background()
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

// DeleteCameraAPI godoc
//
//	@Summary		Soft delete a camera
//	@Description	Soft delete a camera by setting the is_deleted flag to true
//	@Tags			Cameras
//	@Param			id		query		string	true	"Camera ID"
//	@Success		200		{object}	map[string]interface{}	"Camera deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid camera ID"
//	@Failure		500		{object}	map[string]interface{}	"Failed to delete camera"
//	@Router			/fyc/cameras [delete]
func DeleteCameraAPI(c *gin.Context) {
	idStr := c.Query("id")
	ctx := context.Background()

	if idStr == "" {
		log.Error().Msg("No camera ID provided for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "Camera ID must be provided",
			"code":    12,
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid camera ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	log.Info().Int("camera_id", id).Msg("Attempting to soft delete camera")

	_, err = DeleteCamera(ctx, id)
	if err != nil {
		log.Err(err).Int("camera_id", id).Msg("Failed to soft delete camera")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete camera",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	LoadCameralist()
	log.Info().Int("camera_id", id).Msg("Camera deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Camera deleted successfully",
		"code":    0,
	})
}
