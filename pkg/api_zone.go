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

// GetZonesAPI godoc
//
//	@Summary		Get all zones
//	@Description	Get a list of all zones, or a zone by ID if 'id' parameter is provided
//	@Tags			Zones
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Param			id		query		int		false	"Zone ID"
//	@Success		200	{array}		Zone		"List of zones or a single zone"
//	@Router			/fyc/zones [get]
func GetZonesAPI(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		log.Info().Str("Zone ID", idStr).Msg("Fetching zone by ID")

		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		zone, err := GetZoneByID(ctx, id)
		if err != nil {
			log.Err(err).Str("zone id", idStr).Msg("Error retrieving zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("zone_id", idStr).Msg("zone fetched successfully")
		c.JSON(http.StatusOK, zone)
		return
	}

	log.Info().Str("extra", extra_req).Msg("Fetching all cameras")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		zones, err := GetAllZoneExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all zones with extra data ")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all zones with extra data ",
				"code":    10,
			})
			return
		}
		log.Debug().Int("Zones length", len(zones)).Msg("Get Zone api db data")
		if len(zones) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No zones found ",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, zones)
		return
	}

	zoo, err := GetAllZone(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zones")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all zones",
			"code":    10,
		})
		return
	}

	if len(zoo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zones found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, zoo)
}

// GetzonesAPI godoc
//
//	@Summary		Get enabled zones or a specific zone by ID
//	@Description	Get a list of enabled zones or a specific zone by ID with optional extra data
//	@Tags			Zones
//	@Produce		json
//	@Param			id		query		string	false	"Zone ID"
//	@Success		200		{array}	Zone		"List of enabled zones or a single zone"
//	@Router			/fyc/zonesEnabled [get]
func GetZoneEnabledAPI(c *gin.Context) {
	log.Debug().Msg("Get Zone EnabledAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		Zone, err := GetZoneEnabledByID(ctx, id)
		if err != nil {
			log.Err(err).Str("camera_id", idStr).Msg("Error retrieving Zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Zone id", idStr).Msg("Enabled Zone fetched successfully")
		c.JSON(http.StatusOK, Zone)
		return
	}

	// Fetch all enabled cameras
	Zone, err := GetZoneListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled Zones")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled Zones",
			"code":    10,
		})
		return
	}

	if len(Zone) == 0 {
		log.Info().Msg("No enabled Zones found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled Zone found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Zone_count", len(Zone)).Msg("Enabled Zone fetched successfully")
	c.JSON(http.StatusOK, Zone)
}

// GetzonesAPI godoc
//
//	@Summary		Get Deleted zones or a specific zone by ID
//	@Description	Get a list of Deleted zones or a specific zone by ID with optional extra data
//	@Tags			Zones
//	@Produce		json
//	@Param			id		query		string	false	"Zone ID"
//	@Success		200		{object}	Zone		"List of Deleted zones or a single zone"
//	@Router			/fyc/zonesDeleted [get]
func GetZoneDeletedAPI(c *gin.Context) {
	log.Debug().Msg("Get Zone DeletedAPI request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Err(err).Str("id", idStr).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "ID must be a valid integer",
				"code":    12,
			})
			return
		}

		Zone, err := GetZoneDeletedByID(ctx, id)
		if err != nil {
			log.Err(err).Str("Zone_id", idStr).Msg("Error retrieving Zone by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Zone not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Zone_id", idStr).Msg("Deleted Zone fetched successfully")
		c.JSON(http.StatusOK, Zone)
		return
	}

	// Fetch all deleted cameras
	Zone, err := GetZoneListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted Zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted Zone",
			"code":    10,
		})
		return
	}

	if len(Zone) == 0 {
		log.Info().Msg("No deleted Zone found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted Zone found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Zone_count", len(Zone)).Msg("Deleted Zone fetched successfully")
	c.JSON(http.StatusOK, Zone)
}

// ChangeStateAPI godoc
//
//	@Summary		Change Zone state or retrieve Zones by ID
//	@Description	Change the state of a Zone (e.g., enabled/disabled) or retrieve a Zone by ID
//	@Tags			Zones
//	@Produce		json
//	@Param			state	query		bool	false	"Zone State"
//	@Param			id		query		int 	false	"Zone ID"
//	@Success		200		{object}	int		"Number of rows affected by the state change"
//	@Router			/fyc/zoneState [put]
func ChangeZoneStateAPI(c *gin.Context) {
	log.Debug().Msg("Change State API request")
	ctx := context.Background()

	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid zone ID format")
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

	rowsAffected, err := ChangeZoneState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("zone with id %d is already enabled", id) {
			log.Info().Str("zone_id", idStr).Msg("zone is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Error",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("zone_id", idStr).Msg("Error changing zone state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("zone_id", idStr).Msg("Zone not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Zone not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("zone_id", idStr).Bool("state", state).Msg("Zone state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Zone state changed successfully",
		"rowsAffected": rowsAffected,
	})
}

// CreateZoneAPI adds a new zone to the database
// @Summary		Add a new zone
// @Description	Add a new zone to the database
// @Tags			Zones
// @Accept			json
// @Produce		json
// @Param			zone	body		Zone	true	"Zone data"
// @Success		201		{object}	Zone
// @Router			/fyc/zones [post]
func CreateZoneAPI(c *gin.Context) {
	var zone Zone
	ctx := context.Background()

	// Bind the incoming JSON payload to the zone struct
	if err := c.ShouldBindJSON(&zone); err != nil {
		log.Err(err).Msg("Invalid request payload for zone creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	// Check if ZoneID already exists
	if functions.Contains(Zonelist, zone.ZoneID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Zone already exists",
			"message": fmt.Sprintf("Zone with ID %d already exists", zone.ZoneID),
			"code":    9,
		})
		return
	}

	// Create the zone
	if err := CreateZone(ctx, &zone); err != nil {
		log.Err(err).Msg("Error creating new zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	// Return the response with a 201 status
	Loadzonelist()
	c.JSON(http.StatusCreated, zone)
}

// UpdateZoneIdAPI updates a zone by its ID
// @Summary		Update a zone by ID
// @Description	Update an existing zone by ID
// @Tags			Zones
// @Accept			json
// @Produce		json
// @Param			id			path		int			true	"Zone ID"
// @Param			zone		body		Zone		true	"Updated zone data"
// @Success		200		{object}	Zone
// @Router			/fyc/zones/{id} [put]
func UpdateZoneIdAPI(c *gin.Context) {
	// Convert ID param to integer
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Invalid ID format for update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates Zone
	ctx := context.Background()

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for zone update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.ZoneID != id {
		log.Warn().Int("expected_id", id).Int("provided_id", updates.ZoneID).Msg("ID mismatch in update request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the param ID",
			"code":    13,
		})
		return
	}

	// Check if ZoneID exists
	if !functions.Contains(Zonelist, updates.ZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.ZoneID),
			"code":    9,
		})
		return
	}

	// Call the service to update the zone
	rowsAffected, err := UpdateZone(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating zone by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone found with the specified ID",
			"code":    9,
		})
		return
	}
	Loadzonelist()
	c.JSON(http.StatusOK, gin.H{
		"message":       "Zone modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteZoneAPI deletes a zone by its ID
// @Summary		Delete a zone
// @Description	Delete a zone by ID
// @Tags			Zones
// @Param			id	path		int		true	"Zone ID"
// @Success		200	{object}	map[string]interface{}	"Zone deleted successfully"
// @Failure		400		{object}	map[string]interface{}	"Invalid request"
// @Failure		404		{object}	map[string]interface{}	"Zone not found"
// @Router			/fyc/zones/{id} [delete]
func DeleteZoneAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Invalid ID format for deletion")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := DeleteZone(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting Zone")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete Zone",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone found with the specified ID",
			"code":    9,
		})
		return
	}
	Loadzonelist()
	c.JSON(http.StatusOK, gin.H{
		"success":      "Zone deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
