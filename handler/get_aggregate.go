package handler

import (
	"flow-pomodoro/flags"
	"flow-pomodoro/jwt"
	"flow-pomodoro/pomodoro"
	"fmt"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetAggregatedQueryParam struct {
	Start                *string `query:"start" validate:"required,datetime"`
	End                  *string `query:"end" validate:"required,datetime"`
	ProjectId            *uint64 `query:"project_id" validate:"omitempty"`
	IncludeSubProjects   bool    `query:"include_sub_project" validate:"omitempty"`
	AggregationRange     *string `query:"aggregation_range" validate:"omitempty,oneof=hour day week month year"`
	AggregateSubProjects bool    `query:"aggregate_sub_project" validate:"omitempty"`
}

type AggregatedByProjectAndRangePomodoro struct {
	Times     []uint64 `json:"times"`
	ProjectId *uint64  `json:"project_id,omitempty"`
}

func GetAggregated(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind query
	q := new(GetAggregatedQueryParam)
	if err = c.Bind(q); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate query
	if err = c.Validate(q); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}
	var qStart, qEnd *time.Time
	if q.Start != nil {
		startTmp, err := datetimeStrConv(*q.Start)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		qStart = &startTmp
	}
	if q.End != nil {
		endTmp, err := datetimeStrConv(*q.End)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		qEnd = &endTmp
	}
	queryParsed := pomodoro.GetAggregatedQuery{Start: qStart, End: qEnd, ProjectId: q.ProjectId, IncludeSubProjects: q.IncludeSubProjects, AggregateSubProjects: q.AggregateSubProjects}

	// Not multiple aggregation
	if q.AggregationRange == nil {
		// Get aggregatedPomodoro
		aggregatedPomodoro, err := pomodoro.GetAggregated(userId, queryParsed)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		return c.JSONPretty(http.StatusOK, aggregatedPomodoro.Data, "	")
	}

	// Multiple aggregation
	var (
		aggregatedPomodoros []AggregatedByProjectAndRangePomodoro
		projects            map[uint64][]uint64 = map[uint64][]uint64{}
		othersTimes         []uint64
		appearedProjectIds  map[uint64]bool = map[uint64]bool{}
	)

	var rangeInt int
	switch *q.AggregationRange {
	case "hour":
		rangeInt = rangeHour
	case "day":
		rangeInt = rangeDay
	case "week":
		rangeInt = rangeWeek
	case "month":
		rangeInt = rangeMonth
	case "year":
		rangeInt = rangeYear
	}

	var end time.Time
	for start := qStart.UTC(); start.Before(qEnd.UTC()); start = end {
		end, _ = timeRangeEnd(start, rangeInt)
		if end.After(qEnd.UTC()) {
			// Trim overdue
			end = qEnd.UTC()
		}
		queryParsed.Start = &start
		queryParsed.End = &end
		aggregatedPomodoro, err := pomodoro.GetAggregated(userId, queryParsed)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		projectIdNotAppearInColumn := appearedProjectIds
		for _, t := range aggregatedPomodoro.Data {
			if t.ProjectId != nil {
				if !appearedProjectIds[*t.ProjectId] {
					// fill by zero
					for j := 0; j < len(othersTimes); j++ {
						projects[*t.ProjectId] = append(projects[*t.ProjectId], 0)
					}
					// project_id appeared
					appearedProjectIds[*t.ProjectId] = true
				}

				// append time
				projects[*t.ProjectId] = append(projects[*t.ProjectId], t.Time)

				// project_id appeared in column
				projectIdNotAppearInColumn[*t.ProjectId] = false

				if len(othersTimes) != 0 {
					// fill by zero to others
					othersTimes = append(othersTimes, 0)
				}
			} else {
				if len(othersTimes) == 0 {
					// fill by zero
					for j := 0; j < len(othersTimes); j++ {
						othersTimes = append(othersTimes, 0)
					}
				}
				// append time
				othersTimes = append(othersTimes, t.Time)
			}
		}
		if len(aggregatedPomodoro.Data) == 0 {
			othersTimes = append(othersTimes, 0)
		}
		for k, v := range projectIdNotAppearInColumn {
			if v {
				// fill by zero to not appeared project_id
				projects[k] = append(projects[k], 0)
			}
		}
	}

	for i, v := range projects {
		aggregatedPomodoros = append(
			aggregatedPomodoros,
			AggregatedByProjectAndRangePomodoro{
				Times:     v,
				ProjectId: &i,
			},
		)
	}
	if len(othersTimes) != 0 {
		aggregatedPomodoros = append(
			aggregatedPomodoros,
			AggregatedByProjectAndRangePomodoro{
				Times: othersTimes,
			},
		)
	}

	return c.JSONPretty(http.StatusOK, aggregatedPomodoros, "	")
}

const (
	rangeHour = iota + 1
	rangeDay
	rangeWeek
	rangeMonth
	rangeYear
)

func timeRangeEnd(start time.Time, r int) (t time.Time, err error) {
	switch r {
	case rangeHour:
		t = start.Add(time.Hour)
	case rangeDay:
		t = start.AddDate(0, 0, 1)
	case rangeWeek:
		t = start.AddDate(0, 0, 7)
	case rangeMonth:
		t = start.AddDate(0, 1, 0)
	case rangeYear:
		t = start.AddDate(1, 0, 0)
	default:
		err = fmt.Errorf("range %d not defined", r)
	}
	return
}
