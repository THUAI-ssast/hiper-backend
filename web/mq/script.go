package mq

import (
	"github.com/dop251/goja"

	"github.com/THUAI-ssast/hiper-backend/web/model"
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
	// TODO: 与文档不符，注意确认内部字段；修改文档或修改代码
	_, err = CallJSFunction(vm, "onMatchFinished", match.Ais, match.Tag, replay)
	if err != nil {
		return err
	}

	return nil
}

func SetCreateMatch(baseContestID uint) error {
	err := SetGoFuncForJS(baseContestID, "createMatch", func(call goja.FunctionCall) goja.Value {
		// 获取 contestants 参数
		contestantsVal := call.Argument(0)
		contestants, ok := contestantsVal.Export().([]interface{})
		if !ok {
			panic("contestants must be an array")
		}

		// 获取 options 参数
		optionsVal := call.Argument(1)
		options, ok := optionsVal.Export().(map[string]interface{})
		if !ok {
			panic("options must be an object")
		}

		// 调用 createMatch 函数
		err := createMatch(contestants, options, baseContestID)
		if err != nil {
			panic(err)
		}

		// 返回 undefined
		return goja.Undefined()
	})
	if err != nil {
		return err
	}

	return nil
}

func SetGetContestantsByRanking(baseContestID uint) error {
	err := SetGoFuncForJS(baseContestID, "getContestantsByRanking", func(call goja.FunctionCall) goja.Value {
		// 获取 filter 参数
		filterVal := call.Argument(0)
		filter, ok := filterVal.Export().(string)
		if !ok {
			panic("filter must be a string")
		}

		// 调用 getContestantsByRanking 函数
		contestants, err := getContestantsByRanking(filter, baseContestID)
		if err != nil {
			panic(err)
		}

		vm := baseContestIDToRuntime[baseContestID]

		// 将 contestants 转换为 goja.Value
		contestantsVal := vm.ToValue(contestants)

		// 返回 contestantsVal
		return contestantsVal
	})
	if err != nil {
		return err
	}

	return nil
}

func SetUpdateContestant(baseContestID uint) error {
	err := SetGoFuncForJS(baseContestID, "updateContestant", func(call goja.FunctionCall) goja.Value {
		// 获取 contestant 参数
		contestantVal := call.Argument(0)
		contestant, ok := contestantVal.Export().(map[string]interface{})
		if !ok {
			panic("contestant must be an object")
		}

		// 获取 body 参数
		bodyVal := call.Argument(1)
		body, ok := bodyVal.Export().(map[string]interface{})
		if !ok {
			panic("body must be an object")
		}

		// 调用 updateContestant 函数
		err := updateContestant(contestant, body, baseContestID)
		if err != nil {
			panic(err)
		}

		// 返回 undefined
		return goja.Undefined()
	})
	if err != nil {
		return err
	}

	return nil
}
