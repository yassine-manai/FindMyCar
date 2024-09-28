package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Zone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int                    `bun:"id,autoincrement" json:"id"`
	ZoneID        int                    `bun:"zone_id,pk" json:"zone_id" binding:"required"`
	Name          string                 `bun:"name" json:"name" binding:"required" `
	MaxCapacity   *int                   `bun:"max_capacity" json:"max_capacity" binding:"required"`
	FreeCapacity  *int                   `bun:"free_capacity" json:"free_capacity" binding:"required"`
	LastUpdated   string                 `bun:"last_update,type:date" json:"last_update" binding:"required"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"is_enabled" default:"false"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"is_deleted" default:"false"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required" swaggertype:"object"`
}

type ResponseZone struct {
	bun.BaseModel `json:"-" bun:"table:zone"`
	ID            int    `bun:"id" json:"id"`
	ZoneID        *int   `bun:"zone_id" json:"zone_id"`
	Name          string `bun:"name" json:"name"`
	MaxCapacity   *int   `bun:"max_capacity" json:"max_capacity"`
	FreeCapacity  *int   `bun:"free_capacity" json:"free_capacity"`
	LastUpdated   string `bun:"last_update" json:"last_update"`
	IsEnabled     bool   `bun:"is_enabled" json:"is_enabled" default:"false"`
	IsDeleted     bool   `bun:"is_deleted" json:"is_deleted" default:"false"`
}

func GetZoneData(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := Dbg.NewSelect().Model(&zone).Scan(ctx)
	if err != nil {

		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}
	return zone, nil
}

// Get all Zones with extra data
func GetAllZoneExtra(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Column().
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all Zones with Extra: %w", err)
	}
	return zone, nil
}

// Get all zone
func GetAllZone(ctx context.Context) ([]ResponseZone, error) {
	var EZ []ResponseZone
	err := Dbg.NewSelect().
		Model(&EZ).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all zones : %w", err)
	}
	return EZ, nil
}

// Gt zone by id
func GetZoneByID(ctx context.Context, id int) (*Zone, error) {
	var zone Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", false).
		Where("zone_id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id : %w", err)
	}
	return &zone, nil
}

func GetZoneListEnabled(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Where("is_enabled = ?", true).
		Scan(ctx, &zone)

	if err != nil {
		return nil, fmt.Errorf("error getting enabled zone list : %w", err)
	}

	return zone, nil
}

func GetZoneListDeleted(ctx context.Context) ([]Zone, error) {
	var zone []Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Where("is_deleted = ?", true).
		Scan(ctx, &zone)
	if err != nil {
		return nil, fmt.Errorf("error fetching deleted zones: %w", err)
	}
	return zone, nil
}

func GetZoneEnabledByID(ctx context.Context, id int) (*Zone, error) {
	var zone Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Column("is_enabled").
		Where("zone_id = ?", id).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx, &zone)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id: %w", err)
	}

	return &zone, nil
}

func GetZoneDeletedByID(ctx context.Context, id int) (*Zone, error) {
	var zone Zone
	err := Dbg.NewSelect().
		Model(&zone).
		Column("is_deleted").
		Where("is_deleted = ?", true).
		Where("zone_id = ?", id).
		Scan(ctx, &zone)

	if err != nil {
		return nil, fmt.Errorf("error getting zone by id: %w", err)
	}

	return &zone, nil
}

// create a new zone
func CreateZone(ctx context.Context, zone *Zone) error {
	// Insert and get the auto-generated ID from the database
	_, err := Dbg.NewInsert().Model(zone).Returning("zone_id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a zone : %w", err)
	}
	log.Debug().Msgf("New zone added with ID: %d", zone.ID)

	return nil
}

// Update a zone by ID
func UpdateZone(ctx context.Context, zone_id int, updates *Zone) (int64, error) {
	res, err := Dbg.NewUpdate().
		Model(updates).
		Where("zone_id = ?", zone_id).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating zone with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated zone with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a zone by ID
func DeleteZone(ctx context.Context, zone_id int) (int64, error) {
	var zone Zone
	res, err := Dbg.NewUpdate().
		Model(zone).
		Set("is_deleted = ?", true).
		Where("zone_id = ?", zone_id).
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted zone with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}

func ChangeZoneState(ctx context.Context, zoneID int, newState bool) (int64, error) {
	existingZone := new(Zone)

	err := Dbg.NewSelect().
		Model(existingZone).
		Where("zone_id = ?", zoneID).
		Where("is_deleted = ?", false).
		Column("is_enabled").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0, fmt.Errorf("error retrieving Zone state with id %d: %w", zoneID, err)
	}

	if existingZone.IsEnabled == newState {
		stateMessage := "already"
		if !newState {
			stateMessage = "disabled"
		} else {
			stateMessage = "enabled"
		}
		return 0, fmt.Errorf("Zone with id %d is already %s", zoneID, stateMessage)
	}

	// Step 2: Change the state since it's different
	res, err := Dbg.NewUpdate().
		Model(&Zone{}).
		Set("is_enabled = ?", newState).
		Where("is_deleted = ?", false).
		Where("zone_id = ?", zoneID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error changing Zone state with id %d: %w", zoneID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Changed Zone State with ID: %d, rows affected: %d", zoneID, rowsAffected)

	return rowsAffected, nil
}
