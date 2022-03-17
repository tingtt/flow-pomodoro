package main

import (
	"flow-pomodoro/jwt"
	"flow-pomodoro/pomodoro"
	"net/http"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type GetQueryParam struct {
	Start              *string `query:"start" validate:"required,datetime"`
	End                *string `query:"end" validate:"required,datetime"`
	ProjectId          *uint64 `query:"project_id" validate:"omitempty"`
	IncludeSubProjects bool    `query:"include_sub_project" validate:"omitempty"`
	TodoId             *uint64 `query:"todo_id" validate:"omitempty"`
}

func getList(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
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
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}
	var start, end *time.Time
	if q.Start != nil {
		startTmp, err := datetimeStrConv(*q.Start)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		start = &startTmp
	}
	if q.End != nil {
		endTmp, err := datetimeStrConv(*q.End)
		if err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}
		end = &endTmp
	}
	queryParsed := pomodoro.GetListQuery{Start: start, End: end, ProjectId: q.ProjectId, IncludeSubProjects: q.IncludeSubProjects, TodoId: q.TodoId}

	// Get pomodoros
	pomodoros, err := pomodoro.GetList(userId, queryParsed)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	if pomodoros == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, pomodoros, "	")
}
