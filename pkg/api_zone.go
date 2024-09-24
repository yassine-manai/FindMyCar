package pkg

import (
	"context"
	"fmc/functions"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetZones godoc
//
//	@Summary		Get all zones
//	@Description	Get a list of all zones
//	@Tags			Zones
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200	{array}		Zone
//	@Router			/fyc/zones [get]
func GetZonesAPI(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

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
		log.Debug().Interface("Zones", zones).Msg("Get Zone api db dat")
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

// GetZoneByID godoc
//
//	@Summary		Get zone by ID
//	@Description	Get a specific zone by ID
//	@Tags			Zones
//	@Produce		json
//	@Param			id	path		int	true	"Zone ID"
//	@Success		200	{object}	Zone
//	@Router			/fyc/zones/{id} [get]
func GetZoneByIDAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid Zone ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Zone ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	zone, err := GetZoneByID(ctx, id)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error retrieving Zone by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Zone not found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, zone)
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

	// Check if CarParkID is valid
	if !functions.Contains(CarParkList, *zone.CarParkID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Carpark not found",
			"message": fmt.Sprintf("Carpark with ID %d is not found", *zone.CarParkID),
			"code":    9,
		})
		return
	}

	// Check if ZoneID already exists
	if functions.Contains(Zonelist, *zone.ZoneID) {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Zone already exists",
			"message": fmt.Sprintf("Zone with ID %d already exists", *zone.ZoneID),
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
// @Failure		400		{object}	map[string]interface{}	"Invalid request"
// @Failure		404		{object}	map[string]interface{}	"Zone not found"
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

	if *updates.ZoneID != id {
		log.Warn().Int("expected_id", id).Int("provided_id", *updates.ZoneID).Msg("ID mismatch in update request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the param ID",
			"code":    13,
		})
		return
	}

	// Check if CarParkID is valid
	if !functions.Contains(CarParkList, *updates.CarParkID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Carpark not found",
			"message": fmt.Sprintf("Carpark with ID %d does not exist", *updates.CarParkID),
			"code":    9,
		})
		return
	}

	// Check if ZoneID exists
	if !functions.Contains(Zonelist, *updates.ZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *updates.ZoneID),
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
