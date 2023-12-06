package model

import "gorm.io/gorm"

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

type TaskState string

const (
	TaskStatePending     TaskState = "pending"
	TaskStateRunning     TaskState = "running"
	TaskStateFinished    TaskState = "finished"
	TaskStateInputError  TaskState = "input_error"
	TaskStateStopped     TaskState = "stopped"
	TaskStateSystemError TaskState = "system_error"
)

func (ts *TaskState) BeforeCreate(tx *gorm.DB) (err error) {
	*ts = TaskStatePending
	return
}

type TaskStatus struct {
	State TaskState
	Msg   string
}

func (ts *TaskStatus) BeforeSave(tx *gorm.DB) (err error) {
	const maxLen = 1000
	if len(ts.Msg) > maxLen {
		ts.Msg = ts.Msg[:maxLen]
	}
	return
}
