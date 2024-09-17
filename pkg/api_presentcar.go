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
	"github.com/uptrace/bun"
)

// GetAllPresentCars godoc
//
//	@Summary		Get all present cars
//	@Description	Get a list of all present cars
//	@Tags			PresentCars
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//
// @Success		200	{array}		PresentCar
// @Router			/fyc/presentcars [get]
func GetPresentCarsAPI(c *gin.Context, db *bun.DB) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		cars, err := GetAllPresentExtra(ctx, db)
		if err != nil {
			log.Err(err).Msg("Error getting all present cars with extra data ")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "An unexpected error occurred",
				"message": "Error getting all present cars with extra data ",
				"code":    10,
			})
			return
		}

		if len(cars) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "No present cars found",
				"code":    9,
			})
			return
		}

		c.JSON(http.StatusOK, cars)
		return
	}

	Pcars, err := GetAllPresentCars(ctx, db)
	if err != nil {
		log.Err(err).Msg("Error getting all present cars")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all present cars",
			"code":    10,
		})
		return
	}

	if len(Pcars) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present cars found",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, Pcars)
}

// GetPresentCarByLPN godoc
//
//	@Summary		Get present car by LPN
//	@Description	Get a specific present car by LPN
//	@Tags			PresentCars
//	@Produce		json
//	@Param			lpn		path		string	true	"License Plate Number"
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200	{object}	PresentCar
//	@Router			/fyc/presentcars/{lpn} [get]
func GetPresentCarByLPNAPI(c *gin.Context, db *bun.DB) {

	lpn := c.Param("lpn")
	extra_req := c.Query("extra")

	ctx := context.Background()
	car, err := GetPresentCarByLPN(ctx, db, lpn)
	if err != nil {
		log.Err(err).Str("lpn", lpn).Msg("Error retrieving present car by LPN")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Present car not found",
			"code":    9,
		})
		return
	}

	if extra_req == "yes" {
		responseExtra := PresentCar{
			ID:              car.ID,
			CarDetailsID:    car.CarDetailsID,
			CameraID:        car.CameraID,
			Confidence:      car.Confidence,
			CurrZoneID:      car.CurrZoneID,
			LastZoneID:      car.LastZoneID,
			Direction:       car.Direction,
			LPN:             car.LPN,
			TransactionDate: car.TransactionDate,
			Extra:           car.Extra,
		}

		c.JSON(http.StatusOK, responseExtra)

	} else {
		response := ResponsePC{
			ID:              car.ID,
			CarDetailsID:    car.CarDetailsID,
			CameraID:        car.CameraID,
			Confidence:      car.Confidence,
			CurrZoneID:      car.CurrZoneID,
			LastZoneID:      car.LastZoneID,
			Direction:       car.Direction,
			LPN:             car.LPN,
			TransactionDate: car.TransactionDate,
		}

		c.JSON(http.StatusOK, response)

	}

}

// CreatePresentCar godoc
//
//	@Summary		Add a new present car
//	@Description	Add a new present car to the database
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			presentCar	body		PresentCar	true	"Present Car data"
//	@Success		201		{object}	PresentCar
//	@Router			/fyc/presentcars [post]
func CreatePresentCarAPI(c *gin.Context, db *bun.DB) {
	var car PresentCar
	ctx := context.Background()

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(Zonelist, *car.CurrZoneID) {
		log.Debug().Msg("CurrZoneID not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "CurrZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", car.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(Zonelist, *car.LastZoneID) {
		log.Debug().Msg("LastZoneID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "LastZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", car.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(CameraList, *car.CameraID) {
		log.Debug().Msg("Camera ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera ID not found ",
			"message": fmt.Sprintf("Camera with ID %d does not exist", car.CameraID),
			"code":    9,
		})
		return
	}

	// Insert the new car into the database
	if err := CreatePresentCar(ctx, db, &car); err != nil {
		log.Err(err).Msg("Error creating present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	response := PresentCar{
		ID:              car.ID,
		CarDetailsID:    car.CarDetailsID,
		CameraID:        car.CameraID,
		Confidence:      car.Confidence,
		CurrZoneID:      car.CurrZoneID,
		LastZoneID:      car.LastZoneID,
		Direction:       car.Direction,
		LPN:             car.LPN,
		TransactionDate: car.TransactionDate,
		Extra:           car.Extra,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdatePresentCarById godoc
//
//	@Summary		Update a present car by ID
//	@Description	Update an existing present car by ID
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int			true	"Present Car ID"
//	@Param			presentCar	body		PresentCar	true	"Updated present car data"
//	@Success		200		{object}	PresentCar
//	@Router			/fyc/presentcars/{id} [put]
func UpdatePresentCarByIdAPI(c *gin.Context, db *bun.DB) {
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

	var updates PresentCar
	// Bind JSON payload
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if *updates.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the param",
			"code":    13,
		})
		return
	}

	if !functions.Contains(Zonelist, *updates.CurrZoneID) {
		log.Debug().Msg("CurrZoneID not found")

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "CurrZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(Zonelist, *updates.LastZoneID) {
		log.Debug().Msg("LastZoneID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "LastZoneID not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(CameraList, *updates.CameraID) {
		log.Debug().Msg("Camera ID not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera ID not found ",
			"message": fmt.Sprintf("Camera with ID %d does not exist", updates.CameraID),
			"code":    9,
		})
		return
	}

	// Call the service to update the present car
	rowsAffected, err := UpdatePresentCar(ctx, db, id, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating present car by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified ID",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Present car modified successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      updates,
	})
}

// UpdatePresentCarByLpn godoc
//
//	@Summary		Update a present car by LPN
//	@Description	Update an existing present car by lpn
//	@Tags			PresentCars
//	@Accept			json
//	@Produce		json
//	@Param			lpn query string     true  "string default"     default(A)
//	@Param			presentCar	body		PresentCar	true	"Updated present car data by lpn"
//	@Success		200		{object}	PresentCar
//	@Router			/fyc/presentcars [put]
func UpdatePresentCarBylpnAPI(c *gin.Context, db *bun.DB) {

	lpn := c.Query("lpn")
	log.Info().Msgf("provided parameters :%v", lpn)
	var updates PresentCar
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	if !functions.Contains(CarParkList, *updates.CurrZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Carpark not found in CarParkList",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.CurrZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(CameraList, *updates.LastZoneID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Carpark not found in CarParkList",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.LastZoneID),
			"code":    9,
		})
		return
	}

	if !functions.Contains(CarParkList, *updates.CameraID) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Carpark not found in CarParkList",
			"message": fmt.Sprintf("Camera with ID %d does not exist", updates.CameraID),
			"code":    9,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := UpdatePresentCarByLpn(ctx, db, lpn, &updates)
	if err != nil {
		log.Err(err).Msg("Error updating present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified Licence Plate",
			"code":    9,
		})
		return
	}

	response := ResponsePC{
		ID:              updates.ID,
		CarDetailsID:    updates.CarDetailsID,
		CameraID:        updates.CameraID,
		Confidence:      updates.Confidence,
		CurrZoneID:      updates.CurrZoneID,
		LastZoneID:      updates.LastZoneID,
		Direction:       updates.Direction,
		LPN:             lpn,
		TransactionDate: updates.TransactionDate,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Present car modified succesfully successfully",
		"rows_affected": rowsAffected,
		"code":          8,
		"response":      response,
	})

}

// DeletePresentCar godoc
//
//	@Summary		Delete a present car
//	@Description	Delete a present car by ID
//	@Tags			PresentCars
//	@Param			id	path		int		true	"Present Car ID"
//	@Success		200	{string}	string	"Present car deleted successfully"
//	@Success		200			{object}	string	"Success"
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		404			{object}	string	"Not Found"
//	@Router			/fyc/presentcars/{id} [delete]
func DeletePresentCarAPI(c *gin.Context, db *bun.DB) {

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
	rowsAffected, err := DeletePresentCar(ctx, db, id)
	if err != nil {
		log.Err(err).Msg("Error deleting present car")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete present car",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present car found with the specified ID ------ affected rows 0 ",
			"code":    9,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Present car deleted successfully",
		"message": rowsAffected,
		"code":    8,
	})
}
