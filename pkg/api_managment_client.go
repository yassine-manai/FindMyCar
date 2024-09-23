package pkg

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetAllClientCreds godoc
//
//	@Summary		Get all client credentials
//	@Description	Get a list of all client credentials
//	@Tags			Client Credentials
//	@Produce		json
//	@Param			extra	query		string	false	"Include extra information if 'yes'"
//	@Success		200		{array}		ApiManage
//	@Router			/fyc/clientCreds [get]
func GetAllClientCredsApi(c *gin.Context) {
	ctx := context.Background()

	clCred, err := GetAllClientCred(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all client credentials")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all client credentials",
			"code":    10,
		})
		return
	}

	if len(clCred) == 0 {
		log.Warn().Msg("No client credentials found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No client credentials found",
			"code":    9,
		})
		return
	}

	log.Info().Msg("Returning client credentials")
	c.JSON(http.StatusOK, clCred)
}

// GetClientCredByID godoc
//
//	@Summary		Get client credential by ID
//	@Description	Get a specific client credential by ID
//	@Tags			Client Credentials
//	@Produce		json
//	@Param			client_id	path		string	true	"Client ID"
//	@Success		200	{object}	ApiManage
//	@Router			/fyc/clientCreds/{id} [get]
func GetClientCredByIDAPI(c *gin.Context) {
	idStr := c.Param("client_id")

	ctx := context.Background()
	clientCred, err := GetClientCredById(ctx, idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error retrieving client credential by ID")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Client credential not found",
			"code":    9,
		})
		return
	}

	log.Info().Msg("Returning client credential by ID")
	c.JSON(http.StatusOK, clientCred)
}

// AddClientCred godoc
//
//	@Summary		Add a new client credential
//	@Description	Add a new client credential to the database
//	@Tags			Client Credentials
//	@Accept			json
//	@Produce		json
//	@Param			clientCred	body		ApiManage	true	"Client credential data"
//	@Success		201	{object}	ApiManage
//	@Router			/fyc/clientCreds [post]
func AddClientCredAPI(c *gin.Context) {
	var clientCred ApiManage
	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Err(err).Msg("Invalid request payload for client credential creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := AddClientCred(ctx, &clientCred); err != nil {
		log.Err(err).Msg("Error creating client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Msg("Client credential created successfully")
	c.JSON(http.StatusCreated, clientCred)
}

// UpdateClientCred godoc
//
//	@Summary		Update a client credential
//	@Description	Update an existing client credential by ID
//	@Tags			Client Credentials
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string		true	"Client ID"
//	@Param			clientCred	body		ApiManage	true	"Updated client credential data"
//	@Success		200	{object}	ApiManage
//	@Router			/fyc/clientCreds/{id} [put]
func UpdateClientCredAPI(c *gin.Context) {
	idStr := c.Param("id")

	var clientCred ApiManage
	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Err(err).Msg("Invalid request payload for client credential update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": "Invalid client credential data",
			"code":    12,
		})
		return
	}

	if clientCred.ClientID != idStr {
		log.Warn().Str("id_param", idStr).Str("id_body", clientCred.ClientID).Msg("ID mismatch between path and body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The ID in the request body does not match the ID in the query parameter",
			"code":    13,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := UpdateClientCred(ctx, idStr, &clientCred)
	if err != nil {
		log.Err(err).Msg("Error updating client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("id", idStr).Msg("No client credential found to update")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No client credential found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Msg("Client credential updated successfully")
	c.JSON(http.StatusOK, clientCred)
}

// DeleteClientCred godoc
//
//	@Summary		Delete a client credential
//	@Description	Delete a client credential by ID
//	@Tags			Client Credentials
//	@Param			id	path		string	true	"Client ID"
//	@Success		200	{string}	string	"Client credential deleted successfully"
//	@Router			/fyc/clientCreds/{id} [delete]
func DeleteClientCredAPI(c *gin.Context) {
	idStr := c.Param("id")

	ctx := context.Background()
	rowsAffected, err := DeleteClientCred(ctx, idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("Error deleting client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("id", idStr).Msg("No client credential found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No client credential found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("id", idStr).Msg("Client credential deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": "Client credential deleted successfully",
		"code":    8,
	})
}
