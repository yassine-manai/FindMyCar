package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ApiKey struct {
	bun.BaseModel `json:"-" bun:"table:api_key"`
	ID            int    `bun:"id,autoincrement" json:"id"`
	ClientName    string `bun:"client_name" json:"client_name"`
	ClientID      string `bun:"client_id,pk" binding:"required" json:"client_id"`
	ClientSecret  string `bun:"client_secret,pk" binding:"required" json:"client_secret"`
	ApiKey        string `bun:"api_key" binding:"required" json:"api_key"`
	FuzzyLogic    bool   `bun:"fuzzy_logic,type:bool" json:"fuzzy" default:"false"`
	IsEnabled     bool   `bun:"is_enabled,type:bool" json:"is_enabled" default:"false"`
	IsDeleted     bool   `bun:"is_deleted,type:bool" json:"is_deleted" default:"false"`
}

func AddClientCred(ctx context.Context, apimgnt *ApiKey) error {
	_, err := Dbg.NewInsert().Model(apimgnt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding api_cred: %w", err)
	}

	return nil
}

func GetClientCredById(ctx context.Context, clientID string) (*ApiKey, error) {
	api := new(ApiKey)

	err := Dbg.NewSelect().Model(api).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving client cred with ClientID %s: %w", clientID, err)
	}
	return api, nil
}

func GetClientCredBySecret(ctx context.Context, clientSecret string) (*ApiKey, error) {
	api := new(ApiKey)

	err := Dbg.NewSelect().Model(api).Where("client_secret = ?", clientSecret).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with Client Secret %s not found", clientSecret)
		}
		return nil, fmt.Errorf("error retrieving client cred with Client Secret %s: %w", clientSecret, err)
	}
	return api, nil
}

func GetAllClientCred(ctx context.Context) ([]ApiKey, error) {
	var apm []ApiKey
	err := Dbg.NewSelect().Model(&apm).Where("is_deleted = ?", false).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Client Credentials: %w", err)
	}
	return apm, nil
}

func UpdateClientCred(ctx context.Context, clientID string, updatedClientCred *ApiKey) (int64, error) {
	log.Debug().Msgf("Updating clientcred with ClientID: %s\n", clientID)
	result, err := Dbg.NewUpdate().
		Model(updatedClientCred).
		Where("is_deleted = ?", false).
		Where("client_id = ?", clientID).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating client cred with ClientID %s: %w", clientID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteClientCred(ctx context.Context, clientID string) (int64, error) {
	log.Debug().Msgf("Deleting Client with ClientID: %s", clientID)

	result, err := Dbg.NewUpdate().
		Model(&ApiKey{}).
		Set("is_deleted = ?", true).
		Where("client_id = ?", clientID).
		Exec(ctx)

	if err != nil {

		return 0, fmt.Errorf("error deleting client cred with ClientID %s: %w", clientID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func GetClientListEnabled(ctx context.Context) ([]ApiKey, error) {
	var api []ApiKey
	err := Dbg.NewSelect().
		Model(&api).
		Where("is_enabled = ?", true).
		Scan(ctx, &api)

	if err != nil {
		return nil, fmt.Errorf("error getting enabled API list : %w", err)
	}

	return api, nil
}

func GetClientEnabledByID(ctx context.Context, id string) (*ApiKey, error) {
	var api_key ApiKey
	err := Dbg.NewSelect().
		Model(&api_key).
		Where("client_id = ?", id).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx, &api_key)

	if err != nil {
		return nil, fmt.Errorf("error getting client API by id %s: %w", id, err)
	}

	return &api_key, nil
}

func GetClientListDeleted(ctx context.Context) ([]ApiKey, error) {
	var api []ApiKey
	err := Dbg.NewSelect().
		Model(&api).
		Where("is_deleted = ?", true).
		Scan(ctx, &api)
	if err != nil {
		return nil, fmt.Errorf("error fetching deleted api list : %w", err)
	}
	return api, nil
}

func GetClientDeletedByID(ctx context.Context, id string) (*ApiKey, error) {
	var apk ApiKey
	err := Dbg.NewSelect().
		Model(&apk).
		Where("is_deleted = ?", true).
		Where("client_id = ?", id).
		Scan(ctx, &apk)

	if err != nil {
		return nil, fmt.Errorf("error getting Client API by id: %w", err)
	}

	return &apk, nil
}

func ChangeApiKeyState(ctx context.Context, clientID string, newState bool) (int64, error) {
	existingApiKey := new(ApiKey)

	err := Dbg.NewSelect().
		Model(existingApiKey).
		Where("client_id = ?", clientID).
		Where("is_deleted = ?", false).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0, fmt.Errorf("error retrieving Client state with id %s: %w", clientID, err)
	}

	if existingApiKey.IsEnabled == newState {
		stateMessage := "already"
		if !newState {
			stateMessage = "disabled"
		} else {
			stateMessage = "enabled"
		}
		return 0, fmt.Errorf("client with id %s is already %s", clientID, stateMessage)
	}

	// Step 2: Change the state since it's different
	res, err := Dbg.NewUpdate().
		Model(&ApiKey{}).
		Set("is_enabled = ?", newState).
		Where("is_deleted = ?", false).
		Where("client_id = ?", clientID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error changing client_id state with id %s: %w", clientID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Changed client_id State with ID: %s, rows affected: %d", clientID, rowsAffected)

	return rowsAffected, nil
}
