package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type PresentCar struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
	ID              *int                   `bun:"id,pk,autoincrement" json:"ID"`
	CarDetailsID    *int                   `bun:"carDetailsID" json:"carDetailsID" binding:"required"`
	CameraID        *int                   `bun:"cameraID" json:"cameraID" binding:"required"`
	Confidence      *int                   `bun:"confidence" json:"confidence" binding:"required"`
	CurrZoneID      *int                   `bun:"currZoneID" json:"currZoneID" binding:"required"`
	Direction       string                 `bun:"direction" json:"direction" binding:"required"`
	LastZoneID      *int                   `bun:"lastZoneID" json:"lastZoneID" binding:"required"`
	LPN             string                 `bun:"lpn" json:"lpn" binding:"required"`
	TransactionDate string                 `bun:"transactionDate" json:"transactionDate" binding:"required"`
	Extra           map[string]interface{} `bun:"extra,type:jsonb" time_format:"01/01/2005" json:"extra" binding:"required"`
}

type ResponsePC struct {
	bun.BaseModel   `json:"-" bun:"table:presentcar"`
	ID              *int   `bun:"id,pk,autoincrement"`
	CarDetailsID    *int   `bun:"carDetailsID"`
	CameraID        *int   `bun:"cameraID"`
	Confidence      *int   `bun:"confidence"`
	CurrZoneID      *int   `bun:"currZoneID"`
	Direction       string `bun:"direction"`
	LastZoneID      *int   `bun:"lastZoneID"`
	LPN             string `bun:"lpn"`
	TransactionDate string `bun:"transactionDate"`
}

type CustomTime struct {
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
}

// Get all present cars
func GetAllPresentExtra(ctx context.Context, db *bun.DB) ([]PresentCar, error) {
	var cars []PresentCar
	err := db.NewSelect().Model(&cars).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars with Extra: %w", err)
	}
	return cars, nil
}

// Get all present cars
func GetAllPresentCars(ctx context.Context, db *bun.DB) ([]ResponsePC, error) {
	var Pcars []ResponsePC
	err := db.NewSelect().Model(&Pcars).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all present cars: %w", err)
	}
	return Pcars, nil
}

// Get present car by LPN
func GetPresentCarByLPN(ctx context.Context, db *bun.DB, lpn string) (*PresentCar, error) {
	car := new(PresentCar)
	err := db.NewSelect().Model(car).Where("lpn = ?", lpn).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting present car by LPN: %w", err)
	}
	return car, nil
}

// Create a new present car
func CreatePresentCar(ctx context.Context, db *bun.DB, car *PresentCar) error {
	// Insert and get the auto-generated ID from the database
	_, err := db.NewInsert().Model(car).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating present car: %w", err)
	}
	log.Debug().Msgf("New present car added with ID: %d", car.ID)

	return nil
}

// Update a present car by ID and return rows affected
func UpdatePresentCar(ctx context.Context, db *bun.DB, id int, updates *PresentCar) (int64, error) {
	res, err := db.NewUpdate().Model(updates).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected() // Get the number of rows affected
	log.Debug().Msgf("Updated present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}

// update by LPN
func UpdatePresentCarByLpn(ctx context.Context, db *bun.DB, lpn string, updates *PresentCar) (int64, error) {
	log.Debug().Str("lpn", lpn).Msgf("Update Present Car by LPN:%v", updates)
	res, err := db.NewUpdate().Model(updates).Where("lpn = ?", lpn).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated present car with LPN: %s, rows affected: %d", lpn, rowsAffected)

	return rowsAffected, nil
}

// Delete a present car by ID and return rows affected
func DeletePresentCar(ctx context.Context, db *bun.DB, id int) (int64, error) {
	res, err := db.NewDelete().Model(&PresentCar{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting present car: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted present car with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
