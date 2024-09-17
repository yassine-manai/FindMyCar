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

type PresentCarAPI struct {
	PresentCarService *PresentCarOp
}

func NewPresentCarAPI(db *bun.DB) *PresentCarAPI {
	return &PresentCarAPI{
		PresentCarService: NewPresent(db),
	}
}

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
func (api *PresentCarAPI) GetPresentCars(c *gin.Context) {
	ctx := context.Background()
	extra_req := c.DefaultQuery("extra", "false")

	if strings.ToLower(extra_req) == "true" || strings.ToLower(extra_req) == "1" || strings.ToLower(extra_req) == "yes" {
		cars, err := api.PresentCarService.GetAllPresentExtra(ctx)
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

	Pcars, err := api.PresentCarService.GetAllPresentCars(ctx)
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
func (api *PresentCarAPI) GetPresentCarByLPN(c *gin.Context) {

	lpn := c.Param("lpn")
	extra_req := c.Query("extra")

	ctx := context.Background()
	car, err := api.PresentCarService.GetPresentCarByLPN(ctx, lpn)
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
func (api *PresentCarAPI) CreatePresentCar(c *gin.Context) {
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

	ZoneService := NewZone(api.PresentCarService.DB)
	_, err := ZoneService.GetZoneByID(ctx, car.CurrZoneID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", car.CurrZoneID),
			"code":    14,
		})
		return
	}

	_, errr := ZoneService.GetZoneByID(ctx, car.LastZoneID)
	if errr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", car.LastZoneID),
			"code":    14,
		})
		return
	}

	CamService := NewCamera(api.PresentCarService.DB)
	_, errCam := CamService.GetCameraByID(ctx, car.CameraID)
	if errCam != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera not found",
			"message": fmt.Sprintf("Camera with ID %d does not exist", car.CameraID),
			"code":    14,
		})
		return
	}

	// Insert the new car into the database
	if err := api.PresentCarService.CreatePresentCar(ctx, &car); err != nil {
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
func (api *PresentCarAPI) UpdatePresentCarById(c *gin.Context) {
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

	// Check if the ID in the request body matches the URL ID
	if updates.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the param",
			"code":    13,
		})
		return
	}

	ZoneService := NewZone(api.PresentCarService.DB)
	ctx := context.Background()
	_, errup := ZoneService.GetZoneByID(ctx, updates.CurrZoneID)
	if errup != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Current Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.CurrZoneID),
			"code":    14,
		})
		return
	}

	_, errr := ZoneService.GetZoneByID(ctx, updates.LastZoneID)
	if errr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Last Zone not found",
			"message": fmt.Sprintf("Zone with ID %d does not exist", updates.LastZoneID),
			"code":    14,
		})
		return
	}

	CamService := NewCamera(api.PresentCarService.DB)
	_, errCam := CamService.GetCameraByID(ctx, updates.CameraID)
	if errCam != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Camera not found",
			"message": fmt.Sprintf("Camera with ID %d does not exist", updates.CameraID),
			"code":    14,
		})
		return
	}

	// Call the service to update the present car
	rowsAffected, err := api.PresentCarService.UpdatePresentCar(ctx, id, &updates)
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
func (api *PresentCarAPI) UpdatePresentCarBylpn(c *gin.Context) {

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

	ctx := context.Background()
	rowsAffected, err := api.PresentCarService.UpdatePresentCarByLpn(ctx, lpn, &updates)
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
func (api *PresentCarAPI) DeletePresentCar(c *gin.Context) {

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
	rowsAffected, err := api.PresentCarService.DeletePresentCar(ctx, id)
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
