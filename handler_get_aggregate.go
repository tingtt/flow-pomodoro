package main

import (
	"flow-pomodoros/jwt"
	"flow-pomodoros/pomodoro"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetAggregatedQueryParam struct {
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

	// Get aggregatedPomodoro
	aggregatedPomodoro, err := pomodoro.GetAggregated(userId, q.Start, q.End, q.ProjectId, q.IncludeSubProjects)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	return c.JSONPretty(http.StatusOK, aggregatedPomodoro, "	")
}
