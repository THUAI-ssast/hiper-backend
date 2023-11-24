package api

type ErrorFor422 struct {
	Code   ErrorCodeFor422 `json:"code"`
	Field  string          `json:"field"`
	Detail string          `json:"detail"`
}

type ErrorCodeFor422 string

const (
	MissingField  ErrorCodeFor422 = "missing_field"
	Invalid       ErrorCodeFor422 = "invalid"
	AlreadyExists ErrorCodeFor422 = "already_exists"
)
