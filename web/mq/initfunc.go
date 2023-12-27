package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/dop251/goja"

	"github.com/THUAI-ssast/hiper-backend/web/model"
)

var Ctx_callback = context.Background()

// 定义全局 map
var (
	baseContestIDToRuntime = make(map[uint]*goja.Runtime)
	runtimeToBaseContestID = make(map[*goja.Runtime]uint)
	mutex                  = &sync.Mutex{} // 用于保护 map 的并发访问
)

func InitMq() {
	go ListenMsgForMatchFinished(Ctx_callback, "match_result")
}

func InitGameMq(baseContestID uint) {
	SetCreateMatch(baseContestID)
	SetGetContestantsByRanking(baseContestID)
	SetUpdateContestant(baseContestID)
	SendBuildGameLogicMsg(Ctx_callback, baseContestID)
}

func SetGoFuncForJS(baseContestID uint, funcName string, goFunc func(goja.FunctionCall) goja.Value) error {
	// 打开 JavaScript 文件
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		return err
	}
	script := baseContest.Script

	// 查找 map 中对应的 runtime
	mutex.Lock()
	vm, exists := baseContestIDToRuntime[baseContestID]
	mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			return err
		}

		// 再次查找 runtime
		mutex.Lock()
		vm = baseContestIDToRuntime[baseContestID]
		mutex.Unlock()
	}

	// 映射 Go 函数
	vm.Set(funcName, goFunc)

	// 执行 JavaScript 代码
	_, err = vm.RunString(script)
	if err != nil {
		return err
	}

	return nil
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
