package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Sign struct {
	bun.BaseModel `json:"-" bun:"table:sign"`
	ID            int    `bun:"id,autoincrement" json:"id"`
	SignID        int    `bun:"sign_id,pk" binding:"required" json:"sign_id"`
	SignName      string `bun:"sign_name" binding:"required" json:"sign_name"`
	SignType      string `bun:"sign_type" binding:"required" json:"sign_type"`
	SignIP        string `bun:"sign_ip" binding:"required" json:"sign_ip"`
	SignPort      int    `bun:"sign_port" binding:"required" json:"sign_port"`
	ZoneID        int    `bun:"zone_id" binding:"required" json:"zone_id"`
	IsEnabled     bool   `bun:"is_enabled,type:bool" json:"is_enabled" default:"false"`
	IsDeleted     bool   `bun:"is_deleted,type:bool" json:"is_deleted" default:"false"`
}

func CreateSign(ctx context.Context, sign *Sign) error {
	_, err := Dbg.NewInsert().Model(sign).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding sign: %w", err)
	}
	return nil
}

func GetSignById(ctx context.Context, signID int) (*Sign, error) {
	sign := new(Sign)

	err := Dbg.NewSelect().Model(sign).
		Where("sign_id = ?", signID).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("sign with SignID %d not found", signID)
		}
		return nil, fmt.Errorf("error retrieving sign with SignID %d: %w", signID, err)
	}
	return sign, nil
}

func GetAllSigns(ctx context.Context) ([]Sign, error) {
	var signs []Sign
	err := Dbg.NewSelect().Model(&signs).Where("is_deleted = ?", false).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all signs: %w", err)
	}
	return signs, nil
}

func UpdateSign(ctx context.Context, signID int, updatedSign *Sign) (int64, error) {
	log.Debug().Msgf("Updating sign with SignID: %d\n", signID)
	result, err := Dbg.NewUpdate().
		Model(updatedSign).
		Where("is_deleted = ?", false).
		Where("sign_id = ?", signID).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating sign with SignID %d: %w", signID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteSign(ctx context.Context, signID int) (int64, error) {
	log.Debug().Msgf("Deleting Sign with SignID: %d", signID)

	result, err := Dbg.NewUpdate().
		Model(&Sign{}).
		Set("is_deleted = ?", true).
		Where("sign_id = ?", signID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error deleting sign with SignID %d: %w", signID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func GetSignListEnabled(ctx context.Context) ([]Sign, error) {
	var signs []Sign
	err := Dbg.NewSelect().
		Model(&signs).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting enabled sign list: %w", err)
	}

	return signs, nil
}

func GetEnabledSignByID(ctx context.Context, signID int) (*Sign, error) {
	var sign Sign
	err := Dbg.NewSelect().
		Model(&sign).
		Where("sign_id = ?", signID).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting enabled sign by ID %d: %w", signID, err)
	}

	return &sign, nil
}

func GetSignListDeleted(ctx context.Context) ([]Sign, error) {
	var signs []Sign
	err := Dbg.NewSelect().
		Model(&signs).
		Where("is_deleted = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching deleted sign list: %w", err)
	}
	return signs, nil
}

func GetDeletedSignByID(ctx context.Context, signID int) (*Sign, error) {
	var sign Sign
	err := Dbg.NewSelect().
		Model(&sign).
		Where("is_deleted = ?", true).
		Where("sign_id = ?", signID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting deleted sign by ID: %w", err)
	}

	return &sign, nil
}

func ChangeSignState(ctx context.Context, signID int, newState bool) (int64, error) {
	existingSign := new(Sign)

	err := Dbg.NewSelect().
		Model(existingSign).
		Where("sign_id = ?", signID).
		Where("is_deleted = ?", false).
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0, fmt.Errorf("error retrieving Sign state with ID %d: %w", signID, err)
	}

	if existingSign.IsEnabled == newState {
		stateMessage := "already"
		if !newState {
			stateMessage = "disabled"
		} else {
			stateMessage = "enabled"
		}
		return 0, fmt.Errorf("sign with ID %d is already %s", signID, stateMessage)
	}

	// Change the state
	res, err := Dbg.NewUpdate().
		Model(&Sign{}).
		Set("is_enabled = ?", newState).
		Where("is_deleted = ?", false).
		Where("sign_id = ?", signID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error changing sign state with ID %d: %w", signID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Changed Sign State with ID: %d, rows affected: %d", signID, rowsAffected)

	return rowsAffected, nil
}
