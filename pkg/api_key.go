package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetAllClientCreds godoc
//
// @Summary		Get all client credentials
// @Description	Get a list of all client credentials
// @Tags			Client Credentials
// @Produce		json
// @Param			id		query		string	false	"Client ID"
// @Success		200		{array}		ApiKey
// @Router			/fyc/clientCreds [get]
func GetAllClientCredsApi(c *gin.Context) {
	ctx := context.Background()
	id := c.Query("id")
	log.Info().Msg("Fetching all client credentials")

	if id != "" {
		log.Info().Str("clientID", id).Msg("Fetching Client by ID")

		camera, err := GetClientCredById(ctx, id)
		if err != nil {
			log.Err(err).Str("Client_id", id).Msg("Error retrieving Client by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("camera_id", id).Msg("Client fetched successfully")
		c.JSON(http.StatusOK, camera)
		return
	}

	clCred, err := GetAllClientCred(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting all client credentials")
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

// AddClientCred godoc
//
//	@Summary		Add a new client credential
//	@Description	Add a new client credential to the database
//	@Tags			Client Credentials
//	@Accept			json
//	@Produce		json
//	@Param			clientCred	body		ApiKey	true	"Client credential data"
//	@Success		201	{object}	ApiKey
//	@Router			/fyc/clientCreds [post]
func AddClientCredAPI(c *gin.Context) {
	var clientCred ApiKey

	log.Info().Msg("Attempting to add new client credential")

	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for client credential creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := AddClientCred(ctx, &clientCred); err != nil {
		log.Error().Err(err).Msg("Error creating client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Str("client_id", clientCred.ClientID).Msg("Client credential created successfully")
	c.JSON(http.StatusCreated, clientCred)
}

// UpdateClientCred godoc
//
//	@Summary		Update a client credential
//	@Description	Update an existing client credential by ID
//	@Tags			Client Credentials
//	@Accept			json
//	@Produce		json
//	@Param			id			query		string		true	"Client ID"
//	@Param			clientCred	body		ApiKey	true	"Updated client credential data"
//	@Success		200	{object}	ApiKey
//	@Router			/fyc/clientCreds/{id} [put]
func UpdateClientCredAPI(c *gin.Context) {
	idStr := c.Query("id")

	log.Info().Str("client_id", idStr).Msg("Attempting to update client credential")

	var clientCred ApiKey
	if err := c.ShouldBindJSON(&clientCred); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for client credential update")
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
		log.Error().Err(err).Str("client_id", idStr).Msg("Error updating client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", idStr).Msg("No client credential found to update")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No client credential found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", idStr).Msg("Client credential updated successfully")
	c.JSON(http.StatusOK, clientCred)
}

// DeleteClientCred godoc
//
//	@Summary		Delete a client credential
//	@Description	Delete a client credential by ID
//	@Tags			Client Credentials
//	@Param			id	query		string	true	"Client ID"
//	@Success		200	{string}	string	"Client credential deleted successfully"
//	@Router			/fyc/clientCreds [delete]
func DeleteClientCredAPI(c *gin.Context) {
	idStr := c.Query("id")

	log.Info().Str("client_id", idStr).Msg("Attempting to delete client credential")

	ctx := context.Background()
	rowsAffected, err := DeleteClientCred(ctx, idStr)
	if err != nil {
		log.Error().Err(err).Str("client_id", idStr).Msg("Error deleting client credential")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete client credential",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("client_id", idStr).Msg("No client credential not found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No client credential found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", idStr).Msg("Client credential deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": "Client credential deleted successfully",
		"code":    8,
	})
}

// ChangeStateAPI godoc
//
//	@Summary		Change Client state or retrieve Client by ID
//	@Description	Change the state of a Client (e.g., enabled/disabled) or retrieve a client by ID
//	@Tags			Client Credentials
//	@Produce		json
//	@Param			state	query		bool	false	"Client State"
//	@Param			id		query		int 	false	"Client ID"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Router			/fyc/clientState [put]
func ChangeClientStateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()
	id := c.Query("id")

	stateStr := c.Query("state")
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		log.Err(err).Str("state", stateStr).Msg("Invalid state format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid state format",
			"message": "State must be a boolean value (true/false)",
			"code":    13,
		})
		return
	}

	rowsAffected, err := ChangeApiKeyState(ctx, id, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("client with id %s is already enabled", id) {
			log.Info().Str("client_id", id).Msg("client is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("client_id", id).Msg("Error changing client state")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "An unexpected error occurred",
			"message": fmt.Sprintf("client with id %s is already enabled", id),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("client_id", id).Msg("client not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Client not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("client_id", id).Bool("state", state).Msg("Client state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "Client state changed successfully",
		"rowsAffected": rowsAffected,
	})
}

// GetClientAPI godoc
//
//	@Summary		Get enabled clients or a specific client by clientID
//	@Description	Get a list of enabled clients or a specific client by ID with optional extra data
//	@Tags			Client Credentials
//	@Produce		json
//	@Param			id		query		string	false	"Client ID"
//	@Success		200		{object}	ApiKey		"List of enabled clients or a single client"
//	@Router			/fyc/clientEnabled [get]
func GetClientEnabledAPI(c *gin.Context) {
	log.Debug().Msg("Get Enabled API request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		Client, err := GetClientEnabledByID(ctx, idStr)
		if err != nil {
			log.Err(err).Str("Client_id", idStr).Msg("Error retrieving Client by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Client_id", idStr).Msg("Enabled Client fetched successfully")
		c.JSON(http.StatusOK, Client)
		return
	}

	// Fetch all enabled Clients
	Clients, err := GetClientListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled Clients")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled Clients",
			"code":    10,
		})
		return
	}

	if len(Clients) == 0 {
		log.Info().Msg("No enabled Clients found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled Clients found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Client_count", len(Clients)).Msg("Enabled Clients fetched successfully")
	c.JSON(http.StatusOK, Clients)
}

// GetClientAPI godoc
//
//	@Summary		Get deleted Clients or a specific Client by ID
//	@Description	Get a list of deleted Client or a specific Client by ID with optional extra data
//	@Tags			Client Credentials
//	@Produce		json
//	@Param			id		query		string	false	"Client ID"
//	@Success		200		{object}	ApiKey		"List of deleted Clients or a Client Client"
//	@Router			/fyc/clientsDeleted [get]
func GetClientDeletedAPI(c *gin.Context) {
	log.Debug().Msg("Get Client Deleted API request")
	ctx := context.Background()
	idStr := c.Query("id")

	if idStr != "" {
		Client, err := GetClientDeletedByID(ctx, idStr)
		if err != nil {
			log.Err(err).Str("Client_ID", idStr).Msg("Error retrieving Client by ID")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Client not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("Client_ID", idStr).Msg("Deleted Client fetched successfully")
		c.JSON(http.StatusOK, Client)
		return
	}

	// Fetch all deleted Clients
	Clients, err := GetClientListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted Clients")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted Clients",
			"code":    10,
		})
		return
	}

	if len(Clients) == 0 {
		log.Info().Msg("No deleted Clients found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted Clients found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Client_count", len(Clients)).Msg("Deleted Clients fetched successfully")
	c.JSON(http.StatusOK, Clients)
}
