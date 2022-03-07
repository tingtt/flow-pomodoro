package main

import (
	"flow-pomodoros/jwt"
	"flow-pomodoros/pomodoro"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetQueryParam struct {
	ProjectId          *uint64   `query:"project_id" validate:"omitempty"`
	IncludeSubProjects bool      `query:"include_sub_project" validate:"omitempty"`
	TodoId             *uint64   `query:"todo_id" validate:"omitempty"`
	Start              time.Time `query:"start" validate:"required"`
	End                time.Time `query:"end" validate:"required"`
}

func get(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Bind query
	q := new(GetQueryParam)
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
	if q.ProjectId != nil && q.TodoId != nil {
		// 422: Unprocessable entity
		c.Logger().Debug("`project_id` and `todo_id` cannnot query at the same time")
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": "`project_id` and `todo_id` cannnot query at the same time"}, "	")
	}
	if q.Start.After(q.End) {
		// 422: Unprocessable entity
		c.Logger().Debug("`start` must before `end`")
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": "`start` must before `end`"}, "	")
	}

	// Get pomodoros
	var pomodoros []pomodoro.Pomodoro
	if q.ProjectId != nil {
		pomodoros, err = pomodoro.GetListProjectId(userId, q.Start, q.End, *q.ProjectId, q.IncludeSubProjects)
	} else if q.TodoId != nil {
		pomodoros, err = pomodoro.GetListTodo(userId, q.Start, q.End, *q.TodoId)
	} else {
		pomodoros, err = pomodoro.GetList(userId, q.Start, q.End)
	}
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	if pomodoros == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, pomodoros, "	")
}
