package mq

import (
	"fmt"
	"hiper-backend/model"
	"sync"

	"github.com/dop251/goja"
)

// 定义全局 map
var (
	baseContestIDToRuntime = make(map[uint]*goja.Runtime)
	runtimeToBaseContestID = make(map[*goja.Runtime]uint)
	mutex                  = &sync.Mutex{} // 用于保护 map 的并发访问
)

func CreateRuntimeWithJSFile(baseContestID uint) error {
	// 打开 JavaScript 文件
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		return err
	}
	script := baseContest.Script

	// 创建一个新的 goja 运行时
	vm := goja.New()

	// 执行 JavaScript 代码
	_, err = vm.RunString(script)
	if err != nil {
		return err
	}

	// 添加到全局 map 中
	mutex.Lock()
	baseContestIDToRuntime[baseContestID] = vm
	runtimeToBaseContestID[vm] = baseContestID
	mutex.Unlock()

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

// TODO: 在aiassign api中调用
func CallOnAIAssigned(contestant model.Contestant) error {
	baseContestID := contestant.BaseContestID

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

	// 调用 onAiAssigned 函数
	_, err := CallJSFunction(vm, "onAiAssigned", contestant)
	if err != nil {
		return err
	}

	return nil
}

func CallOnMatchFinished(matchID uint, replay string) error {
	// 根据 matchID 获取 match 对象
	match, err := model.GetMatchByID(matchID, true)
	if err != nil {
		return err
	}

	// 从 match 对象中获取 baseContestID
	baseContestID := match.BaseContestID

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

	// 调用 onMatchFinished 函数
	_, err = CallJSFunction(vm, "onMatchFinished", match.Players, match.Tag, replay)
	if err != nil {
		return err
	}

	return nil
}
