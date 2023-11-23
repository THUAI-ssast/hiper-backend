package user

import (
	"hiper-backend/mail"
	"hiper-backend/model"
)

// TODO
func SendVerificationCode(email string) error {
	// 生成验证码
	var code string

	// 发送邮件
	if err := mail.SendVerificationCode(code, email); err != nil {
		return err
	}

	// 保存验证码
	if err := model.SaveVerificationCode(code, email, 5); err != nil {
		return err
	}
	return nil
}
