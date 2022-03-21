package pomodoro

import "time"

type AggregatedPomodoro struct {
	Data     []Time    `json:"data"`
	BaseTime time.Time `json:"base_time"`
}

type Time struct {
	Time      uint64  `json:"time"`
	ProjectId *uint64 `json:"project_id,omitempty"`
}

type GetAggregatedQuery struct {
	Start                *time.Time `query:"start" validate:"required"`
	End                  *time.Time `query:"end" validate:"required"`
	ProjectId            *uint64    `query:"project_id" validate:"omitempty"`
	IncludeSubProjects   bool       `query:"include_sub_project" validate:"omitempty"`
	AggregateSubProjects bool       `query:"aggregate_sub_project" validate:"omitempty"`
}

func GetAggregated(userId uint64, q GetAggregatedQuery) (a AggregatedPomodoro, err error) {
	startTmp := q.Start.UTC()
	q.Start = &startTmp
	endTmp := q.End.UTC()
	q.End = &endTmp

	pomodoros, err := GetList(userId, GetListQuery{Start: q.Start, End: q.End, ProjectId: q.ProjectId, IncludeSubProjects: q.IncludeSubProjects})
	if err != nil {
		return
	}

	var projects map[uint64]uint64 = map[uint64]uint64{}
	var othersTime uint64 = 0

	for _, pomo := range pomodoros {
		if pomo.End == nil {
			continue
		}

		// Trim overdue
		if pomo.Start.Before(*q.Start) {
			pomo.Start = *q.Start
		}
		if pomo.End.After(*q.End) {
			pomo.End = q.End
		}

		if pomo.ProjectId != nil {
			if q.AggregateSubProjects && pomo.ParentProjectId != nil {
				projects[*pomo.ParentProjectId] += uint64(pomo.End.Sub(pomo.Start).Seconds())
			} else {
				projects[*pomo.ProjectId] += uint64(pomo.End.Sub(pomo.Start).Seconds())
			}
		} else {
			othersTime += uint64(pomo.End.Sub(pomo.Start).Seconds())
		}
	}

	var timeAry = []Time{}

	for k, v := range projects {
		timeAry = append(timeAry, Time{v, &k})
	}
	if othersTime != 0 {
		timeAry = append(timeAry, Time{Time: othersTime})
	}

	a.BaseTime = *q.Start
	a.Data = timeAry
	return
}
