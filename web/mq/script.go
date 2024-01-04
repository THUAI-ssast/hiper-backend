package mq

import (
	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
)

func CreateRuntimeWithJSFile(baseContestID uint) error {

	// 创建一个新的 eventloop
	loop := eventloop.NewEventLoop()

	// 添加到全局 map 中
	Mutex.Lock()
	BaseContestIDToRuntime[baseContestID] = loop
	RuntimeToBaseContestID[loop] = baseContestID
	Mutex.Unlock()

	return nil
}

func CallOnAIAssigned(contestant model.Contestant) error {
	baseContestID := contestant.BaseContestID

	// 查找 map 中对应的 runtime
	Mutex.Lock()
	loop, exists := BaseContestIDToRuntime[baseContestID]
	Mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			return err
		}

		// 再次查找 runtime
		Mutex.Lock()
		loop = BaseContestIDToRuntime[baseContestID]
		Mutex.Unlock()
	}

	// 调用 onAiAssigned 函数
	_, err := CallJSFunction(loop, "onAiAssigned", contestant)
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
	Mutex.Lock()
	loop, exists := BaseContestIDToRuntime[baseContestID]
	Mutex.Unlock()

	// 如果没有找到，创建一个新的 runtime
	if !exists {
		err := CreateRuntimeWithJSFile(baseContestID)
		if err != nil {
			return err
		}

		// 再次查找 runtime
		Mutex.Lock()
		loop = BaseContestIDToRuntime[baseContestID]
		Mutex.Unlock()
	}

	var players []map[string]interface{}

	for i, ai := range match.Ais {
		contestant, _ := model.GetContestant(map[string]interface{}{
			"base_contest_id": baseContestID,
			"assigned_ai_id":  ai.ID,
		}, nil)
		players = append(players, map[string]interface{}{
			"contestant": contestant,
			"score":      match.Scores[i],
		})
	}

	// 调用 onMatchFinished 函数
	_, err = CallJSFunction(loop, "onMatchFinished", players, match.Tag, replay)
	if err != nil {
		return err
	}

	return nil
}

func SetCreateMatch(baseContestID uint, vm *goja.Runtime) error {
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
			options = make(map[string]interface{})
		}

		// 调用 createMatch 函数
		err := createMatch(contestants, options, baseContestID)
		if err != nil {
			panic(err)
		}

		return goja.Undefined()
	}, vm)
	if err != nil {
		return err
	}

	return nil
}

func SetGetContestantsByRanking(baseContestID uint, vm *goja.Runtime) error {
	err := SetGoFuncForJS(baseContestID, "getContestantsByRanking", func(call goja.FunctionCall) goja.Value {
		// 获取 filter 参数
		filterVal := call.Argument(0)
		filter, ok := filterVal.Export().(string)
		if !ok {
			filter = "survived"
		}

		// 调用 getContestantsByRanking 函数
		contestants, err := getContestantsByRanking(filter, baseContestID)
		if err != nil {
			panic(err)
		}

		contestantsjs := make([]interface{}, 0, len(contestants))
		for _, contestant := range contestants {
			contestantsjs = append(contestantsjs, map[string]interface{}{
				"username":           contestant.User.Username,
				"assignedAiId":       contestant.AssignedAiID,
				"points":             contestant.Points,
				"performance":        contestant.Performance,
				"assignAiEnabled":    contestant.Permissions.AssignAiEnabled,
				"publicMatchEnabled": contestant.Permissions.PublicMatchEnabled,
			})
		}

		// 创建一个新的 JavaScript 数组
		contestantsVal := vm.ToValue(contestantsjs)

		// 返回 contestantsVal
		return contestantsVal
	}, vm)
	if err != nil {
		return err
	}

	return nil
}

func SetUpdateContestant(baseContestID uint, vm *goja.Runtime) error {
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
	}, vm)
	if err != nil {
		return err
	}

	return nil
}
