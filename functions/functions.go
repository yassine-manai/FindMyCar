package functions

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func Contains(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func CreateTables(ctx context.Context, db *bun.DB, models []interface{}) error {
	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("error creating table for model %T: %w", model, err)
		}
	}
	return nil
}
