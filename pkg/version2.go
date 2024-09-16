package pkg

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenRequester struct {
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	GrantType    string `form:"grant_type" binding:"required" default:"client_credentials"`
}

type TokenResponseTest struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type CarLocation struct {
	Facility    string `json:"facility"`
	SpotID      string `json:"spot_id"`
	PictureName string `json:"picture_name"`
}

type FindMyCarResponse struct {
	ResponseCode int           `json:"response_code"`
	Locations    []CarLocation `json:"locations"`
}

type PictureResponse struct {
	ImageData string `json:"image_data"`
}

const AccessTokenFake = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

// Package-level variables for car locations
var carLocationsEN = []CarLocation{
	{
		Facility:    "Downtown Parking Garage",
		SpotID:      "A-123",
		PictureName: "car_ABCD_1.jpg",
	},
	{
		Facility:    "Airport Long-Term Parking",
		SpotID:      "B-456",
		PictureName: "car_ABCD_2.jpg",
	},
}

var carLocationsAR = []CarLocation{
	{
		Facility:    "موقف سيارات وسط المدينة",
		SpotID:      "أ-123",
		PictureName: "car_ABCD_1.jpg",
	},
	{
		Facility:    "مواقف السيارات طويلة الأمد في المطار",
		SpotID:      "ب-456",
		PictureName: "car_ABCD_2.jpg",
	},
}

// Helper function to check if the token is valid
func checkToken(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authorization token is required",
			"code":    "13",
		})
		return false
	}

	// Extract the token from the "Bearer <token>" format
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != AccessTokenFake {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid access token",
			"code":    "13",
		})
		return false
	}

	return true
}

// @Summary Get an access token
// @Description Get an access token using client credentials
// @Description Valide PARAMS ---------- Client ID : 6 / Client Secret : 4711 / grant_type : client_credentials
// @Tags Car Location
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   client_id     formData   string  true  "Client ID"
// @Param   client_secret formData   string  true  "Client Secret"
// @Param   grant_type    formData   string  true  "Grant Type"
// @Success 200 {object} TokenResponseTest
// @Router /token [post]
func getToken(c *gin.Context) {
	var TokenRequester TokenRequester

	// Hardcoded values for client_id, client_secret, and grant_type
	hardcodedClientID := "6"
	hardcodedClientSecret := "4711"
	hardcodedGrantType := "client_credentials"

	if err := c.ShouldBind(&TokenRequester); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the provided values match the hardcoded ones
	if TokenRequester.ClientID != hardcodedClientID || TokenRequester.ClientSecret != hardcodedClientSecret || TokenRequester.GrantType != hardcodedGrantType {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid client credentials or grant type",
			"code":    "13",
		})
		return
	}

	// Mock token response if values are correct
	c.JSON(http.StatusOK, TokenResponseTest{
		AccessToken: AccessTokenFake,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	})
}

// @Summary Find a car by license plate
// @Description Find a car using the license plate number --- Token is eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
// @Description Valid LICENCE PLATE : ABCD ---
// @Tags Car Location
// @Accept  json
// @Produce  json
// @Param   license_plate query string true "License Plate"
// @Param   language      query string true "Language"
// @Param   fuzzy_logic   query bool   true "Fuzzy Logic"
// @Param   Authorization header string true "Bearer Token"
// @Success 200 {object} FindMyCarResponse
// @Router /findmycar [get]
func findMyCar(c *gin.Context) {

	if !checkToken(c) {
		return
	}

	licensePlate := c.Query("license_plate")
	language := c.Query("language")
	fuzzyLogic := c.Query("fuzzy_logic")

	simpleAR := CarLocation{
		Facility:    "مواقف السيارات طويلة الأمد في المطار",
		SpotID:      "ب-456",
		PictureName: "car_" + licensePlate + "_1.jpg",
	}

	simpleEN := CarLocation{
		Facility:    "Airport Long-Term Parking",
		SpotID:      "B-456",
		PictureName: "car_" + licensePlate + "_2.jpg",
	}

	if licensePlate == "ABCD" {

		if fuzzyLogic == "true" {
			if language == "EN" {
				c.JSON(http.StatusOK, FindMyCarResponse{
					ResponseCode: 200,
					Locations:    carLocationsEN,
				})
			} else if language == "AR" {
				c.JSON(http.StatusOK, FindMyCarResponse{
					ResponseCode: 200,
					Locations:    carLocationsAR,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Bad Request",
					"message": "Unsupported language code",
					"code":    "12",
				})
				return
			}
		} else if fuzzyLogic == "false" {
			if language == "EN" {
				c.JSON(http.StatusOK, FindMyCarResponse{
					ResponseCode: 200,
					Locations:    []CarLocation{simpleEN},
				})
			} else if language == "AR" {
				c.JSON(http.StatusOK, FindMyCarResponse{
					ResponseCode: 200,
					Locations:    []CarLocation{simpleAR},
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Bad Request",
					"message": "Unsupported language code",
					"code":    "12",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid fuzzy logic value",
				"code":    "12",
			})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "License plate not found",
			"code":    9,
		})
		return
	}
}

// @Summary Get a picture by picture name
// @Description Get an image using the picture name --- Token is eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
// @Tags Car Location
// @Produce  json
// @Param   picture_name query string true "Picture Name"
// @Param   language     query string true "Language"
// @Param   Authorization header string true "Bearer Token"
// @Success 200 {object} PictureResponse
// @Router /getpicture [get]
func getPicture(c *gin.Context) {
	pictureName := c.Query("picture_name")
	lang := c.Query("language")

	if !checkToken(c) {
		return
	}

	// Mock image data
	imageDataAR := "iVBORw0KGgoAAAANSUhdUgAAAAEAAAABCAYAAAAfFcSJAAAACklEQVR4nGMAAQAABQABDQ..."
	imageDataEN := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYA....AAAfFcSJAAAAAARDSDHFDFDFGsdsdS54DQ..."

	var carLocations []CarLocation
	switch lang {
	case "EN":
		carLocations = carLocationsEN
	case "AR":
		carLocations = carLocationsAR
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Unsupported language code",
			"code":    "12",
		})
		return
	}

	var imageData string
	for _, location := range carLocations {
		if location.PictureName == pictureName {
			if lang == "EN" {
				imageData = imageDataEN
			} else if lang == "AR" {
				imageData = imageDataAR
			}
			c.JSON(http.StatusOK, PictureResponse{
				ImageData: imageData,
			})
			return
		}
	}

	// If no picture found
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "No picture found",
		"code":    9,
	})
}
