package pkg

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

type Settings struct {
	bun.BaseModel      `json:"-" bun:"table:settings"`
	CarParkID          int    `bun:"carpark_id" json:"carpark_id" binding:"required"`
	CarParkName        string `bun:"carpark_name" binding:"required" json:"carpark_name"`
	AppLogo            string `bun:"app_logo,type:bytea" binding:"required" json:"app_logo"`
	DefaultLang        string `bun:"default_lang" binding:"required" json:"default_lang" default:"EN"`
	TimeOutScreenKisok *int   `bun:"timeout_screenKiosk" binding:"required" json:"timeout_screenKiosk"`
	FycCleanCron       string `bun:"fyc_clean_cron" binding:"required" json:"fyc_clean_cron"`
	CountingCleanCron  string `bun:"couting_clean_cron" binding:"required" json:"couting_clean_cron"`
}

func CreateSettings(ctx context.Context, settings *Settings) error {
	_, err := Dbg.NewInsert().Model(settings).Exec(ctx)
	return err
}

// GetSettings fetches a settings entry by CarParkID
func GetSettings(ctx context.Context, carParkID int) (*Settings, error) {
	settings := new(Settings)
	err := Dbg.NewSelect().Model(settings).Where("carpark_id = ?", carParkID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func GetAllSettings(ctx context.Context) ([]Settings, error) {
	var settings []Settings
	err := Dbg.NewSelect().Model(&settings).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateSettings updates a settings entry by CarParkID
func UpdateSettings(ctx context.Context, settings *Settings) error {
	res, err := Dbg.NewUpdate().Model(settings).Where("carpark_id = ?", settings.CarParkID).Exec(ctx)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
