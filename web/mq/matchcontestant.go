package mq

import (
	"hiper-backend/model"
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
	return contestants, err
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
