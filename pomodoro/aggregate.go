package pomodoro

import "time"

type AggregatedPomodoro struct {
	Data     []Time    `json:"data"`
	BaseTime time.Time `json:"base_time"`
}

type Time struct {
	Time      uint    `json:"time"`
	ProjectId *uint64 `json:"project_id,omitempty"`
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
	var othersTime uint = 0

	for _, pomo := range pomodoros {
		if pomo.End == nil {
			continue
		}

		// Trim overdue
		if pomo.Start.Before(start) {
			pomo.Start = start
		}
		if pomo.End.After(end) {
			pomo.End = &end
		}

		if pomo.ProjectId != nil {
			projects[*pomo.ProjectId] += uint(pomo.End.Sub(pomo.Start).Seconds())
		} else {
			othersTime += uint(pomo.End.Sub(pomo.Start).Seconds())
		}
	}

	var timeAry = []Time{}

	for k, v := range projects {
		timeAry = append(timeAry, Time{v, &k})
	}
	if othersTime != 0 {
		timeAry = append(timeAry, Time{Time: othersTime})
	}

	a.BaseTime = start
	a.Data = timeAry
	return
}
