package repository

type DomainType string
type OperationType string

const (
	GameLogicDomain DomainType = "game_logic"
	AiDomain        DomainType = "ai"
)

const (
	BuildOperation OperationType = "build"
	RunOperation   OperationType = "run"
)
