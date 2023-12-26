package mq

import (
	"hiper-backend/model"
	"sort"
)

func getContestantsByRanking(filter string, baseContestID uint) (contestants []model.Contestant, err error) {
	baseContest, err := model.GetBaseContestByID(baseContestID)
	if err != nil {
		return nil, err
	}
	preloads := []model.PreloadQuery{
		{
			Table:   "User",
			Columns: []string{},
		},
		{
			Table:   "Ai",
			Columns: []string{},
		},
	}
	contestants, err = baseContest.GetContestants(preloads)
	sort.Slice(contestants, func(i, j int) bool {
		return contestants[i].Points > contestants[j].Points
	})
	if filter == "all" {
		return contestants, err
	} else if filter == "survived" {
		survivedContestants := make([]model.Contestant, 0)
		for _, contestant := range contestants {
			if contestant.Permissions.PublicMatchEnabled {
				survivedContestants = append(survivedContestants, contestant)
			}
		}
		return survivedContestants, nil
	} else {
		eliminatedContestants := make([]model.Contestant, 0)
		for _, contestant := range contestants {
			if !contestant.Permissions.PublicMatchEnabled {
				eliminatedContestants = append(eliminatedContestants, contestant)
			}
		}
		return eliminatedContestants, nil
	}
}

func createMatch(contestantsjs []interface{}, options map[string]interface{}, baseContestID uint) (err error) {
	Ais := []uint{}
	for _, contestantjs := range contestantsjs {
		contestantjsm := contestantjs.(map[string]interface{})
		contestantID := uint(contestantjsm["id"].(float64))
		contestant, err := model.GetContestantByID(contestantID, nil)
		if err != nil {
			return err
		}
		Ais = append(Ais, contestant.AssignedAi.ID)
	}
	tag := options["tag"].(string)
	extraInfo := options["extraInfo"].(map[string]interface{})
	_, err = AddMatch(Ais, baseContestID, tag, extraInfo)
	return err
}

func updateContestant(contestantjs interface{}, body map[string]interface{}, baseContestID uint) (err error) {
	contestantjsm := contestantjs.(map[string]interface{})
	contestantID := uint(contestantjsm["id"].(float64))
	contestant, err := model.GetContestantByID(contestantID, nil)
	if err != nil {
		return err
	}

	// 获取 body 中的字段
	performance, ok := body["Performance"]
	if ok {
		contestant.Performance = performance.(string)
	}
	assignAiEnabled, ok := body["AssignAiEnabled"]
	if ok {
		contestant.Permissions.AssignAiEnabled = assignAiEnabled.(bool)
	}
	publicMatchEnabled, ok := body["PublicMatchEnabled"]
	if ok {
		contestant.Permissions.PublicMatchEnabled = publicMatchEnabled.(bool)
	}
	points, ok := body["Points"]
	if ok {
		contestant.Points = points.(int)
	}

	// 更新 contestant
	err = model.UpdateContestantByID(contestantID, map[string]interface{}{
		"performance": contestant.Performance,
		"permissions": contestant.Permissions,
		"points":      contestant.Points,
	})
	if err != nil {
		return err
	}

	return nil
}

func AddMatch(playerIDs []uint, baseContestID uint, tag string, extraInfo map[string]interface{}) (matchID uint, err error) {
	match := model.Match{BaseContestID: baseContestID, Tag: tag}
	err = match.Create(playerIDs)
	if err != nil {
		return 0, err
	}
	SendBuildMatchMsg(model.Ctx, match.ID, extraInfo)
	return match.ID, nil
}

// 在ChangeMatch中调用
func StartMatch(matchID uint) {
	match, err := model.GetMatchByID(matchID, true)
	if err != nil {
		return
	}
	aiids := []uint{}
	for _, ai := range match.Players {
		aiids = append(aiids, ai.ID)
	}
	SendRunMatchMsg(model.Ctx, matchID)
	SendAIIDsMsg(model.Ctx, aiids)
}

func ChangeMatch(tag string, state string, matchID uint) (err error) {
	match, err := model.GetMatchByID(matchID, true)
	if err != nil {
		return err
	}
	if state == "running" && match.State != "running" {
		err = model.UpdateMatchByID(matchID, map[string]interface{}{"tag": tag, "state": state})
		if err != nil {
			return err
		}
		StartMatch(matchID)
		return nil
	}
	err = model.UpdateMatchByID(matchID, map[string]interface{}{"tag": tag, "state": state})
	if err != nil {
		return err
	}
	SendChangeMatchMsg(model.Ctx, matchID)
	return nil
}
