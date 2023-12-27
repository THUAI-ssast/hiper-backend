package basecontest

import (
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
)

func ConvertStruct(s interface{}) map[string]interface{} {
	var m map[string]interface{}
	if structs.IsStruct(s) {
		m = structs.Map(s)
	} else if vm, ok := s.(map[string]interface{}); ok {
		m = vm
	}
	for k, v := range m {
		snake := strcase.ToSnake(k)
		if snake != k {
			delete(m, k)
			m[snake] = v
		}

		// 如果值是 struct，递归地转换它
		if structs.IsStruct(v) {
			m[snake] = ConvertStruct(v)
		} else if vm, ok := v.(map[string]interface{}); ok {
			// 如果值是 map，也递归地转换它
			m[snake] = ConvertStruct(vm)
		}
	}
	return m
}
