package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Zone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	ZoneID        int                    `bun:"zone_id" json:"zone_id" binding:"required"`
	MaxCapacity   int                    `bun:"max_capacity" json:"max_capacity" binding:"required"`
	Present       int                    `bun:"present" json:"present" binding:"required"`
	Name          map[string]interface{} `bun:"name,type:json" json:"name" binding:"required"`
	Description   string                 `bun:"description" json:"description" binding:"required"`
	CarParkID     int                    `bun:"carpark_id" json:"carpark_id" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseZone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	ZoneID        int                    `bun:"zone_id"`
	MaxCapacity   int                    `bun:"max_capacity"`
	Present       int                    `bun:"present"`
	Name          map[string]interface{} `bun:"name"`
	Description   string                 `bun:"description"`
	CarParkID     int                    `bun:"carpark_id"`
}

type ZoneOp struct {
	DB *bun.DB
}

func NewZone(db *bun.DB) *ZoneOp {
	return &ZoneOp{
		DB: db,
	}
}

// Get all Zones with extra data
func (znop *ZoneOp) GetAllZoneExtra(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := znop.DB.NewSelect().Model(&zone).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}
	return zone, nil
}

// Get all zone
func (znop *ZoneOp) GetAllZone(ctx context.Context) ([]ResponseZone, error) {
	var EZ []ResponseZone
	err := znop.DB.NewSelect().Model(&EZ).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all zones : %w", err)
	}
	return EZ, nil
}

// Gt zone by id
func (znop *ZoneOp) GetZoneByID(ctx context.Context, id int) (*Zone, error) {
	zone := new(Zone)
	err := znop.DB.NewSelect().Model(zone).Where("zone_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}
	return zone, nil
}

// create a new zone
func (znop *ZoneOp) CreateZone(ctx context.Context, zone *Zone) error {
	// Insert and get the auto-generated ID from the database
	_, err := znop.DB.NewInsert().Model(zone).Returning("zone_id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a zone : %w", err)
	}
	log.Debug().Msgf("New zone added with ID: %d", zone.ID)

	return nil
}

// Update a zone by ID
func (znop *ZoneOp) UpdateZone(ctx context.Context, zone_id int, updates *Zone) (int64, error) {
	res, err := znop.DB.NewUpdate().Model(updates).Where("zone_id = ?", zone_id).ExcludeColumn("ID").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating zone with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated zone with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a zone by ID
func (znop *ZoneOp) DeleteZone(ctx context.Context, zone_id int) (int64, error) {
	res, err := znop.DB.NewDelete().Model(&Zone{}).Where("zone_id = ?", zone_id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted zone with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}
