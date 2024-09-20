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

// GetCarDetails godoc
//
//	@Summary		Get all car details
//	@Description	Get a list of all car details
//	@Tags			Car Details
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{array}		CarDetail
//	@Router			/fyc/carDetails [get]
func GetCarDetailsAPI(c *gin.Context) {
	ctx := context.Background()
	extraReq := c.DefaultQuery("extra", "false")

	if strings.ToLower(extraReq) == "true" || strings.ToLower(extraReq) == "1" || strings.ToLower(extraReq) == "yes" {
		carDetailsExtra, err := GetAllCarDetailExtra(ctx)
		if err != nil {
			log.Err(err).Msg("Error getting all car details with extra data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all car details with extra data",
				"code":    10,
			})
			return
		}

		if len(carDetailsExtra) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No car details found",
				"code":    9,
			})
			return
		}

		for i := range carDetailsExtra {
			if carDetailsExtra[i].Image1 != "" {
				carDetailsExtra[i].Image1 = functions.ByteaToBase64([]byte(carDetailsExtra[i].Image1))
			}

			if carDetailsExtra[i].Image2 != "" {
				carDetailsExtra[i].Image2 = functions.ByteaToBase64([]byte(carDetailsExtra[i].Image2))
			}

			if carDetailsExtra[i].Image3 != "" {
				carDetailsExtra[i].Image3 = functions.ByteaToBase64([]byte(carDetailsExtra[i].Image3))
			}
		}

		c.JSON(http.StatusOK, carDetailsExtra)
		return
	}

	carDetails, err := GetAllCarDetail(ctx)
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

	for i := range carDetails {
		if carDetails[i].Image1 != "" {
			carDetails[i].Image1 = functions.ByteaToBase64([]byte(carDetails[i].Image1))
		}

		if carDetails[i].Image2 != "" {
			carDetails[i].Image2 = functions.ByteaToBase64([]byte(carDetails[i].Image2))
		}

		if carDetails[i].Image3 != "" {
			carDetails[i].Image3 = functions.ByteaToBase64([]byte(carDetails[i].Image3))
		}
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
func CreateCarDetailAPI(c *gin.Context) {
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

	Image1Enc, err := functions.Base64ToBytea(carDetail.Image1)
	if err != nil {
		log.Err(err).Msg("Error converting image 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 1",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image2Enc, err := functions.Base64ToBytea(carDetail.Image2)
	if err != nil {
		log.Err(err).Msg("Error converting image 2")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 2",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image3Enc, err := functions.Base64ToBytea(carDetail.Image3)
	if err != nil {
		log.Err(err).Msg("Error converting image 3")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 3",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	carDetail.Image1 = string(Image1Enc)
	carDetail.Image2 = string(Image2Enc)
	carDetail.Image3 = string(Image3Enc)

	log.Debug().Msg(carDetail.Image1)
	log.Debug().Msg(carDetail.Image2)
	log.Debug().Msg(carDetail.Image3)
	fmt.Scan()

	ctx := context.Background()
	if err := CreateCarDetail(ctx, &carDetail); err != nil {
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

// GetCarDetailByID godoc
//
//	@Summary		Get cardetail by ID
//	@Description	Get a specific carDetail by ID
//	@Tags			Car Details
//	@Produce		json
//	@Param			id	path		int	true	"CarDetail ID"
//	@Success		200	{object}	CarDetail
//	@Router			/fyc/carDetails/{id} [get]
func GetCarDetailsByIdAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Invalid carDetail ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "carDetail ID must be a valid integer",
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	carDetail, err := GetCarDetailByID(ctx, id)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error retrieving carDetail by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "CarDetail not found",
			"code":    9,
		})
		return
	}

	carDetail.Image1 = functions.ByteaToBase64([]byte(carDetail.Image1))
	carDetail.Image2 = functions.ByteaToBase64([]byte(carDetail.Image2))
	carDetail.Image3 = functions.ByteaToBase64([]byte(carDetail.Image3))

	c.JSON(http.StatusOK, carDetail)
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
func UpdateCarDetailByIdAPI(c *gin.Context) {
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

	Image1Enc, err := functions.Base64ToBytea(updates.Image1)
	if err != nil {
		log.Err(err).Msg("Error converting image 1")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 1",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image2Enc, err := functions.Base64ToBytea(updates.Image2)
	if err != nil {
		log.Err(err).Msg("Error converting image 2")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 2",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	Image3Enc, err := functions.Base64ToBytea(updates.Image3)
	if err != nil {
		log.Err(err).Msg("Error converting image 3")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error converting image 3",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	updates.Image1 = string(Image1Enc)
	updates.Image2 = string(Image2Enc)
	updates.Image3 = string(Image3Enc)

	// Call the service to update the car detail
	ctx := context.Background()
	rowsAffected, err := UpdateCarDetail(ctx, id, &updates)
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
func DeleteCarDetailAPI(c *gin.Context) {

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
	rowsAffected, err := DeleteCarDetail(ctx, id)
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
