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

// GetZonesImages godoc
//
//	@Summary		Get all zones images
//	@Description	Get a list of all zones images
//	@Tags			Zones Image
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200	{array}		ImageZone
//	@Router			/fyc/zonesImages [get]
func GetAllImageZonesAPI(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		//zoneimage return list of object
		zonesImage, err := GetAllZoneImageExtra(ctx)
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
				"message": "No zones found",
				"code":    9,
			})
			return
		}

		for i := range zonesImage {
			if zonesImage[i].ImageSm != "" {
				zonesImage[i].ImageSm = functions.ByteaToBase64([]byte(zonesImage[i].ImageSm))
			}

			if zonesImage[i].ImageLg != "" {
				zonesImage[i].ImageLg = functions.ByteaToBase64([]byte(zonesImage[i].ImageLg))
			}
		}

		c.JSON(http.StatusOK, zonesImage)
		return
	}

	zoimg, err := GetAllZoneImage(ctx)
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

	for i := range zoimg {
		if zoimg[i].ImageSm != "" {
			zoimg[i].ImageSm = functions.ByteaToBase64([]byte(zoimg[i].ImageSm))
		}

		if zoimg[i].ImageLg != "" {
			zoimg[i].ImageLg = functions.ByteaToBase64([]byte(zoimg[i].ImageLg))
		}
	}

	c.JSON(http.StatusOK, zoimg)
}

// GetZoneImageByIDAPI godoc
//
//	@Summary		Get zoneimage by field
//	@Description	Get a specific zoneimage by either ID or Zone
//	@Tags			Zones Image
//	@Produce		json
//	@Param			field	query		string	true	"Search by field (id or zone)"
//	@Param			value	query		string	true	"Value of the selected field"
//	@Success		200	{object}	ImageZone
//	@Router			/fyc/zonesImage [get]
func GetZoneImageByIDAPI(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")

	// Validate the field input
	if field != "id" && field != "zone" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid field",
			"message": "Field must be either 'id' or 'zone'",
			"code":    11,
		})
		return
	}

	// Validate the value parameter
	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing value",
			"message": "Value parameter is required",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	var zoneImage *ImageZone
	//var err error

	// Query based on the selected field (either id or zone)
	if field == "id" {
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Err(err).Str("id", value).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "Zone ID must be a valid integer",
				"code":    13,
			})
			return
		}
		zoneImage, _ = GetZoneImageByID(ctx, id)
		c.JSON(http.StatusOK, zoneImage)

	} else if field == "zone" {
		var zone []ImageZone

		zone_id, err := strconv.Atoi(value)
		if err != nil {
			log.Err(err).Str("id", value).Msg("Invalid Zone ID format")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid ID format",
				"message": "Zone ID must be a valid integer",
				"code":    13,
			})
			return
		}

		zone, _ = GetZoneImageByZoneID(ctx, zone_id)

		if len(zone) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No zones images found",
				"code":    9,
			})
			return
		}

		for i := range zone {
			if zone[i].ImageSm != "" {
				zone[i].ImageSm = functions.ByteaToBase64([]byte(zone[i].ImageSm))
			}

			if zone[i].ImageLg != "" {
				zone[i].ImageLg = functions.ByteaToBase64([]byte(zone[i].ImageLg))
			}
		}

		c.JSON(http.StatusOK, zone)
	}

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
func CreateZoneImageAPI(c *gin.Context) {
	var zoneImage ImageZone

	if err := c.ShouldBindJSON(&zoneImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(Zonelist, *zoneImage.ZoneID) {
		log.Debug().Msg("Zone not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *zoneImage.ZoneID),
			"code":    9,
		})
		return
	}

	ImageSmEnc, err := functions.DecodeBase64ToByteArray(zoneImage.ImageSm)
	if err != nil {
		log.Err(err).Msg("Error converting image SM")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image SM",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ImageLgEnc, err := functions.DecodeBase64ToByteArray(zoneImage.ImageLg)
	if err != nil {
		log.Err(err).Msg("Error converting image LG")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	zoneImage.ImageSm = string(ImageSmEnc)
	fmt.Println("-*-*-*-*-*************************************--------------------------")
	var a = zoneImage.ImageSm

	var byteArr [64]byte
	// Fill byteArr with some data
	copy(byteArr[:], []byte(a))

	// Convert [64]byte to a string by slicing the array first
	aString := string(byteArr[:])
	fmt.Println(copy(byteArr[:], []byte(a)))

	fmt.Println("-*-*-*-*-*************************************--------------------------")

	fmt.Print(aString)
	//zoneImage.ImageLg = ImageLgEnc

	//fmt.Println(string(ImageSmEnc))
	//fmt.Println(ImageLgEnc)

	fmt.Println("-*-*-*-*-*************************************--------------------------")

	ctx := context.Background()

	if err := CreateZoneImage(ctx, &zoneImage); err != nil {
		log.Err(err).Msg("Error creating new zone image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new zone image",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	fmt.Println(zoneImage.ImageSm)
	fmt.Println(zoneImage.ImageLg)
	fmt.Print(ImageLgEnc)
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
func UpdateZoneImageByIdAPI(c *gin.Context) {
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

	if *updates.ZoneID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the param ID",
			"code":    13,
		})
		return
	}

	if !functions.Contains(Zonelist, *updates.ZoneID) {
		log.Debug().Msg("Zone not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", *updates.ZoneID),
			"code":    9,
		})
		return
	}

	ImageSmEnc, err := functions.Base64ToBytea(updates.ImageSm)
	if err != nil {
		log.Err(err).Msg("Error converting image SM")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image SM",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ImageLgEnc, err := functions.Base64ToBytea(updates.ImageLg)
	if err != nil {
		log.Err(err).Msg("Error converting image LG")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	updates.ImageSm = string(ImageSmEnc)
	updates.ImageLg = string(ImageLgEnc)

	// Call the service to update the present car
	rowsAffected, err := UpdateZoneImage(ctx, id, &updates)
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
func DeleteZoneImageAPI(c *gin.Context) {

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
	rowsAffected, err := DeleteZoneImage(ctx, id)
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
