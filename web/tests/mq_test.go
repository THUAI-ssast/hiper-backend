package tests

import (
	"fmt"
	"testing"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/web/mq"

	"github.com/dop251/goja"
)

func TestInitGameMqAndRunScript(t *testing.T) {
	// 设置 baseContestID
	baseContestID := uint(2)

	// 获取 baseContest 的 script
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		t.Fatalf("Failed to get baseContest by ID: %v", err)
	}

	script := baseContest.Script

	// 查找 map 中对应的 runtime
	mq.Mutex.Lock()
	loop, exists := mq.BaseContestIDToRuntime[baseContestID]
	mq.Mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := mq.CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			t.Fatalf("cannot create loop: %v", err)
		}

		// 再次查找 runtime
		mq.Mutex.Lock()
		loop = mq.BaseContestIDToRuntime[baseContestID]
		mq.Mutex.Unlock()
	}

	// 使用 vm 运行 baseContest.script
	loop.Run(func(vm *goja.Runtime) {
		mq.InitGameMq(baseContestID, vm)
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
		t.Fatalf("Failed to run script: %v", err)
	}
}
