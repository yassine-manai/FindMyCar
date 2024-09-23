package pkg

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
	Facility    int    `json:"facility"`
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
		Facility:    1,
		SpotID:      "A-123",
		PictureName: "car_ABCD_1.jpg",
	},
	{
		Facility:    1,
		SpotID:      "B-456",
		PictureName: "car_ABCD_2.jpg",
	},
}

var carLocationsAR = []CarLocation{
	{
		Facility:    2,
		SpotID:      "أ-123",
		PictureName: "car_ABCD_1.jpg",
	},
	{
		Facility:    5,
		SpotID:      "ب-456",
		PictureName: "car_ABCD_2.jpg",
	},
}

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
func GetToken(c *gin.Context) {
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
func FindMyCar(c *gin.Context) {

	if !checkToken(c) {
		return
	}

	licensePlate := c.Query("license_plate")
	language := c.Query("language")
	fuzzyLogic := c.DefaultQuery("fuzzy_logic", "false")

	ctx := context.Background()

	car, err := GetPresentCarByLPN(ctx, licensePlate)
	if err != nil {
		log.Err(err).Str("license_plate", licensePlate).Msg("Error retrieving car by LPN")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Car Not Found",
			"message": "No car found with the provided license plate",
			"code":    404,
		})
		return
	}

	log.Info().Msg("Car found with license plate")
	log.Debug().Int("zone", *car.CurrZoneID).Msg("Current Zone ID")

	zoneImages, err := GetZoneImageByZoneID(ctx, *car.CurrZoneID)
	if err != nil {
		log.Err(err).Msg("Error retrieving zone image")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Zone Image Not Found",
			"message": "No image found for the current zone",
			"code":    404,
		})
		return
	}

	Response := CarLocation{
		Facility:    *car.CurrZoneID,
		SpotID:      "B-456",
		PictureName: "",
	}

	for _, zoneImage := range zoneImages {
		if language == zoneImage.Lang || fuzzyLogic == "false" {
			Response.PictureName = zoneImage.ImageLg
			c.JSON(http.StatusOK, gin.H{
				"ResponseCode": 200,
				"Locations":    Response,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "No Match",
		"message": "No matching car location found for the given language",
		"code":    9,
	})
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
func GetPicture(c *gin.Context) {
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

	// If not found
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "No picture found",
		"code":    9,
	})
}

type Settings struct {
	Logo               string `json:"logo"`
	TimeOutScreenKiosk int    `json:"timeout_screenKiosk"`
	DesfaultLanguage   string `json:"def_lang"`
	FuzzyLogic         bool   `json:"fuzzyLogic"`
}

// @Summary Get Settings
// @Description Get Settings
// @Description Get settings data --- Token is eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
// @Tags Car Location
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   Authorization header string true "Bearer Token"
// @Success 200 {object} Settings
// @Router /getSettings [get]
func Getsettings(c *gin.Context) {
	TestLogo := "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTkiIGhlaWdodD0iMzYiIHZpZXdCb3g9IjAgMCAxOSAzNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGQ9Ik0xMS4wMjU3IDMyLjYzODdDNi42MTE2MiAzMi42Mzg3IDMuMDMxOTIgMjkuMDU0IDMuMDMxOTIgMjQuNjMzOEwzLjA0NDY2IDE0LjczNDZDMy4wNDQ2NiAxNC43MDkgMy4wNDQ2NiAxNC42ODM1IDMuMDQ0NjYgMTQuNjUxNkMzLjA0NDY2IDEzLjc1ODcgMi4zNjMxMSAxMy4wMjUxIDEuNTIyMzMgMTMuMDI1MUMwLjY4MTU0NCAxMy4wMjUxIDAgMTMuNzU4NyAwIDE0LjY1MTZDMCAxNC42NzcxIDAgMTQuNzA5IDAgMTQuNzM0NlYyNC42MzM4QzAgMjkuMjU4MSAyLjg0MDgzIDMzLjIxMjcgNi44NjY0IDM0Ljg1ODRDMTYuOTExMiAzOC4yMzg5IDE4LjI0ODggMjguMDY1NCAxOC4yNDg4IDI4LjA2NTRDMTYuOTYyMiAzMC43Njk4IDE0LjIxNjkgMzIuNjM4NyAxMS4wMjU3IDMyLjYzODdaIiBmaWxsPSIjMDBBQjU0Ii8+CiAgICA8cGF0aCBkPSJNMTYuMTM0MiAxMy4xNDY0QzE1LjQ3MTggMTMuNDI3MSAxNS4yMDQzIDE0LjA1ODUgMTUuMjA0MyAxNC43MzQ2QzE1LjIwNDMgMTYuMDAzOSAxNS4yMjM0IDI0LjgxMjUgMTUuMTM0MiAyNS4zNTQ2QzE0Ljc4MzkgMjcuNDA4NSAxMi45MDQ5IDI4LjkyMDIgMTAuODM0NyAyOC44MjQ1QzguNzc3MzcgMjguNzI4OCA3LjA0NDg1IDI3LjA3NjggNi44NjAxMyAyNS4wMjNWNC43MjY5M0wzLjgwMjczIDguMjc5NjlWMjMuOTUxNEMzLjgwMjczIDI0LjkwODIgMy44MzQ1OCAyNS44NTIyIDQuMTIxMjEgMjYuNzgzNEM0LjcwNzIxIDI4LjY5MDUgNi4wOTU3OCAzMC4yOTc5IDcuODg1NjMgMzEuMTY1NEM5LjczOTE4IDMyLjA1ODMgMTEuOTQzIDMyLjExNTcgMTMuODM0OCAzMS4zMTIxQzE1LjY2OTIgMzAuNTMzOSAxNy4xMjc5IDI4Ljk5NjcgMTcuODA5NCAyNy4xMjc4QzE4LjE5MTYgMjYuMDc1NCAxOC4yNDI2IDI1LjAwMzggMTguMjQyNiAyMy45MDA0VjE0LjU0OTZDMTguMjQ4OSAxMy40NzgxIDE3LjEyMTUgMTIuNzI1NCAxNi4xMzQyIDEzLjE0NjRaIiBmaWxsPSIjMkY1RkFDIi8+CiAgICA8cGF0aCBkPSJNMTQuNDU5MiA4LjQ1ODI1TDExLjQwMTggNC45MDU0OUwxMS40MDgxIDI0LjUxOUMxMS40MDgxIDI0LjcyOTUgMTEuMjM2MSAyNC44OTU0IDExLjAzMjMgMjQuODk1NEMxMC44MjIxIDI0Ljg5NTQgMTAuNjU2NSAyNC43MjMxIDEwLjY1NjUgMjQuNTE5TDEwLjY2MjkgNC4wMDYxM0w3LjU5OTEyIDAuNDUzMzY5QzcuNTk5MTIgMC40NTMzNjkgNy42MTE4NiAyNC40NzQ0IDcuNjExODYgMjQuNTE5QzcuNjExODYgMjYuNDEzNCA5LjE0NjkzIDI3Ljk0NDIgMTEuMDMyMyAyNy45NDQyQzEyLjcxMzkgMjcuOTQ0MiAxNC4xMTUyIDI2LjcyNiAxNC4zOTU1IDI1LjExODZDMTQuNDI3MyAyNC45OTc0IDE0LjQ0NjQgMjQuODYzNSAxNC40NDY0IDI0LjcyOTVIMTQuNDUyOEMxNC40NDY0IDI0LjM2NTkgMTQuNDU5MiA4LjQ1ODI1IDE0LjQ1OTIgOC40NTgyNVoiIGZpbGw9IiMyRjVGQUMiLz4KPC9zdmc+Cg=="

	if !checkToken(c) {
		return
	}

	c.JSON(http.StatusOK, Settings{
		Logo:               TestLogo,
		TimeOutScreenKiosk: 10,
		DesfaultLanguage:   "EN",
		FuzzyLogic:         false,
	})
}
