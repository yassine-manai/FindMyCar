package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `json:"-" bun:"table:user"`
	ID            int    `bun:"id,autoincrement" json:"id"`
	UserName      string `bun:"username,pk" binding:"required" json:"username"`
	Password      string `bun:"password" binding:"required" json:"password"`
	FirstName     string `bun:"first_name" binding:"required" json:"first_name"`
	LastName      string `bun:"last_name" binding:"required" json:"last_name"`
	Role          string `bun:"role" binding:"required" json:"role"`
	IsEnabled     bool   `bun:"is_enabled,type:bool" json:"is_enabled" default:"false"`
	IsDeleted     bool   `bun:"is_deleted,type:bool" json:"is_deleted" default:"false"`
}

func AddUser(ctx context.Context, user *User) error {
	_, err := Dbg.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	return nil
}

func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := new(User)

	err := Dbg.NewSelect().Model(user).
		Where("username = ?", username).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", username)
		}
		return nil, fmt.Errorf("error retrieving User cred with username %s: %w", username, err)
	}
	return user, nil
}

func GetAllUsers(ctx context.Context) ([]User, error) {
	var user []User
	err := Dbg.NewSelect().Model(&user).Where("is_deleted = ?", false).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Users: %w", err)
	}
	return user, nil
}

func UpdateUser(ctx context.Context, username string, updatedUser *User) (int64, error) {
	log.Debug().Msgf("Updating user with Username: %s\n", username)
	result, err := Dbg.NewUpdate().
		Model(updatedUser).
		Where("is_deleted = ?", false).
		Where("username = ?", username).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating User cred with username %s: %w", username, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteUser(ctx context.Context, username string) (int64, error) {
	log.Debug().Msgf("Deleting User with username: %s", username)

	result, err := Dbg.NewUpdate().
		Model(&User{}).
		Set("is_deleted = ?", true).
		Where("username = ?", username).
		Exec(ctx)

	if err != nil {

		return 0, fmt.Errorf("error deleting client cred with ClientID %s: %w", username, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func GetUserListEnabled(ctx context.Context) ([]User, error) {
	var user []User
	err := Dbg.NewSelect().
		Model(&user).
		Where("is_enabled = ?", true).
		Scan(ctx, &user)

	if err != nil {
		return nil, fmt.Errorf("error getting enabled user list : %w", err)
	}

	return user, nil
}

func GetUserEnabledByID(ctx context.Context, username string) (*User, error) {
	var user User
	err := Dbg.NewSelect().
		Model(&user).
		Where("username = ?", username).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx, &user)

	if err != nil {
		return nil, fmt.Errorf("error getting client by Username %s: %w", username, err)
	}

	return &user, nil
}

func GetUserListDeleted(ctx context.Context) ([]User, error) {
	var user []User
	err := Dbg.NewSelect().
		Model(&user).
		Where("is_deleted = ?", true).
		Scan(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("error fetching deleted user list : %w", err)
	}
	return user, nil
}

func GetUserDeletedByID(ctx context.Context, username string) (*User, error) {
	var user User
	err := Dbg.NewSelect().
		Model(&user).
		Where("is_deleted = ?", true).
		Where("username = ?", username).
		Scan(ctx, &user)

	if err != nil {
		return nil, fmt.Errorf("error getting User by username: %w", err)
	}

	return &user, nil
}

func ChangeUserState(ctx context.Context, username string, newState bool) (int64, error) {
	user := new(User)

	err := Dbg.NewSelect().
		Model(user).
		Where("client_id = ?", username).
		Where("is_deleted = ?", false).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0, fmt.Errorf("error retrieving User state with username %s: %w", username, err)
	}

	if user.IsEnabled == newState {
		stateMessage := "already"
		if !newState {
			stateMessage = "disabled"
		} else {
			stateMessage = "enabled"
		}
		return 0, fmt.Errorf("client with id %s is already %s", username, stateMessage)
	}

	res, err := Dbg.NewUpdate().
		Model(&ApiKey{}).
		Set("is_enabled = ?", newState).
		Where("is_deleted = ?", false).
		Where("username = ?", username).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error changing User state with Username %s: %w", username, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Changed User State with UserName: %s, rows affected: %d", username, rowsAffected)

	return rowsAffected, nil
}
