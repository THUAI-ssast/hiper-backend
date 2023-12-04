package model

import (
	"errors"
	"fmt"
  "context"
	"time"
)

var ctx = context.Background()

var (
	ErrVerificationCode         = errors.New("verification error")
	ErrVerificationCodeNotExist = fmt.Errorf("%w: verification code not exist", ErrVerificationCode)
	ErrVerificationCodeExpired  = fmt.Errorf("%w: verification code expired", ErrVerificationCode)
)

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
