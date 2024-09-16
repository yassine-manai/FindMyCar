package pkg

import (
	"fmc/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type TokenRequest struct {
	ClientID     string `json:"clientId" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
}

// TokenResponse represents the structure for token response
type TokenResponse struct {
	Token string `json:"token"`
}

// Claims defines the structure for JWT claims
type Claims struct {
	ClientID string `json:"clientId"`
	jwt.RegisteredClaims
}

// Hardcoded valid credentials for example purposes (should ideally be stored securely)
var validClientID = "6"
var validClientSecret = "4711"

// GenerateJWT generates a signed JWT token using the provided clientId.
func GenerateJWT(clientId string) (string, error) {
	var configvar config.ConfigFile
	if err := configvar.Load(); err != nil {
		log.Err(err).Msgf("Error loading config: %v", err)
	} else {
		fmt.Println("Success fetching config data ", configvar.App.JSecret)
	}

	var jwtKey = []byte(configvar.App.JSecret)

	// Set the token expiration time (1 hour in this case)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Define the claims, including the clientId and expiration time
	claims := &Claims{
		ClientID: clientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate a new token with claims and the HMAC SHA256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TokenHandler godoc
//
//	@Summary		Generate a JWT token
//	@Description	Generates a JWT token using client ID and client secret.
//	@Tags			Test_Version1
//	@Accept			json
//	@Produce		json
//	@Param			TokenRequest	body		TokenRequest	true	"Client credentials"
//	@Success		200	{object}	TokenResponse	"Token generated successfully"
//	@Failure		400	{object}	string			"Invalid request payload"
//	@Failure		401	{object}	string			"Invalid client credentials"
//	@Failure		405	{object}	string			"Method not allowed"
//	@Router			/fyc/v1/Auth/token [post]
func TokenHandler(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Method not allowed",
			"code":    11,
		})
		return
	}

	var req TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Invalid request payload",
			"code":    12,
		})
		return
	}

	if req.ClientID != validClientID || req.ClientSecret != validClientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Invalid client credentials",
			"code":    13,
		})
		return
	}

	token, err := GenerateJWT(req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Could not generate token",
			"code":    10,
		})

		return
	}

	// Return the generated token in a JSON response
	c.JSON(http.StatusOK, TokenResponse{
		Token: token,
	})
}
