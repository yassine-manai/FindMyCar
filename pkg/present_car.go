package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type PresentCar struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
	ID              *int                   `bun:"id,pk,autoincrement" json:"id"`
	TransactionDate string                 `bun:"transaction_date" json:"transaction_date" binding:"required"`
	CameraID        *int                   `bun:"camera_id" json:"came ra_id" binding:"required"`
	LPN             string                 `bun:"lpn" json:"lpn" binding:"required"`
	CurrZoneID      *int                   `bun:"curr_zone_id" json:"currZoneID" binding:"required"`
	LastZoneID      *int                   `bun:"last_zone_id" json:"last_zone_id" binding:"required"`
	Direction       string                 `bun:"direction" json:"direction" binding:"required"`
	Confidence      *int                   `bun:"confidence" json:"confidence" binding:"required"`
	CarDetailsID    *int                   `bun:"car_details_id" json:"car_details_id" binding:"required"`
	Extra           map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required" swaggertype:"object"`
}

type ResponsePC struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
	ID              *int   `bun:"id" json:"id"`
	CarDetailsID    *int   `bun:"car_details_id" json:"car_details_id"`
	CameraID        *int   `bun:"camera_id" json:"camera_id"`
	Confidence      *int   `bun:"confidence" json:"confidence"`
	CurrZoneID      *int   `bun:"curr_zone_id" json:"curr_zone_id"`
	LastZoneID      *int   `bun:"last_zone_id" json:"last_zone_id"`
	Direction       string `bun:"direction" json:"direction"`
	LPN             string `bun:"lpn" json:"lpn"`
	TransactionDate string `bun:"transaction_date" json:"transaction_date"`
}

/* type CustomTime struct {
	time.Time
}

func (ct *CustomTime) Scan(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("could not scan type %T into CustomTime", value)
	}
	parsedTime, err := time.Parse("02/01/2006", strVal)
	if err != nil {
		return fmt.Errorf("could not parse date: %v", err)
	}
	ct.Time = parsedTime
	return nil
} */

// Get all present cars
func GetAllPresentExtra(ctx context.Context) ([]PresentCar, error) {
	var cars []PresentCar
	err := Dbg.NewSelect().Model(&cars).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars with Extra: %w", err)
	}
	return cars, nil
}

// Get all present cars
func GetAllPresentCars(ctx context.Context) ([]ResponsePC, error) {
	var Pcars []ResponsePC
	err := Dbg.NewSelect().Model(&Pcars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars: %w", err)
	}
	return Pcars, nil
}

// Get present car by LPN
func GetPresentCarByLPN(ctx context.Context, lpn string) (*PresentCar, error) {
	car := new(PresentCar)
	err := Dbg.NewSelect().Model(car).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting present car by LPN: %w", err)
	}
	return car, nil
}

// Create a new present car
func CreatePresentCar(ctx context.Context, car *PresentCar) error {
	// Insert and get the auto-generated ID from the database
	_, err := Dbg.NewInsert().Model(car).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating present car: %w", err)
	}
	log.Debug().Msgf("New present car added with ID: %d", car.ID)

	return nil
}

// Update a present car by ID and return rows affected
func UpdatePresentCar(ctx context.Context, id int, updates *PresentCar) (int64, error) {
	res, err := Dbg.NewUpdate().Model(updates).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected() // Get the number of rows affected
	log.Debug().Msgf("Updated present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}

// update by LPN
func UpdatePresentCarByLpn(ctx context.Context, lpn string, updates *PresentCar) (int64, error) {
	log.Debug().Str("lpn", lpn).Msgf("Update Present Car by LPN:%v", updates)
	res, err := Dbg.NewUpdate().Model(updates).Where("lpn = ?", lpn).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated present car with LPN: %s, rows affected: %d", lpn, rowsAffected)

	return rowsAffected, nil
}

// Delete a present car by ID and return rows affected
func DeletePresentCar(ctx context.Context, id int) (int64, error) {
	res, err := Dbg.NewDelete().Model(&PresentCar{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
