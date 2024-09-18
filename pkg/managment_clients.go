package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ApiManage struct {
	bun.BaseModel `json:"-" bun:"table:api_management"`

	ClientID     string `bun:"client_id,pk" json:"client_id"`
	ClientSecret string `bun:"client_secret" binding:"required" json:"client_secret"`
	Fuzzy        string `bun:"fuzzy" binding:"required" json:"fuzzy"`
}

func AddClientCred(ctx context.Context, apimgnt *ApiManage) error {
	_, err := Dbg.NewInsert().Model(apimgnt).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding api_cred: %w", err)
	}

	return nil
}

func GetClientCredById(ctx context.Context, clientID string) (*ApiManage, error) {
	api := new(ApiManage)

	err := Dbg.NewSelect().Model(api).Where("client_id = ?", clientID).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with ClientID %s not found", clientID)
		}
		return nil, fmt.Errorf("error retrieving client cred with ClientID %s: %w", clientID, err)
	}
	return api, nil
}

func GetClientCredBySecret(ctx context.Context, clientSecret string) (*ApiManage, error) {
	api := new(ApiManage)

	err := Dbg.NewSelect().Model(api).Where("client_secret = ?", clientSecret).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("client with Client Secret %s not found", clientSecret)
		}
		return nil, fmt.Errorf("error retrieving client cred with Client Secret %s: %w", clientSecret, err)
	}
	return api, nil
}

func GetAllClientCred(ctx context.Context) ([]ApiManage, error) {
	var apm []ApiManage
	err := Dbg.NewSelect().Model(&apm).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Client Credentials: %w", err)
	}
	return apm, nil
}

func UpdateClientCred(ctx context.Context, clientID string, updatedClientCred *ApiManage) (int64, error) {
	log.Debug().Msgf("Updating clientcred with ClientID: %s\n", clientID)
	result, err := Dbg.NewUpdate().
		Model(updatedClientCred).
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

	result, err := Dbg.NewDelete().Model((*ApiManage)(nil)).Where("client_id = ?", clientID).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting client cred with ClientID %s: %w", clientID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}
