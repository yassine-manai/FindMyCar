package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ErrorMessage struct {
	bun.BaseModel `json:"-" bun:"table:error_message"`
	Code          string `bun:"code,pk" json:"code" binding:"required"`
	Language      string `bun:"language" ,default:"en" json:"language"`
	Message       string `bun:"message" json:"message" binding:"required"`
}

func GetAllErrors(ctx context.Context) ([]ErrorMessage, error) {
	var errors []ErrorMessage
	err := Dbg.NewSelect().Model(&errors).Scan(ctx)
	if err != nil {
		log.Err(err).Msg("Error getting all error messages")
		return nil, fmt.Errorf("error getting all error messages: %w", err)
	}
	log.Info().Int("count", len(errors)).Msg("Successfully retrieved all error messages")
	return errors, nil
}

func GetErrorMessageByCode(ctx context.Context, code string) (*ErrorMessage, error) {
	var errMsg ErrorMessage
	err := Dbg.NewSelect().
		Model(&errMsg).
		Where("code = ?", code).
		Scan(ctx)
	if err != nil {
		log.Err(err).Str("code", code).Msg("Error getting error message by code")
		return nil, fmt.Errorf("error getting error message with code %s: %w", code, err)
	}
	log.Info().Str("code", code).Msg("Successfully retrieved error message")
	return &errMsg, nil
}

func CreateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Dbg.NewInsert().Model(errMsg).Exec(ctx)
	if err != nil {
		log.Err(err).Str("code", errMsg.Code).Msg("Failed to insert error message")
		return fmt.Errorf("failed to insert error message: %w", err)
	}
	log.Info().Str("code", errMsg.Code).Msg("Successfully inserted error message")
	return nil
}

func UpdateErrorMessage(ctx context.Context, errMsg *ErrorMessage) error {
	_, err := Dbg.NewUpdate().
		Model(errMsg).
		Where("code = ?", errMsg.Code).
		Exec(ctx)
	if err != nil {
		log.Err(err).Str("code", errMsg.Code).Msg("Failed to update error message")
		return fmt.Errorf("failed to update error message with code %s: %w", errMsg.Code, err)
	}
	log.Info().Str("code", errMsg.Code).Msg("Successfully updated error message")
	return nil
}

func DeleteErrorMessage(ctx context.Context, code string) error {
	_, err := Dbg.NewDelete().
		Model((*ErrorMessage)(nil)).
		Where("code = ?", code).
		Exec(ctx)
	if err != nil {
		log.Err(err).Str("code", code).Msg("Failed to delete error message")
		return fmt.Errorf("failed to delete error message with code %s: %w", code, err)
	}
	log.Info().Str("code", code).Msg("Successfully deleted error message")
	return nil
}
