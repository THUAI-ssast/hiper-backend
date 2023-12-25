package mq

import (
	"fmt"
	"hiper-backend/model"
	"io"
	"os"

	"github.com/dop251/goja"
)

func CreateRuntimeWithJSFile(jsFilePath string) (*goja.Runtime, error) {
	// 打开 JavaScript 文件
	file, err := os.Open(jsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取 JavaScript 文件
	jsCode, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 创建一个新的 goja 运行时
	vm := goja.New()

	// 执行 JavaScript 代码
	_, err = vm.RunString(string(jsCode))
	if err != nil {
		return nil, err
	}

	// 返回运行时
	return vm, nil
}

func CallJSFunction(vm *goja.Runtime, funcName string, args ...interface{}) (goja.Value, error) {
	// 获取 JavaScript 函数
	jsFunc, ok := goja.AssertFunction(vm.Get(funcName))
	if !ok {
		return goja.Undefined(), fmt.Errorf("function %s not found in JavaScript code", funcName)
	}

	// 准备函数参数
	funcArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		funcArgs[i] = vm.ToValue(arg)
	}

	// 调用 JavaScript 函数
	result, err := jsFunc(goja.Undefined(), funcArgs...)
	if err != nil {
		return goja.Undefined(), err
	}

	// 返回结果
	return result, nil
}

// TODO: 在aiassign api中调用
func CallOnAIAssigned(gameID uint, contestant model.Contestant) error {
	return nil
}

func CallOnMatchFinished(matchID uint, replay string) error {
	return nil
}
