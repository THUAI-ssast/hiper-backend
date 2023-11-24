package model

// TODO: add more fields
type VerificationCode struct {
}

// TODO
func SaveVerificationCode(code string, email string, expireInMinutes int) error {
	// 保存验证码到 redis，设置过期时间
	return nil
}
