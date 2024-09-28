package pkg

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type Camera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int                    `bun:"id,autoincrement" json:"id"`
	CamID         int                    `bun:"cam_id,pk" json:"cam_id" binding:"required"`
	CamName       string                 `bun:"cam_name" json:"cam_name" binding:"required"`
	CamType       string                 `bun:"cam_type" json:"cam_type" binding:"required"`
	CamIP         string                 `bun:"cam_ip" json:"cam_ip" binding:"required"`
	CamPORT       int                    `bun:"cam_port" json:"cam_port" binding:"required"`
	CamUser       string                 `bun:"cam_user" json:"cam_user" binding:"required"`
	CamPass       string                 `bun:"cam_password" json:"cam_password" binding:"required"`
	ZoneIdIn      *int                   `bun:"zone_id_in" json:"zone_id_in" binding:"required"`
	ZoneIdOut     *int                   `bun:"zone_id_out" json:"zone_id_out" binding:"required"`
	Direction     string                 `bun:"direction" json:"direction" binding:"required"`
	IsEnabled     bool                   `bun:"is_enabled,type:bool" json:"is_enabled" default:"false"`
	IsDeleted     bool                   `bun:"is_deleted,type:bool" json:"is_deleted" default:"false"`
	Extra         map[string]interface{} `bun:"extra,type:jsonb" json:"extra" binding:"required" swaggertype:"object"`
}

type ResponseCamera struct {
	bun.BaseModel `json:"-" bun:"table:camera"`
	ID            int    `bun:"id" json:"id"`
	CamID         int    `bun:"cam_id" json:"cam_id"`
	CamName       string `bun:"cam_name" json:"cam_name"`
	CamType       string `bun:"cam_type" json:"cam_type"`
	CamIP         string `bun:"cam_ip" json:"cam_ip"`
	CamPORT       int    `bun:"cam_port" json:"cam_port" `
	CamUser       string `bun:"cam_user"  json:"cam_user"`
	CamPass       string `bun:"cam_password" json:"cam_password" `
	ZoneIdIn      *int   `bun:"zone_id_in"  json:"zone_id_in"`
	ZoneIdOut     *int   `bun:"zone_id_out" json:"zone_id_out" `
	Direction     string `bun:"direction" json:"direction" `
	IsEnabled     bool   `bun:"is_enabled" json:"is_enabled" `
	IsDeleted     bool   `bun:"is_deleted" json:"is_deleted" `
}

// Get all camera Data
func GetDataCamera(ctx context.Context) ([]Camera, error) {
	var camData []Camera
	err := Dbg.NewSelect().Model(&camData).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}
	return camData, nil
}

// Get all camera with extra data
func GetAllCameraExtra(ctx context.Context) ([]Camera, error) {
	var camera []Camera
	err := Dbg.NewSelect().Model(&camera).
		Where("is_deleted = ?", false).
		Column().
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all Camera with Extra Data: %w", err)
	}
	return camera, nil
}

// Get all camera
func GetAllCamera(ctx context.Context) ([]ResponseCamera, error) {
	var cam []ResponseCamera
	err := Dbg.NewSelect().
		Model(&cam).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting all cameras : %w", err)
	}
	return cam, nil
}

// Get camera by ID
func GetCameraByID(ctx context.Context, id int) (*Camera, error) {
	var cam Camera
	err := Dbg.NewSelect().
		Model(&cam).
		Where("cam_id = ?", id).
		Where("is_deleted = ?", false).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id %d: %w", id, err)
	}

	return &cam, nil
}

func GetCameraListEnabled(ctx context.Context) ([]Camera, error) {
	var camera []Camera
	err := Dbg.NewSelect().
		Model(&camera).
		Where("is_enabled = ?", true).
		Scan(ctx, &camera)

	if err != nil {
		return nil, fmt.Errorf("error getting Enabled camera List: %w", err)
	}

	return camera, nil
}

func GetCameraListDeleted(ctx context.Context) ([]Camera, error) {
	var cameras []Camera
	err := Dbg.NewSelect().
		Model(&cameras).
		Where("is_deleted = ?", true).
		Scan(ctx, &cameras)
	if err != nil {
		return nil, fmt.Errorf("error fetching deleted cameras: %w", err)
	}
	return cameras, nil
}
func GetCameraEnabledByID(ctx context.Context, id int) (*Camera, error) {
	var cam Camera
	err := Dbg.NewSelect().
		Model(&cam).
		Column("is_enabled").
		Where("cam_id = ?", id).
		Where("is_deleted = ?", false).
		Where("is_enabled = ?", true).
		Scan(ctx, &cam)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id: %w", err)
	}

	return &cam, nil
}

func GetCameraDeletedByID(ctx context.Context, id int) (*Camera, error) {
	var cam Camera
	err := Dbg.NewSelect().
		Model(&cam).
		Column("is_deleted").
		Where("is_deleted = ?", true).
		Where("cam_id = ?", id).
		Scan(ctx, &cam)

	if err != nil {
		return nil, fmt.Errorf("error getting camera by id: %w", err)
	}

	return &cam, nil
}

// CreateCamera creates a new camera if the ID is not found, and checks if the database contains data.
func CreateCamera(ctx context.Context, newcam *Camera) error {

	// Proceed with the camera creation if no camera exists
	_, err := Dbg.NewInsert().Model(newcam).Returning("cam_id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating a new camera: %w", err)
	}

	log.Debug().Msgf("New camera added with ID: %d", newcam.ID)
	return nil
}

// Update a camera by ID
func UpdateCamera(ctx context.Context, cam_id int, updates *Camera) (int64, error) {
	res, err := Dbg.NewUpdate().Model(updates).Where("cam_id = ?", cam_id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating camera with id %d: %w", cam_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated camera with ID: %d, rows affected: %d", cam_id, rowsAffected)

	return rowsAffected, nil
}

// Delete a camera by ID (soft delete: sets is_deleted to true)
func DeleteCamera(ctx context.Context, Cam_id int) (int64, error) {
	res, err := Dbg.NewUpdate().
		Model(&Camera{}).
		Set("is_deleted = ?", true).
		Where("cam_id = ?", Cam_id).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error deleting Camera with id %d: %w", Cam_id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Soft-deleted Camera with ID: %d, rows affected: %d", Cam_id, rowsAffected)

	return rowsAffected, nil
}

func ChangeCameraState(ctx context.Context, camID int, newState bool) (int64, error) {
	existingCam := new(Camera)

	err := Dbg.NewSelect().
		Model(existingCam).
		Where("cam_id = ?", camID).
		Where("is_deleted = ?", false).
		Column("is_enabled").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return 0, fmt.Errorf("error retrieving Camera state with id %d: %w", camID, err)
	}

	if existingCam.IsEnabled == newState {
		stateMessage := "already"
		if !newState {
			stateMessage = "disabled"
		} else {
			stateMessage = "enabled"
		}
		return 0, fmt.Errorf("camera with id %d is already %s", camID, stateMessage)
	}

	// Step 2: Change the state since it's different
	res, err := Dbg.NewUpdate().
		Model(&Camera{}).
		Set("is_enabled = ?", newState).
		Where("is_deleted = ?", false).
		Where("cam_id = ?", camID).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("error changing Camera state with id %d: %w", camID, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Changed Camera State with ID: %d, rows affected: %d", camID, rowsAffected)

	return rowsAffected, nil
}
