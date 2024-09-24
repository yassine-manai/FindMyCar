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

// GetAllCarparks godoc
//
//	@Summary		Get all carparks
//	@Description	Get a list of all carparks
//	@Tags			Carparks
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200	{object}	[]Carpark
//	@Router			/fyc/carparks [get]
func GetAllCarparksAPI(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

	log.Info().Str("extra_req", extra_req).Msg("Fetching all carparks")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		carparkEx, err := GetAllCarparksExtra(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching all carparks with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error fetching all carparks with extra data",
				"code":    10,
			})
			return
		}

		if len(carparkEx) == 0 {
			log.Warn().Msg("No carparks found with extra data")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No carparks found",
				"code":    9,
			})
			return
		}

		log.Info().Int("count", len(carparkEx)).Msg("Carparks with extra data retrieved successfully")
		c.JSON(http.StatusOK, carparkEx)
		return
	}

	carpark, err := GetAllCarparks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching all carparks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error fetching all carparks",
			"code":    10,
		})
		return
	}

	if len(carpark) == 0 {
		log.Warn().Msg("No carparks found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No carparks found",
			"code":    9,
		})
		return
	}

	log.Info().Int("count", len(carpark)).Msg("Carparks retrieved successfully")
	c.JSON(http.StatusOK, carpark)
}

// GetCarparkByID godoc
//
//	@Summary		Get carpark by ID
//	@Description	Get a specific carpark by ID
//	@Tags			Carparks
//	@Produce		json
//	@Param			id	path		int	true	"Carpark ID"
//	@Success		200	{object}	Carpark
//	@Router			/fyc/carparks/{id} [get]
func GetCarparkByIDAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Str("id", idStr).Msg("Invalid carpark ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Carpark ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	carpark, err := GetCarparkByID(ctx, id)
	if err != nil {
		log.Error().Str("id", idStr).Err(err).Msg("Error retrieving carpark by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Carpark not found",
			"code":    9,
		})
		return
	}

	log.Info().Int("id", id).Msg("Carpark retrieved successfully")
	c.JSON(http.StatusOK, carpark)
}

// AddCarpark godoc
//
//	@Summary		Add a new carpark
//	@Description	Add a new carpark to the database
//	@Tags			Carparks
//	@Accept			json
//	@Produce		json
//	@Param			carpark	body		Carpark	true	"Carpark data"
//	@Success		201	{object}	Carpark
//	@Router			/fyc/carparks [post]
func AddCarparkAPI(c *gin.Context) {
	var carpark Carpark
	if err := c.ShouldBindJSON(&carpark); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for carpark creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := AddCarpark(ctx, &carpark); err != nil {
		log.Error().Err(err).Msg("Error creating new carpark")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create carpark",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("id", carpark.ID).Msg("Carpark created successfully")
	LoadCarparklist()
	c.JSON(http.StatusCreated, carpark)
}

// UpdateCarpark godoc
//
//	@Summary		Update a carpark
//	@Description	Update an existing carpark by ID
//	@Tags			Carparks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Carpark ID"
//	@Param			carpark	body		Carpark	true	"Updated carpark data"
//	@Success		200	{object}	Carpark
//	@Router			/fyc/carparks/{id} [put]
func UpdateCarparkAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Str("id", idStr).Msg("Invalid carpark ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var carpark Carpark
	if err := c.ShouldBindJSON(&carpark); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for carpark update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": "Invalid carpark data",
			"code":    12,
		})
		return
	}

	if carpark.ID != id {
		log.Error().Int("body_id", carpark.ID).Int("path_id", id).Msg("Carpark ID mismatch")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the query parameter",
			"code":    13,
		})
		return
	}

	if !functions.Contains(CarParkList, carpark.ID) {
		log.Warn().Int("id", carpark.ID).Msg("Carpark ID not found in the list")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": fmt.Sprintf("Carpark with ID %d does not exist", carpark.ID),
			"code":    9,
		})
		return
	}

	ctx := context.Background()
	_, err = UpdateCarpark(ctx, id, &carpark)
	if err != nil {
		log.Error().Err(err).Msg("Error updating carpark")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update carpark",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Int("id", carpark.ID).Msg("Carpark updated successfully")
	LoadCarparklist()
	c.JSON(http.StatusOK, carpark)
}

// DeleteCarpark godoc
//
//	@Summary		Delete a carpark
//	@Description	Delete a carpark by ID
//	@Tags			Carparks
//	@Param			id	path		int		true	"Carpark ID"
//	@Success		200	{string}	string	"Carpark deleted successfully"
//	@Router			/fyc/carparks/{id} [delete]
func DeleteCarparkAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Str("id", idStr).Msg("Invalid carpark ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := DeleteCarpark(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting carpark")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete carpark",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Int("id", id).Msg("No carpark found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No carpark found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Int("id", id).Msg("Carpark deleted successfully")
	LoadCarparklist()
	c.JSON(http.StatusOK, gin.H{
		"success": "Carpark deleted successfully",
		"message": rowsAffected,
		"code":    8,
	})
}
