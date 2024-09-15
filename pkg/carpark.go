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
	ID            int                    `bun:"id,pk,autoincrement"`
	Type          string                 `bun:"type" binding:"required"`
	Name          string                 `bun:"name" binding:"required"`
	Capacity      int                    `bun:"capacity,default:900" binding:"required"`
	Language      map[string]interface{} `bun:"language,type:jsonb" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" binding:"required"`
}

type CarparkReposne struct {
	// Structure
	bun.BaseModel `json:"-" bun:"table:carpark"`
	ID            int                    `bun:"id"`
	Type          string                 `bun:"type"`
	Name          string                 `bun:"name"`
	Capacity      int                    `bun:"capacity,default:900"`
	Language      map[string]interface{} `bun:"language,type:jsonb" `
}

type CarparkOp struct {
	DB *bun.DB
}

func NewCarpark(db *bun.DB) *CarparkOp {
	return &CarparkOp{
		DB: db,
	}
}

func (s *CarparkOp) AddCarpark(ctx context.Context, carpark *Carpark) error {
	log.Debug().Msgf("Adding Carpark....")

	_, err := s.DB.NewInsert().Model(carpark).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding carpark: %w", err)
	}
	log.Debug().Msgf("New carpark added with ID: %d \n", carpark.ID)

	return nil
}

func (s *CarparkOp) GetCarparkByID(ctx context.Context, id int) (*Carpark, error) {
	carpark := new(Carpark)

	log.Debug().Msgf("Fetching carpark  with ID: %d \n", id)

	err := s.DB.NewSelect().Model(carpark).Where("id = ?", id).Scan(ctx)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("carpark with ID %d not found", id)
		}
		return nil, fmt.Errorf("error retrieving carpark with ID %d: %w", id, err)
	}
	return carpark, nil
}

func (s *CarparkOp) GetAllCarparks(ctx context.Context) ([]CarparkReposne, error) {
	var cprk []CarparkReposne
	err := s.DB.NewSelect().Model(&cprk).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all carparks: %w", err)
	}
	return cprk, nil
}

func (s *CarparkOp) GetAllCarparksExtra(ctx context.Context) ([]Carpark, error) {
	var carparks []Carpark
	err := s.DB.NewSelect().Model(&carparks).Order("id ASC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving all carparks with extra data: %w", err)
	}
	return carparks, nil
}

func (s *CarparkOp) UpdateCarpark(ctx context.Context, id int, updatedCarpark *Carpark) (int64, error) {
	log.Debug().Msgf("updating carpark  with ID: %d \n", id)
	result, err := s.DB.NewUpdate().
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

func (s *CarparkOp) DeleteCarpark(ctx context.Context, id int) (int64, error) {
	log.Debug().Msgf("Deleting Carpark with Id %v", id)

	result, err := s.DB.NewDelete().Model((*Carpark)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting carpark with ID %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error fetching rows affected: %w", err)
	}

	return rowsAffected, nil
}
