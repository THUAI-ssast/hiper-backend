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

func paginate(tx *gorm.DB, query QueryParams, result interface{}) (count int64, err error) {
	tx = tx.Select(query.Fields).Where(query.Filter)
	tx = tx.Session(&gorm.Session{})

	if err = tx.Limit(query.Limit).Offset(query.Offset).Find(result).Error; err != nil {
		return 0, err
	}
	if err = tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
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
	err := Rdb.Set(ctx, email, code, time.Duration(expireInMinutes)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetVerificationCode(email string) (string, error) {
	code, err := Rdb.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}
