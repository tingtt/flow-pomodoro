package main

import (
	"flow-pomodoro/jwt"
	"flow-pomodoro/pomodoro"
	"fmt"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetAggregatedQueryParam struct {
	AggregationRange   *string `query:"aggregation_range" validate:"omitempty,oneof=hour day week month year"`
	ProjectId          *uint64 `query:"project_id" validate:"omitempty"`
	IncludeSubProjects bool    `query:"include_sub_project" validate:"omitempty"`
	Start              string  `query:"start" validate:"required,datetime"`
	End                string  `query:"end" validate:"required,datetime"`
}

func getAggregated(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
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
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}
	qStart, _ := datetimeStrConv(q.Start)
	qEnd, _ := datetimeStrConv(q.End)
	if qStart.After(qEnd) {
		// 422: Unprocessable entity
		c.Logger().Debug("`start` must before `end`")
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": "`start` must before `end`"}, "	")
	}

	// Not multiple aggregation
	if q.AggregationRange == nil {
		// Get aggregatedPomodoro
		aggregatedPomodoro, err := pomodoro.GetAggregated(userId, qStart, qEnd, q.ProjectId, q.IncludeSubProjects)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		return c.JSONPretty(http.StatusOK, aggregatedPomodoro, "	")
	}

	// Multiple aggregation
	var aggregatedPomodoros []pomodoro.AggregatedPomodoro

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
	for start := qStart; start.Before(qEnd); start = end {
		end, _ = timeRangeEnd(start, rangeInt)
		if end.After(qEnd) {
			// Trim overdue
			end = qEnd
		}
		aggregatedPomodoro, err := pomodoro.GetAggregated(userId, start, end, q.ProjectId, q.IncludeSubProjects)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		aggregatedPomodoros = append(aggregatedPomodoros, aggregatedPomodoro)
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
