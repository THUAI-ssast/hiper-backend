package model

type BaseContest struct {
	GameId   uint
	Metadata Metadata      `gorm:"embedded"`
	States   ContestStates `gorm:"embedded"`
	Script   string
}

type Metadata struct {
	CoverUrl string
	Readme   string
	Title    string
}

type ContestStates struct {
	AssignAiEnabled                 bool `gorm:"default: false"`
	CommitAiEnabled                 bool `gorm:"default: false"`
	ContestScriptEnvironmentEnabled bool `gorm:"default: false"`
	PrivateMatchEnabled             bool `gorm:"default: false"`
	PublicMatchEnabled              bool `gorm:"default: false"`
	TestMatchEnabled                bool `gorm:"default: false"`
}
