package pkg

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type CarDetailAPI struct {
	CarDetailService *CarDetailOp
}

func NewCarDetailAPI(db *bun.DB) *CarDetailAPI {
	return &CarDetailAPI{
		CarDetailService: NewCarDetail(db),
	}
}

// GetCarDetails godoc
//
//	@Summary		Get all car details
//	@Description	Get a list of all car details
//	@Tags			Car Details
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{array}		CarDetail
//	@Router			/fyc/carDetails [get]
func (api *CarDetailAPI) GetCarDetails(c *gin.Context) {
	ctx := context.Background()
	extraReq := c.DefaultQuery("extra", "false")

	if strings.ToLower(extraReq) == "true" || strings.ToLower(extraReq) == "1" || strings.ToLower(extraReq) == "yes" {
		carDetails, err := api.CarDetailService.GetAllCarDetailExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all car details with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all car details with extra data",
				"code":    10,
			})
			return
		}

		if len(carDetails) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No car details found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, carDetails)
		return
	}

	carDetails, err := api.CarDetailService.GetAllCarDetail(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all car details")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all car details",
			"code":    10,
		})
		return
	}

	if len(carDetails) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No car details found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, carDetails)
}

// CreateCarDetail godoc
//
//	@Summary		Add a new car detail
//	@Description	Add a new car detail to the database
//	@Tags			Car Details
//	@Accept			json
//	@Produce		json
//	@Param			CarDetail	body		CarDetail	true	"Car detail data"
//	@Success		201		{object}	CarDetail
//	@Router			/fyc/carDetails [post]
func (api *CarDetailAPI) CreateCarDetail(c *gin.Context) {
	var carDetail CarDetail

	if err := c.ShouldBindJSON(&carDetail); err != nil {
		log.Err(err).Msg("Invalid request payload for car detail creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := api.CarDetailService.CreateCarDetail(ctx, &carDetail); err != nil {
		log.Err(err).Msg("Error creating new car detail")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create a new car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	c.JSON(http.StatusCreated, carDetail)
}

// UpdateCarDetailById godoc
//
//	@Summary		Update a car detail by ID
//	@Description	Update an existing car detail by ID
//	@Tags			Car Details
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int			true	"Car ID"
//	@Param			CarDetail	body		CarDetail	true	"Updated car detail data"
//	@Success		200		{object}	CarDetail
//	@Failure		400		{object}	map[string]interface{}	"Invalid request"
//	@Failure		404		{object}	map[string]interface{}	"Car detail not found"
//	@Router			/fyc/carDetails/{id} [put]
func (api *CarDetailAPI) UpdateCarDetailById(c *gin.Context) {
	// Convert ID param to integer
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid ID format for car detail update")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	var updates CarDetail

	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Err(err).Msg("Invalid request payload for car detail update")

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

	// Call the service to update the car detail
	ctx := context.Background()
	rowsAffected, err := api.CarDetailService.UpdateCarDetail(ctx, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating car detail by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No car detail found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Car detail modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// DeleteCarDetail godoc
//
//	@Summary		Delete a car detail
//	@Description	Delete a car detail by ID
//	@Tags			Car Details
//	@Param			id	path		int		true	"Car detail ID"
//	@Success		200	{object}	map[string]interface{}	"Car detail deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Invalid request"
//	@Failure		404		{object}	map[string]interface{}	"Car detail not found"
//	@Router			/fyc/carDetails/{id} [delete]
func (api *CarDetailAPI) DeleteCarDetail(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Msg("Error ID Format")

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := api.CarDetailService.DeleteCarDetail(ctx, id)
	if err != nil {
		log.Err(err).Msg("Error deleting car detail")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete car detail",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No car detail found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      "Car detail deleted successfully",
		"rowsAffected": rowsAffected,
		"code":         8,
	})
}
