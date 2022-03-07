package pomodoro

import "time"

type AggregatedPomodoro struct {
	Data     []Time    `json:"data"`
	BaseTime time.Time `json:"base_time"`
}

type Time struct {
	Time      uint   `json:"time"`
	ProjectId uint64 `json:"project_id"`
}
