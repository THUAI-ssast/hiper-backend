package mq

import (
	"hiper-backend/model"
)

func AddMatch(userIDs []uint, baseContestID uint, tag string) (matchID uint, err error) {
	match := model.Match{BaseContestID: baseContestID, Tag: tag}
	err = match.Create(userIDs)
	if err != nil {
		return 0, err
	}
	SendBuildMatchMsg(model.Ctx, match.ID)
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
