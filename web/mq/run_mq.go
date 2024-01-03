package mq

import (
	"fmt"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/dop251/goja"
)

func InitGameMqAndRunScript(baseContestID uint) {
	// 获取 baseContest 的 script
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		fmt.Println("Failed to get baseContest by ID")
	}

	script := baseContest.Script

	// 查找 map 中对应的 runtime
	Mutex.Lock()
	loop, exists := BaseContestIDToRuntime[baseContestID]
	Mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			fmt.Println("cannot create loop")
		}

		// 再次查找 runtime
		Mutex.Lock()
		loop = BaseContestIDToRuntime[baseContestID]
		Mutex.Unlock()
	}

	// 使用 vm 运行 baseContest.script
	loop.Run(func(vm *goja.Runtime) {
		InitGameMq(baseContestID, vm)
		// 定义一个 Go 函数
		printFunc := func(call goja.FunctionCall) goja.Value {
			// 获取函数参数
			message := call.Argument(0).String()

			// 输出语句
			fmt.Println(message)

			// 返回 undefined
			return vm.ToValue(nil)
		}

		// 将 Go 函数映射到 JavaScript 环境
		vm.Set("print", printFunc)
		_, err = vm.RunString(script)
	})
	if err != nil {
		fmt.Println("Failed to run script")
	}
}
