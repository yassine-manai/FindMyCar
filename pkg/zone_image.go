package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type ImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int                    `bun:"id,autoincrement" json:"ID"`
	ZoneID        int                    `bun:"zone_id" json:"zone_id" binding:"required"`
	Lang          string                 `bun:"lang" json:"lang" binding:"required"`
	ImageSm       []byte                 `bun:"image_s" json:"image_s" binding:"required"`
	ImageLg       []byte                 `bun:"image_l" json:"image_l" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseImageZone struct {
	bun.BaseModel `json:"-" bun:"table:zone_images"`
	ID            int    `bun:"id,autoincrement"`
	ZoneID        int    `bun:"zone_id"`
	Lang          string `bun:"lang"`
	ImageSm       []byte `bun:"image_s"`
	ImageL        []byte `bun:"image_l"`
}

type ImageZoneOp struct {
	DB *bun.DB
}

func NewImageZone(db *bun.DB) *ImageZoneOp {
	return &ImageZoneOp{
		DB: db,
	}
}

// Get all Zones with extra data
func (znimop *ImageZoneOp) GetAllZoneImageExtra(ctx context.Context) ([]ImageZone, error) {
	var zoneImage []ImageZone
	err := znimop.DB.NewSelect().Model(&zoneImage).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Image Zones with Extra Data: %w", err)
	}
	return zoneImage, nil
}

// Get all zone
func (znimop *ImageZoneOp) GetAllZone(ctx context.Context) ([]ResponseImageZone, error) {
	var EZI []ResponseImageZone
	err := znimop.DB.NewSelect().Model(&EZI).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Zones Images : %w", err)
	}
	return EZI, nil
}

// Gt zone by id
func (znimop *ImageZoneOp) GetZoneImageByID(ctx context.Context, zone_id int) (*ImageZone, error) {
	zoneImg := new(ImageZone)
	err := znimop.DB.NewSelect().Model(zoneImg).Where("zone_id = ?", zone_id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Zone Image by id : %w", err)
	}
	return zoneImg, nil
}

// create a new zone
func (znimop *ImageZoneOp) CreateZoneImage(ctx context.Context, zoneImg *ImageZone) error {
	// Insert and get the auto-generated ID from the database
	_, err := znimop.DB.NewInsert().Model(zoneImg).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a Zone Image : %w", err)
	}
	log.Debug().Msgf("New zone image added with ID: %d", zoneImg.ID)

	return nil
}

// Update a zone img by ID
func (znimop *ImageZoneOp) UpdateZoneImage(ctx context.Context, zone_id int, updates *ImageZone) (int64, error) {
	res, err := znimop.DB.NewUpdate().Model(updates).Where("zone_id = ?", zone_id).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating Zone Image with id %d: %w", zone_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated Zone Image with ID: %d, rows affected: %d", zone_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a zone img by ID
func (znimop *ImageZoneOp) DeleteZoneImage(ctx context.Context, id int) (int64, error) {
	res, err := znimop.DB.NewDelete().Model(&ImageZone{}).Where("ID = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Zone Image with id %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted Zone Image with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
