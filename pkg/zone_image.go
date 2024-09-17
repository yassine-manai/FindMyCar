package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	ZoneID        *int                   `bun:"zone_id" json:"zone_id" binding:"required"`
	Lang          string                 `bun:"lang" json:"lang" binding:"required"`
	ImageSm       []byte                 `bun:"image_s" json:"image_s" binding:"required"`
	ImageLg       []byte                 `bun:"image_l" json:"image_l" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int    `bun:"id,pk,autoincrement"`
	ZoneID        *int   `bun:"zone_id"`
	Lang          string `bun:"lang"`
	ImageSm       []byte `bun:"image_s"`
	ImageL        []byte `bun:"image_l"`
}

// Get all Zones with extra data
func GetAllZoneImageExtra(ctx context.Context, db *bun.DB) ([]ImageZone, error) {
	var zoneImage []ImageZone
	err := db.NewSelect().Model(&zoneImage).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Image Zones with Extra Data: %w", err)
	}
	return zoneImage, nil
}

// Get all zone
func GetAllZoneImage(ctx context.Context, db *bun.DB) ([]ResponseImageZone, error) {
	var EZI []ResponseImageZone
	err := db.NewSelect().Model(&EZI).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones Images : %w", err)
	}
	return EZI, nil
}

// Gt zone by id
func GetZoneImageByID(ctx context.Context, db *bun.DB, zone_id int) (*ImageZone, error) {
	zoneImg := new(ImageZone)
	err := db.NewSelect().Model(zoneImg).Where("zone_id = ?", zone_id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by id : %w", err)
	}
	return zoneImg, nil
}

// create a new zone
func CreateZoneImage(ctx context.Context, db *bun.DB, zoneImg *ImageZone) error {
	// Insert and get the auto-generated ID from the database
	_, err := db.NewInsert().Model(zoneImg).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a Zone Image : %w", err)
	}
	log.Debug().Msgf("New zone image added with ID: %d", zoneImg.ID)

	return nil
}

// Update a zone img by ID
func UpdateZoneImage(ctx context.Context, db *bun.DB, zone_id int, updates *ImageZone) (int64, error) {
	res, err := db.NewUpdate().Model(updates).Where("zone_id = ?", zone_id).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating Zone Image with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated Zone Image with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a zone img by ID
func DeleteZoneImage(ctx context.Context, db *bun.DB, id int) (int64, error) {
	res, err := db.NewDelete().Model(&ImageZone{}).Where("ID = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone Image with id %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted Zone Image with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
