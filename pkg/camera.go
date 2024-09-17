package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Camera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int                    `bun:"id,pk,autoincrement" json:"ID"`
	CamName       string                 `bun:"cam_name" json:"cam_name" binding:"required"`
	CamType       string                 `bun:"cam_type" json:"cam_type" binding:"required"`
	CamIP         string                 `bun:"cam_ip" json:"cam_ip" binding:"required"`
	CamPORT       string                 `bun:"cam_port" json:"cam_port" binding:"required"`
	CamUser       string                 `bun:"cam_user" json:"cam_user" binding:"required"`
	CamPass       string                 `bun:"cam_password" json:"cam_password" binding:"required"`
	ZoneIdIn      *int                   `bun:"zone_id_in" json:"zone_id_in" binding:"required"`
	ZoneIdOut     *int                   `bun:"zone_id_out" json:"zone_id_out" binding:"required"`
	Direction     string                 `bun:"direction" json:"direction" binding:"required"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required"`
}

type ResponseCamera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int    `bun:"id,pk,autoincrement"`
	CamName       string `bun:"cam_name"`
	CamType       string `bun:"cam_type"`
	CamIP         string `bun:"cam_ip"`
	CamPORT       string `bun:"cam_port" `
	CamUser       string `bun:"cam_user" `
	CamPass       string `bun:"cam_password" `
	ZoneIdIn      *int   `bun:"zone_id_in" `
	ZoneIdOut     *int   `bun:"zone_id_out" `
	Direction     string `bun:"direction" `
}

// Get all camera with extra data
func GetAllCameraExtra(ctx context.Context, db *bun.DB) ([]Camera, error) {
	var camera []Camera
	err := db.NewSelect().Model(&camera).Column().Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Camera with Extra Data: %w", err)
	}
	return camera, nil
}

// Get all camera
func GetAllCamera(ctx context.Context, db *bun.DB) ([]ResponseCamera, error) {
	var cam []ResponseCamera
	err := db.NewSelect().Model(&cam).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}
	return cam, nil
}

// Gt camera by id
func GetCameraByID(ctx context.Context, db *bun.DB, id int) (*Camera, error) {
	cm := new(Camera)
	err := db.NewSelect().Model(cm).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting camera by id : %w", err)
	}
	return cm, nil
}

// create a new camera
func CreateCamera(ctx context.Context, db *bun.DB, newcam *Camera) error {
	// Insert and get the auto-generated ID from the database
	_, err := db.NewInsert().Model(newcam).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a camera : %w", err)
	}
	log.Debug().Msgf("New camera added with ID: %d", newcam.ID)

	return nil
}

// Update a camera by ID
func UpdateCamera(ctx context.Context, db *bun.DB, cam_id int, updates *Camera) (int64, error) {
	res, err := db.NewUpdate().Model(updates).Where("zone_id = ?", cam_id).ExcludeColumn("id").Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating camera with id %d: %w", cam_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated camera with ID: %d, rows affected: %d", cam_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a zone img by ID
func DeleteCamera(ctx context.Context, db *bun.DB, id int) (int64, error) {
	res, err := db.NewDelete().Model(&Camera{}).Where("ID = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting Camera with id %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted Camera with ID: %d, rows affected: %d", id, rowsAffected)

	return rowsAffected, nil
}
