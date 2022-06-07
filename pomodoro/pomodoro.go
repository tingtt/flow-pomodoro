package pomodoro

import (
	"time"
)

type Pomodoro struct {
	Id              uint64     `json:"id"`
	Start           time.Time  `json:"start"`
	End             *time.Time `json:"end,omitempty"`
	RemainingTime   *uint      `json:"remaining_time,omitempty"`
	TodoId          uint64     `json:"todo_id"`
	ProjectId       *uint64    `json:"project_id,omitempty"`
	ParentProjectId *uint64    `json:"parent_project_id,omitempty"`
}
