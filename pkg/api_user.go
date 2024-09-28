package pkg

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetAllUser godoc
//
// @Summary		Get all Users
// @Description	Get a list of all Users
// @Tags			Users
// @Produce		json
// @Param			username		query		string	false	"Username"
// @Success		200		{array}		User
// @Router			/fyc/users [get]
func GetAllUserApi(c *gin.Context) {
	ctx := context.Background()
	usernameStr := c.Query("username")
	log.Info().Msg("Fetching all Users Data")

	if usernameStr != "" {
		log.Info().Str("UserName", usernameStr).Msg("Fetching User by ID")

		user, err := GetUserByUsername(ctx, usernameStr)
		if err != nil {
			log.Err(err).Str("username", usernameStr).Msg("Error retrieving User")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "User not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("username", usernameStr).Msg("User fetched successfully")
		c.JSON(http.StatusOK, user)
		return
	}

	users, err := GetAllUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting all Users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error getting all Users",
			"code":    10,
		})
		return
	}

	if len(users) == 0 {
		log.Warn().Msg("No User found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No User found",
			"code":    9,
		})
		return
	}

	log.Info().Msg("Returning User data ")
	c.JSON(http.StatusOK, users)
}

// AddUserCred godoc
//
//	@Summary		Add a new User credential
//	@Description	Add a new User credential to the database
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			User	body		User	true	"User data"
//	@Success		201	{object}	User
//	@Router			/fyc/user [post]
func AddUserAPI(c *gin.Context) {
	var user User

	log.Info().Msg("Attempting to add new user")

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for user creation")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": err.Error(),
			"code":    12,
		})
		return
	}

	ctx := context.Background()
	if err := AddUser(ctx, &user); err != nil {
		log.Error().Err(err).Msg("Error creating user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	log.Info().Str("UserName", user.UserName).Msg("UserName created successfully")
	c.JSON(http.StatusCreated, user)
}

// UpdateClientCred godoc
//
//	@Summary		Update a client credential
//	@Description	Update an existing client credential by ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			username	query		string		true	"Client ID"
//	@Param			clientCred	body		User	true	"Updated client credential data"
//	@Success		200	{object}	User
//	@Router			/fyc/user/{id} [put]
func UpdateUserAPI(c *gin.Context) {
	usernameStr := c.Query("username")

	log.Info().Str("username", usernameStr).Msg("Attempting to update User")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Err(err).Msg("Invalid request payload for client credential update")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"message": "Invalid client credential data",
			"code":    12,
		})
		return
	}

	if user.UserName != usernameStr {
		log.Warn().Str("username param", usernameStr).Str("username", user.UserName).Msg("ID mismatch between path and body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID mismatch",
			"message": "The Username in the request body does not match the Username in the query parameter",
			"code":    13,
		})
		return
	}

	ctx := context.Background()
	rowsAffected, err := UpdateUser(ctx, usernameStr, &user)
	if err != nil {
		log.Error().Err(err).Str("client_id", usernameStr).Msg("Error updating USER")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update User",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("username", usernameStr).Msg("No user found to update")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No  user found with the specified ID",
			"code":    9,
		})
		return
	}

	log.Info().Str("username", usernameStr).Msg("User updated successfully")
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
//
//	@Summary		Delete a cUse
//	@Description	Delete a Use by username
//	@Tags			Users
//	@Param			username	query		string	true	"Username"
//	@Success		200	{string}	string	"User deleted successfully"
//	@Router			/fyc/user [delete]
func DeleteUserCredAPI(c *gin.Context) {
	userStr := c.Query("id")
	log.Info().Str("User", userStr).Msg("Attempting to delete user")
	ctx := context.Background()
	rowsAffected, err := DeleteUser(ctx, userStr)
	if err != nil {
		log.Error().Err(err).Str("user", userStr).Msg("Error deleting user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": err.Error(),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Warn().Str("username", userStr).Msg("No User found to delete")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No user found with the specified Username",
			"code":    9,
		})
		return
	}

	log.Info().Str("username", userStr).Msg("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": "User deleted successfully",
		"code":    8,
	})
}

// ChangeStateAPI godoc
//
//	@Summary		Change user state or retrieve user by username
//	@Description	Change the state of a user (e.g., enabled/disabled) or retrieve a user by username
//	@Tags			Users
//	@Produce		json
//	@Param			state		query		bool	false	"Client State"
//	@Param			username	query		string 	false	"Username"
//	@Success		200		{object}	int64		"Number of rows affected by the state change"
//	@Router			/fyc/userState [put]
func ChangeUserStateAPI(c *gin.Context) {
	log.Debug().Msg("ChangeStateAPI request")
	ctx := context.Background()
	username := c.Query("state")
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

	rowsAffected, err := ChangeApiKeyState(ctx, username, state)
	if err != nil {
		if err.Error() == fmt.Sprintf("User with username %s is already enabled", username) {
			log.Info().Str("username", username).Msg("User is already enabled")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Conflict",
				"message": err.Error(),
				"code":    12,
			})
			return
		}

		log.Err(err).Str("username", username).Msg("Error changing user state")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "An unexpected error occurred",
			"message": fmt.Sprintf("user with id %s is already enabled", username),
			"code":    10,
		})
		return
	}

	if rowsAffected == 0 {
		log.Info().Str("username", username).Msg("User not found or state unchanged")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "User not found or state unchanged",
			"code":    9,
		})
		return
	}

	log.Info().Str("username", username).Bool("state", state).Msg("User state changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message":      "User state changed successfully",
		"rowsAffected": rowsAffected,
	})
}

// GetClientAPI godoc
//
//	@Summary		Get enabled User or a specific User by username
//	@Description	Get a list of enabled users or a specific user by ID with optional extra data
//	@Tags			Users
//	@Produce		json
//	@Param			username		query		string	false	"UserName"
//	@Success		200		{object}	User		"List of enabled Users or a single User"
//	@Router			/fyc/userEnabled [get]
func GetUserEnabledAPI(c *gin.Context) {
	log.Debug().Msg("Get Enabled API request")
	ctx := context.Background()
	username := c.Query("username")

	if username != "" {
		User, err := GetUserEnabledByID(ctx, username)
		if err != nil {
			log.Err(err).Str("username", username).Msg("Error retrieving User by username")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "User not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("username", username).Msg("Enabled User fetched successfully")
		c.JSON(http.StatusOK, User)
		return
	}

	// Fetch all enabled Clients
	User, err := GetUserListEnabled(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving enabled User")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving enabled Users",
			"code":    10,
		})
		return
	}

	if len(User) == 0 {
		log.Info().Msg("No enabled Users found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No enabled Users found",
			"code":    9,
		})
		return
	}

	log.Info().Int("UserCount", len(User)).Msg("Enabled User fetched successfully")
	c.JSON(http.StatusOK, User)
}

// GetClientAPI godoc
//
//	@Summary		Get deleted User or a specific User by username
//	@Description	Get a list of deleted Users or a specific User by username with optional extra data
//	@Tags			Users
//	@Produce		json
//	@Param			username		query		string	false	"username"
//	@Success		200		{object}	User		"List of deleted Users or a User"
//	@Router			/fyc/userDeleted [get]
func GetUserDeletedAPI(c *gin.Context) {
	log.Debug().Msg("Get User Deleted API request")
	ctx := context.Background()
	username := c.Query("username")

	if username != "" {
		User, err := GetUserDeletedByID(ctx, username)
		if err != nil {
			log.Err(err).Str("Username", username).Msg("Error retrieving User by username")
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "User not found",
				"code":    9,
			})
			return
		}

		log.Info().Str("username", username).Msg("Deleted Users fetched successfully")
		c.JSON(http.StatusOK, User)
		return
	}

	// Fetch all deleted Clients
	Users, err := GetUserListDeleted(ctx)
	if err != nil {
		log.Err(err).Msg("Error retrieving deleted Users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "An unexpected error occurred",
			"message": "Error retrieving deleted Users",
			"code":    10,
		})
		return
	}

	if len(Users) == 0 {
		log.Info().Msg("No deleted Clients found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "No deleted Users found",
			"code":    9,
		})
		return
	}

	log.Info().Int("Users_count", len(Users)).Msg("Deleted Users fetched successfully")
	c.JSON(http.StatusOK, Users)
}
