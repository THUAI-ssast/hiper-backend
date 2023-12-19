package tests

import (
	"hiper-backend/user"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserBasicFunc(t *testing.T) {
	t.Run("TestGenValidateCode", TestGenValidateCode)
	t.Run("TestIsValidPassword", TestIsValidPassword)
	t.Run("TestIsValidURL", TestIsValidURL)
}

func TestGenValidateCode(t *testing.T) {
	s := user.GenValidateCode(10)
	assert.Equal(t, len(s), 10)
	i, err := strconv.Atoi(s)
	assert.Nil(t, err)
	assert.LessOrEqual(t, i, 9999999999)
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"validPassword1", true},
		{"invalid", false},
		{"", false},
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
		{"http://example.com", true},
		{"invalid", false},
		{"", true},
	}

	for _, tt := range tests {
		got := user.IsValidURL(tt.url)
		assert.Equal(t, tt.want, got)
	}
}
