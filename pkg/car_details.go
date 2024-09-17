package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type CarDetail struct {
	bun.BaseModel `json:"-" bun:"table:car_detail"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	CamBody       map[string]interface{} `bun:"cam_body,type:jsonb" json:"cam_body" binding:"required"`
	Image1        []byte                 `bun:"image1" json:"image1" binding:"required"`
	Image2        []byte                 `bun:"image2" json:"image2" binding:"required"`
	Image3        []byte                 `bun:"image3" json:"image3" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseCarDetail struct {
	bun.BaseModel `json:"-" bun:"table:car_detail"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	CamBody       map[string]interface{} `bun:"cam_body,type:jsonb" json:"cam_body"`
	Image1        []byte                 `bun:"image1" json:"image1"`
	Image2        []byte                 `bun:"image2" json:"image2"`
	Image3        []byte                 `bun:"image3" json:"image3"`
}

// Get all car details with extra data
func GetAllCarDetailExtra(ctx context.Context, db *bun.DB) ([]CarDetail, error) {
	var cars []CarDetail
	err := db.NewSelect().Model(&cars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all car details with extra data: %w", err)
	}
	return cars, nil
}

// Get all car details
func GetAllCarDetail(ctx context.Context, db *bun.DB) ([]ResponseCarDetail, error) {
	var cars []ResponseCarDetail
	err := db.NewSelect().Model(&cars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all car details: %w", err)
	}
	return cars, nil
}

// Get car detail by ID
func GetCarDetailByID(ctx context.Context, db *bun.DB, id int) (*CarDetail, error) {
	car := new(CarDetail)
	err := db.NewSelect().Model(car).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting car detail by ID: %w", err)
	}
	return car, nil
}

// Create a new car detail
func CreateCarDetail(ctx context.Context, db *bun.DB, newCar *CarDetail) error {
	// Insert and get the auto-generated ID from the database
	_, err := db.NewInsert().Model(newCar).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating car detail: %w", err)
	}
	log.Debug().Msgf("New car detail added with ID: %d", newCar.ID)

	return nil
}

// Update a car detail by ID
func UpdateCarDetail(ctx context.Context, db *bun.DB, carID int, updates *CarDetail) (int64, error) {
	res, err := db.NewUpdate().Model(updates).Where("id = ?", carID).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating car detail with ID %d: %w", carID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated car detail with ID: %d, rows affected: %d", carID, rowsAffected)

	return rowsAffected, nil
}

// Delete a car detail by ID
func DeleteCarDetail(ctx context.Context, db *bun.DB, id int) (int64, error) {
	res, err := db.NewDelete().Model(&CarDetail{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting car detail with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted car detail with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
