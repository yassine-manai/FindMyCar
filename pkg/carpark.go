package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Carpark struct {
	// Structure
	bun.BaseModel `json:"-" bun:"table:carpark"`
	ID            int                    `bun:"id,pk,autoincrement" json:"id"`
	Type          string                 `bun:"type" binding:"required" json:"type"`
	Name          string                 `bun:"name" binding:"required" json:"name"`
	Capacity      *int                   `bun:"capacity,default:900" json:"capacity" binding:"required"`
	Language      map[string]interface{} `bun:"language,type:jsonb" json:"language"  binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type CarparkReposne struct {
	// Structure
	bun.BaseModel `json:"-" bun:"table:carpark"`
	ID            int                    `bun:"id" json:"id"`
	Type          string                 `bun:"type" json:"type"`
	Name          string                 `bun:"name" json:"name"`
	Capacity      *int                   `bun:"capacity" json:"capacity"`
	Language      map[string]interface{} `bun:"language" json:"language"`
}

func AddCarpark(ctx context.Context, carpark *Carpark) error {
	log.Debug().Msgf("Adding Carpark....")

	_, err := Dbg.NewInsert().Model(carpark).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding carpark: %w", err)
	}
	log.Debug().Msgf("New carpark added with ID: %d \n", carpark.ID)

	return nil
}

func GetCarparkByID(ctx context.Context, id int) (*Carpark, error) {
	carpark := new(Carpark)

	log.Debug().Msgf("Fetching carpark  with ID: %d \n", id)

	err := Dbg.NewSelect().Model(carpark).Where("id = ?", id).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("carpark with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving carpark with ID %d: %w", id, err)
	}
	return carpark, nil
}

func GetAllCarparks(ctx context.Context) ([]CarparkReposne, error) {
	var cprk []CarparkReposne
	err := Dbg.NewSelect().Model(&cprk).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all carparks: %w", err)
	}
	return cprk, nil
}

func GetAllCarparksExtra(ctx context.Context) ([]Carpark, error) {
	var carparks []Carpark
	err := Dbg.NewSelect().Model(&carparks).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all carparks with extra data: %w", err)
	}
	return carparks, nil
}

func UpdateCarpark(ctx context.Context, id int, updatedCarpark *Carpark) (int64, error) {
	log.Debug().Msgf("updating carpark  with ID: %d \n", id)
	result, err := Dbg.NewUpdate().
		Model(updatedCarpark).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating carpark with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}

func DeleteCarpark(ctx context.Context, id int) (int64, error) {
	log.Debug().Msgf("Deleting Carpark with Id %v", id)

	result, err := Dbg.NewDelete().Model((*Carpark)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting carpark with ID %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}
