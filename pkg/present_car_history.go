package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type PresentCarHistory struct {
	bun.BaseModel `json:"-" bun:"table:present_car_history"`
	ID            int                    `bun:"id,pk,autoincrement" json:"id"`
	LPN           string                 `bun:"lpn" json:"lpn" binding:"required"`
	CurrZoneID    *int                   `bun:"cur_zone_id" json:"cur_zone_id" binding:"required"`
	LastZoneID    *int                   `bun:"last_zone_id" json:"last_zone_id" binding:"required"`
	CamID         *int                   `bun:"cam_id" json:"cam_id" binding:"required"`
	Image         string                 `bun:"image" json:"image" binding:"required"`
	Confidence    *int                   `bun:"confidence" json:"confidence" binding:"required"`
	CameraBody    map[string]interface{} `bun:"cam_body,type:jsonb" json:"cam_body" binding:"required" swaggertype:"object"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required" swaggertype:"object"`
}

type ResponsePCH struct {
	bun.BaseModel `json:"-" bun:"table:present_car_history"`
	ID            int                    `bun:"id" json:"id"`
	LPN           string                 `bun:"lpn" json:"lpn"`
	CurrZoneID    *int                   `bun:"cur_zone_id" json:"cur_zone_id"`
	LastZoneID    *int                   `bun:"last_zone_id" json:"last_zone_id"`
	CamID         *int                   `bun:"cam_id" json:"cam_id"`
	Image         string                 `bun:"image" json:"image"`
	Confidence    *int                   `bun:"confidence" json:"confidence"`
	CameraBody    map[string]interface{} `bun:"cam_body" json:"cam_body" swaggertype:"object"`
}

// Get all history with extra data
func GetAllPresentHistoryExtra(ctx context.Context) ([]PresentCarHistory, error) {
	var History []PresentCarHistory
	err := Dbg.NewSelect().Model(&History).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all History with Extra Data: %w", err)
	}
	return History, nil
}

// Get all history
func GetAllPresentHistory(ctx context.Context) ([]ResponsePCH, error) {
	var resp []ResponsePCH
	err := Dbg.NewSelect().Model(&resp).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all History : %w", err)
	}
	return resp, nil
}

// Gt history by id
func GetPresentHistoryByID(ctx context.Context, id int) (*PresentCarHistory, error) {
	hist := new(PresentCarHistory)
	err := Dbg.NewSelect().Model(hist).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting History by id : %w", err)
	}
	return hist, nil
}

// Gt history by lpn
func GetPresentHistoryByLPN(ctx context.Context, lpn string) (*PresentCarHistory, error) {
	hist := new(PresentCarHistory)
	err := Dbg.NewSelect().Model(hist).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting History by lpn : %w", err)
	}
	return hist, nil
}

// create a new history
func CreatePresentHistory(ctx context.Context, newHist *PresentCarHistory) error {
	// Insert and get the auto-generated ID from the database
	_, err := Dbg.NewInsert().Model(newHist).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a history : %w", err)
	}
	log.Debug().Msgf("New History added with ID: %d", newHist.ID)

	return nil
}

// Update a history by ID
func UpdatePresentHistory(ctx context.Context, hist_id int, hist *PresentCarHistory) (int64, error) {
	res, err := Dbg.NewUpdate().Model(hist).Where("zone_id = ?", hist_id).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating history with id %d: %w", hist_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated history with ID: %d, rows affected: %d", hist_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a history by ID
func DeletePresentHistory(ctx context.Context, id int) (int64, error) {
	res, err := Dbg.NewDelete().Model(&PresentCarHistory{}).Where("ID = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting History with id %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted History with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
