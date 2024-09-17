package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ZoneImageAPI struct {
	ZoneImageService *ImageZoneOp
}

func NewZoneImageAPI(db *bun.DB) *ZoneImageAPI {
	return &ZoneImageAPI{
		ZoneImageService: NewImageZone(db),
	}
}

// GetZonesImages godoc
//
//	@Summary		Get all zones image
//	@Description	Get a list of all zones images
//	@Tags			Zones Image
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200	{array}		ImageZone
//	@Router			/fyc/zonesImage [get]
func (api *ZoneImageAPI) GetImageZones(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		zonesImage, err := api.ZoneImageService.GetAllZoneImageExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all zones image  with extra data ")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all zones image with extra data ",
				"code":    10,
			})
			return
		}

		if len(zonesImage) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No zones found ",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, zonesImage)
		return
	}

	zoimg, err := api.ZoneImageService.GetAllZone(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all zones images")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all zones images",
			"code":    10,
		})
		return
	}

	if len(zoimg) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zones images found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, zoimg)
}

// GetZoneImageByID godoc
//
//	@Summary		Get zoneimage by ID
//	@Description	Get a specific zoneimage by ID
//	@Tags			Zones Image
//	@Produce		json
//	@Param			id	path		int	true	"ZoneImage ID"
//	@Success		200	{object}	ImageZone
//	@Router			/fyc/zonesImage/{id} [get]
func (api *ZoneImageAPI) GetZoneImageByID(c *gin.Context) {
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
	zoneimage, err := api.ZoneImageService.GetZoneImageByID(ctx, id)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error retrieving zoneimage by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Zone Image not found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, zoneimage)
}

// CreateZone godoc
//
//	@Summary		Add a new zone Image
//	@Description	Add a new zone image to the database
//	@Tags			Zones Image
//	@Accept			json
//	@Produce		json
//	@Param			ImageZone	body		ImageZone	true	"Zone image data"
//	@Success		201		{object}	ImageZone
//	@Router			/fyc/zonesImage [post]
func (api *ZoneImageAPI) CreateZoneImage(c *gin.Context) {
	var zoneImage ImageZone
	ctx := context.Background()

	if err := c.ShouldBindJSON(&zoneImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ZoneService := NewZone(api.ZoneImageService.DB)
	_, errup := ZoneService.GetZoneByID(ctx, zoneImage.ZoneID)
	if errup != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", zoneImage.ZoneID),
			"code":    14,
		})
		return
	}

	if err := api.ZoneImageService.CreateZoneImage(ctx, &zoneImage); err != nil {
		log.Err(err).Msg("Error creating new zone image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	c.JSON(http.StatusCreated, zoneImage)
}

// UpdateZoneImageId godoc
//
//	@Summary		Update a zone image by ID
//	@Description	Update an existing zone image by ID
//	@Tags			Zones Image
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int			true	"Zone ID"
//	@Param			Image Zone		body		ImageZone		true	"Updated zone image data"
//	@Success		200		{object}	Zone
//	@Failure		400		{object}	map[string]interface{}	"Invalid request"
//	@Failure		404		{object}	map[string]interface{}	"Zone image not found"
//	@Router			/fyc/zonesImage/{id} [put]
func (api *ZoneImageAPI) UpdateZoneImageById(c *gin.Context) {
	// Convert ID param to integer
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates ImageZone
	ctx := context.Background()

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if updates.ZoneID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the param ID",
			"code":    13,
		})
		return
	}

	ZoneService := NewZone(api.ZoneImageService.DB)
	_, errup := ZoneService.GetZoneByID(ctx, updates.ZoneID)
	if errup != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.ZoneID),
			"code":    14,
		})
		return
	}

	// Call the service to update the present car
	rowsAffected, err := api.ZoneImageService.UpdateZoneImage(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating zone image by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone image found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Zone Image modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteZoneImage godoc
//
//	@Summary		Delete a zone image
//	@Description	Delete a zone image by ID
//	@Tags			Zones Image
//	@Param			id	path		int		true	"Zone image ID"
//	@Success		200	{object}	map[string]interface{}	"Zone image deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request"
//	@Failure		404		{object}	map[string]interface{}	"Zone image not found"
//	@Router			/fyc/zonesImage/{id} [delete]
func (api *ZoneImageAPI) DeleteZoneImage(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := api.ZoneImageService.DeleteZoneImage(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting zone image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No zone image found with the specified ID ------  affected rows 0 ",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      "Zone Image deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
