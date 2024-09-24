package pkg

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ErrorMessage struct {
	bun.BaseModel `json:"-" bun:"table:errors"`
	Code          int               `bun:"code" json:"code"`
	Messages      map[string]string `bun:"messages, type:json" json:"messages"`
}

// CreateErrorMessage inserts a new error message into the database
func CreateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Dbg.NewInsert().Model(errMsg).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert error message into database")
	} else {
		log.Info().Int("code", errMsg.Code).Msg("Successfully inserted error message")
	}
	return err
}

// getErrorMessage fetches an error message by code and language from the database
func GetErrorMessageByFilter(ctx context.Context, code int, language string) (ErrorMessage, error) {
	var errMsg ErrorMessage
	err := Dbg.NewSelect().
		Model(&errMsg).
		Where("code", code).
		Where("messages ->> ?", language).
		Scan(ctx)

	if err != nil {
		log.Error().Err(err).Int("code", code).Str("language", language).Msg("Failed to fetch error message from database")
	} else {
		log.Info().Int("code", code).Str("language", language).Msg("Successfully fetched error message")
	}
	return errMsg, err
}

// getErrorMessage fetches an error message by code and language from the database
func GetErrorMessage(ctx context.Context) ([]ErrorMessage, error) {
	var errMsg []ErrorMessage
	err := Dbg.NewSelect().
		Model(&errMsg).
		Scan(ctx)

	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch error message from database")
	} else {
		log.Info().Msg("Successfully fetched error message")
	}
	return errMsg, err
}

// UpdateErrorMessage updates an error message's specific language in the database
func UpdateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Dbg.NewUpdate().
		Model(errMsg).
		Set("messages = ?", errMsg.Messages).
		Where("code = ?", errMsg.Code).
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Int("code", errMsg.Code).Msg("Failed to update error message in database")
	} else {
		log.Info().Int("code", errMsg.Code).Msg("Successfully updated error message")
	}
	return err
}

// DeleteErrorMessage deletes a specific language from the messages field of an error message by code
func DeleteErrorMessage(ctx context.Context, code int, language string) (int64, error) {
	// Use JSONB "- key" syntax to remove the language key from the messages field
	res, err := Dbg.NewUpdate().
		Model((*ErrorMessage)(nil)).
		Set("messages = messages - ?", language).
		Where("code = ?", code).
		Exec(ctx)

	if err != nil {
		log.Error().Err(err).Int("code", code).Str("language", language).Msg("Failed to delete language from error message")
		return 0, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Warn().Int("code", code).Str("language", language).Msg("No error message found or no language removed")
	} else {
		log.Info().Int("code", code).Str("language", language).Int64("rowsAffected", rowsAffected).Msg("Successfully deleted language from error message")
	}

	return rowsAffected, nil
}
