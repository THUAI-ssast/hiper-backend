package task

import (
	"fmt"
	"log"

	"github.com/THUAI-ssast/hiper-backend/web/model"
	"github.com/THUAI-ssast/hiper-backend/worker/repository"
)

func getBinds(domain repository.DomainType, id uint) []string {
	switch domain {
	case repository.GameLogicDomain:
		return []string{
			fmt.Sprintf("/var/hiper/games/%d/game_logic:/app", id),
		}
	case repository.AiDomain:
		ai, err := model.GetAiByID(id, false)
		if err != nil {
			log.Fatal(err)
		}
		sdkID := ai.SdkID
		return []string{
			fmt.Sprintf("/var/hiper/ais/%d:/app", id),
			fmt.Sprintf("/var/hiper/sdks/%d:/sdk", sdkID),
		}
	}
	panic("unreachable")
}
