package model

import "gorm.io/gorm"

type BaseContest struct {
	gorm.Model
	GameID uint
	States ContestStates `gorm:"embedded"`
	Script string
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

// BeforeSave truncates the message to 1000 characters.
func (ts *TaskStatus) BeforeSave(tx *gorm.DB) (err error) {
	const maxLen = 1000
	if len(ts.Msg) > maxLen {
		ts.Msg = ts.Msg[:maxLen]
	}
	return
}

type DockerTask struct {
	Dockerfile string
	Status     TaskStatus `gorm:"embedded"`
}

// Association CRUD

// contestant

// Sorted by points in descending order.
// Currently supported preloads: "User", "AssignedAi"
func (bc *BaseContest) GetContestants(preloads []preloadQuery) ([]Contestant, error) {
	return GetContestants(map[string]interface{}{"base_contest_id": bc.ID}, preloads)
}

func (bc *BaseContest) GetContestantByUserID(userID uint, preloads []preloadQuery) (Contestant, error) {
	return GetContestant(map[string]interface{}{"base_contest_id": bc.ID, "user_id": userID}, preloads)
}

func (bc *BaseContest) UpdateContestantByUserID(userID uint, updates map[string]interface{}) error {
	return db.Model(&Contestant{}).Where("base_contest_id = ? AND user_id = ?", bc.ID, userID).Updates(updates).Error
}

func (bc *BaseContest) DeleteContestantByUserID(userID uint) error {
	return db.Where("base_contest_id = ? AND user_id = ?", bc.ID, userID).Delete(&Contestant{}).Error
}

// ai

func (bc *BaseContest) GetAis(query QueryParams, preload bool) ([]Ai, int64, error) {
	query.Filter["base_contest_id"] = bc.ID
	return GetAis(query, preload)
}
