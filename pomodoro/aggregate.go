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

func GetAggregated(userId uint64, start time.Time, end time.Time, projectId *uint64, includeSubProject bool) (a AggregatedPomodoro, err error) {
	var pomodoros []Pomodoro

	if projectId != nil {
		pomodoros, err = GetListProjectId(userId, start, end, *projectId, includeSubProject)
		if err != nil {
			return
		}
	} else {
		pomodoros, err = GetList(userId, start, end)
		if err != nil {
			return
		}
	}

	var projects map[uint64]uint = map[uint64]uint{}

	for _, pomodoro := range pomodoros {
		if pomodoro.End == nil {
			continue
		}

		// Trim overdue
		if pomodoro.Start.Before(start) {
			pomodoro.Start = start
		}
		if pomodoro.End.After(end) {
			pomodoro.End = &end
		}

		projects[pomodoro.ProjectId] += uint(pomodoro.End.Sub(pomodoro.Start).Seconds())
	}

	var timeAry []Time

	for k, v := range projects {
		timeAry = append(timeAry, Time{v, k})
	}

	a.BaseTime = start
	a.Data = timeAry
	return
}
