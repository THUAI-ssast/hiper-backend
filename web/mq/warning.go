package mq

import (
	"fmt"
	"math/rand"

	"github.com/THUAI-ssast/hiper-backend/web/model"

	"github.com/dop251/goja"
)

func WarnCode() {
	// 设置 baseContestID
	baseContestID := uint(1)

	// 获取 baseContest 的 script
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		fmt.Println("Failed to get baseContest")
		return
	}

	script := baseContest.Script

	fmt.Println(script)

	// 查找 map 中对应的 runtime
	Mutex.Lock()
	loop, exists := BaseContestIDToRuntime[baseContestID]
	Mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			fmt.Println("Failed to create runtime")
			return
		}

		// 再次查找 runtime
		Mutex.Lock()
		loop = BaseContestIDToRuntime[baseContestID]
		Mutex.Unlock()
	}

	// 使用 vm 运行 baseContest.script
	loop.Run(func(vm *goja.Runtime) {
		//InitGameMq(baseContestID, vm)
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

		GetFunc := func(call goja.FunctionCall) goja.Value {
			arr := []int{1, 2}
			fmt.Println("abcde")

			// 将数组转换为 goja.Value
			return vm.ToValue(arr)
		}
		vm.Set("gett", GetFunc)
		SetFunc := func(call goja.FunctionCall) goja.Value {
			match := model.Match{BaseContestID: baseContestID}
			//TODO:DELETE!
			// 生成随机的 score
			score := []int{rand.Intn(2), rand.Intn(2)}
			// 确保 score 是 [0, 1] 或 [1, 0]
			if score[0] == score[1] {
				score[1] = 1 - score[0]
			}
			match.Scores = score
			ai1, _ := model.GetAiByID(uint(1), false)
			ai2, _ := model.GetAiByID(uint(2), false)
			c1, _ := model.GetContestant(map[string]interface{}{"base_contest_id": baseContestID, "user_id": ai1.UserID}, nil)
			c2, _ := model.GetContestant(map[string]interface{}{"base_contest_id": baseContestID, "user_id": ai2.UserID}, nil)
			model.UpdateContestantByID(c1.ID, map[string]interface{}{"points": c1.Points + score[0]})
			model.UpdateContestantByID(c2.ID, map[string]interface{}{"points": c2.Points + score[1]})
			match.Create([]uint{uint(1), uint(2)})
			return goja.Undefined()
		}
		vm.Set("createe", SetFunc)
		_, err = vm.RunString(script)
	})
	if err != nil {
		fmt.Println("Failed to run script")
		return
	}
}
