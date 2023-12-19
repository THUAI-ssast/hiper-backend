package model

import (
	"context"
	"time"
)

var ctx = context.Background()

type QueryParams struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Fields []string
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
