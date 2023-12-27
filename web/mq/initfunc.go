package mq

import (
	"context"
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
)

var Ctx_callback = context.Background()

// 定义全局 map
var (
	BaseContestIDToRuntime = make(map[uint]*eventloop.EventLoop)
	RuntimeToBaseContestID = make(map[*eventloop.EventLoop]uint)
	Mutex                  = &sync.Mutex{} // 用于保护 map 的并发访问
)

func InitMq() {
	go ListenMsgForMatchFinished(Ctx_callback, "match_result")
}

func InitGameMq(baseContestID uint, vm *goja.Runtime) {
	SetCreateMatch(baseContestID, vm)
	SetGetContestantsByRanking(baseContestID, vm)
	SetUpdateContestant(baseContestID, vm)
	SendBuildGameLogicMsg(Ctx_callback, baseContestID)
}

func SetGoFuncForJS(baseContestID uint, funcName string, goFunc func(goja.FunctionCall) goja.Value, vm *goja.Runtime) error {
	vm.Set(funcName, goFunc)
	return nil
}

func CallJSFunction(loop *eventloop.EventLoop, funcName string, args ...interface{}) (result goja.Value, err error) {
	loop.Run(func(vm *goja.Runtime) {
		// 获取 JavaScript 函数
		jsFunc, ok := goja.AssertFunction(vm.Get(funcName))
		if !ok {
			err = fmt.Errorf("function %s not found in JavaScript code", funcName)
			return
		}

		// 准备函数参数
		funcArgs := make([]goja.Value, len(args))
		for i, arg := range args {
			funcArgs[i] = vm.ToValue(arg)
		}

		// 调用 JavaScript 函数
		result, err = jsFunc(goja.Undefined(), funcArgs...)
	})

	// 返回结果
	return result, err
}
