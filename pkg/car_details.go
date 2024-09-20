package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type CarDetail struct {
	bun.BaseModel `json:"-" bun:"table:car_detail"`
	ID            int                    `bun:"id,pk,autoincrement" json:"id"`
	CamBody       map[string]interface{} `bun:"cam_body,type:jsonb" json:"cam_body" binding:"required"`
	Image1        string                 `bun:"image1,type:bytea" json:"image1" binding:"required"`
	Image2        string                 `bun:"image2,type:bytea" json:"image2" binding:"required"`
	Image3        string                 `bun:"image3,type:bytea" json:"image3" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseCarDetail struct {
	bun.BaseModel `json:"-" bun:"table:car_detail"`
	ID            int                    `bun:"id" json:"ID"`
	CamBody       map[string]interface{} `bun:"cam_body" json:"cam_body"`
	Image1        string                 `bun:"image1" json:"image1"`
	Image2        string                 `bun:"image2" json:"image2"`
	Image3        string                 `bun:"image3" json:"image3"`
}

// Get all car details with extra data
func GetAllCarDetailExtra(ctx context.Context) ([]CarDetail, error) {
	var cars []CarDetail
	err := Dbg.NewSelect().Model(&cars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all car details with extra data: %w", err)
	}
	return cars, nil
}

// Get all car details
func GetAllCarDetail(ctx context.Context) ([]ResponseCarDetail, error) {
	var cars []ResponseCarDetail
	err := Dbg.NewSelect().Model(&cars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all car details: %w", err)
	}
	return cars, nil
}

// Get car detail by ID
func GetCarDetailByID(ctx context.Context, id int) (*CarDetail, error) {
	car := new(CarDetail)
	err := Dbg.NewSelect().Model(car).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting car detail by ID: %w", err)
	}
	return car, nil
}

// Create a new car detail
func CreateCarDetail(ctx context.Context, newCar *CarDetail) error {
	// Insert and get the auto-generated ID from the database
	_, err := Dbg.NewInsert().Model(newCar).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating car detail: %w", err)
	}
	log.Debug().Msgf("New car detail added with ID: %d", newCar.ID)

	return nil
}

// Update a car detail by ID
func UpdateCarDetail(ctx context.Context, carID int, updates *CarDetail) (int64, error) {
	res, err := Dbg.NewUpdate().Model(updates).Where("id = ?", carID).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating car detail with ID %d: %w", carID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated car detail with ID: %d, rows affected: %d", carID, rowsAffected)

	return rowsAffected, nil
}

// Delete a car detail by ID
func DeleteCarDetail(ctx context.Context, id int) (int64, error) {
	res, err := Dbg.NewDelete().Model(&CarDetail{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting car detail with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted car detail with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
