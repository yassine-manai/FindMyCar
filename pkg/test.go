package pkg

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FYCHandler godoc
//
//	@Summary		Get car information
//	@Description	Get car information based on LPN, language, and fuzzy logic
//	@Tags			Test_Version1
//	@Param			LPN			query		string	true	"License Plate Number"
//	@Param			L			query		string	true	"Language (EN, AR, etc.)"
//	@Param			FuzzyLogic	query		string	true	"Fuzzy logic setting (On/Off)"
//	@Success		200			{object}	string	"Success"
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		401			{object}	string	"Unauthorized"
//	@Failure		404			{object}	string	"Not Found"
//	@Router			/fyc/v1 [get]
func (api PresentCarAPI) FYCHandler(c *gin.Context) {

	/* // Hardcoded access token for validation
	const validAccessToken = "123456789"

	// Check for the access token in the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "Bearer "+validAccessToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid access token",
			"code":    401,
		})
		return
	} */

	// Retrieve query parameters
	LPN := c.Query("LPN")
	Language := c.Query("L")
	FuzzyLogic := c.Query("FuzzyLogic")

	// Check if required query parameters are provided
	if LPN == "" || Language == "" || FuzzyLogic == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Missing required query parameters (LPN, Language, or FuzzyLogic)",
			"code":    400,
		})
		return
	}

	const Fuz = "ON"

	// Retrieve car information by LPN
	ctx := context.Background()
	car, err := api.PresentCarService.GetPresentCarByLPN(ctx, LPN)
	if err != nil {
		log.Err(err).Str("lpn", LPN).Msg("Error retrieving present car by LPN")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Present car not found",
			"code":    9,
		})
		return
	}

	// Validate fuzzy logic setting and LPN
	if FuzzyLogic == Fuz || car.LPN == LPN {
		c.JSON(http.StatusOK, gin.H{
			"LPN":        LPN,
			"Language":   Language,
			"FuzzyLogic": FuzzyLogic,
			"message":    "Success",
			"code":       200,
		})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No present cars found matching LPN",
			"code":    9,
		})
		return
	}
}
