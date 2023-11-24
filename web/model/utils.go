package model

import (
	"errors"
	"fmt"
)

var (
	ErrVerificationCode         = errors.New("verification error")
	ErrVerificationCodeNotExist = fmt.Errorf("%w: verification code not exist", ErrVerificationCode)
	ErrVerificationCodeExpired  = fmt.Errorf("%w: verification code expired", ErrVerificationCode)
)

// TODO
func SaveVerificationCode(code string, email string, expireInMinutes int) error {
	// 保存验证码到 redis，设置过期时间
	return nil
}

func GetVerificationCode(email string) (string, error) {
	// 从 redis 中获取验证码

	// 验证码不存在，返回相应的 error

	// 验证码已过期，返回相应的 error

	// 其他错误，返回 ErrVerificationCode，需要的话可用 Errorf 补充错误信息，需要让上层程序区分的话移到外部 var 定义导出

	// 一切正常，返回验证码
	return "", nil
}
