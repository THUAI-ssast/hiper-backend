package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var ctx = context.Background()

type QueryParams struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Fields []string
}

type preloadQuery struct {
	Table   string
	Columns []string
}

func addPreloads(tx *gorm.DB, preloads []preloadQuery) *gorm.DB {
	for _, preload := range preloads {
		tx = tx.Preload(preload.Table, func(db *gorm.DB) *gorm.DB {
			return db.Select(preload.Columns)
		})
	}
	return tx
}

func SaveVerificationCode(code string, email string, expireInMinutes int) error {
	err := rdb.Set(ctx, email, code, time.Duration(expireInMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetVerificationCode(email string) (string, error) {
	code, err := rdb.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}
