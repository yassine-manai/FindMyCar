package functions

import (
	"context"
	"encoding/base64"
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
	for model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("error creating table for model %T: %w", model, err)
		}
	}
	return nil
}

func Base64ToBytea(base64Str string) (string, error) {
	// Decode the base64-encoded string
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err.Error(), fmt.Errorf("failed to decode base64 string: %w", err)
	}
	return string(data), nil
}

func ByteaToBase64(data []byte) string {
	// Encode the byte slice to a base64 string
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64ToByteArray(base64Image string) ([]byte, error) {
	// Decode the base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Return the byte array (bytea)
	return imageData, nil
}
