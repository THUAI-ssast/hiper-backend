package tests

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/THUAI-ssast/hiper-backend/web/user"
)

func TestBasicFunc(t *testing.T) {
	t.Run("TestGenValidateCode", TestGenValidateCode)
	t.Run("TestIsValidPassword", TestIsValidPassword)
	t.Run("TestIsValidURL", TestIsValidURL)
}

func TestGenValidateCode(t *testing.T) {
	// 测试生成的验证码长度
	for i := 1; i <= 10; i++ {
		s := user.GenValidateCode(i)
		assert.Equal(t, len(s), i)
	}

	// 测试生成的验证码是否只包含数字
	for i := 1; i <= 10; i++ {
		s := user.GenValidateCode(i)
		_, err := strconv.Atoi(s)
		assert.Nil(t, err)
	}

	// 测试生成的验证码是否小于或等于最大值
	for i := 1; i <= 10; i++ {
		s := user.GenValidateCode(i)
		i, err := strconv.Atoi(s)
		assert.Nil(t, err)
		assert.LessOrEqual(t, i, 999999)
	}

	// 测试生成的验证码是否大于或等于最小值
	for i := 1; i <= 10; i++ {
		s := user.GenValidateCode(i)
		i, err := strconv.Atoi(s)
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, i, 0)
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"validPassword1", true},     // 正常情况
		{"invalid", false},           // 错误情况
		{"", false},                  // 边界情况：空字符串
		{"12345678", true},           // 边界情况：最小长度
		{"1234567890123456", true},   // 边界情况：最大长度
		{"12345678901234567", false}, // 边界情况：超过最大长度
		{"password!", true},          // 包含特殊字符
		{"PASSWORD", true},           // 全大写
		{"password", true},           // 全小写
		{"Password", true},           // 大小写混合
		{"password1", true},          // 包含数字
		{"passwordpassword", true},   // 16个字符
		{"pass", false},              // 少于8个字符
		{"pass word", false},         // 包含空格
		{"password中文", false},        // 包含非ASCII字符
	}

	for _, tt := range tests {
		got := user.IsValidPassword(tt.password)
		assert.Equal(t, tt.want, got)
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"http://example.com", true},                   // 正常情况
		{"https://example.com", true},                  // 正常情况
		{"http://www.example.com", true},               // 正常情况
		{"https://www.example.com", true},              // 正常情况
		{"http://example.com/path", true},              // 正常情况
		{"https://example.com/path", true},             // 正常情况
		{"http://example.com/path?query=param", true},  // 正常情况
		{"https://example.com/path?query=param", true}, // 正常情况
		{"http://localhost", true},                     // 正常情况
		{"https://localhost", true},                    // 正常情况
		{"http://127.0.0.1", true},                     // 正常情况
		{"https://127.0.0.1", true},                    // 正常情况
		{"http://[::1]", true},                         // 正常情况
		{"https://[::1]", true},                        // 正常情况
		{"invalid", false},                             // 错误情况
		{"", true},                                     // 边界情况：空字符串
	}

	for _, tt := range tests {
		got := user.IsValidURL(tt.url)
		assert.Equal(t, tt.want, got)
	}
}
