package main

import (
	"flow-pomodoros/jwt"
	"flow-pomodoros/pomodoro"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func postStart(c echo.Context) error {
	// Check `Content-Type`
	if c.Request().Header.Get("Content-Type") != "application/json" &&
		c.Request().Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	post := new(pomodoro.PostStart)
	if err = c.Bind(post); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(post); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// TODO: Check todo id

	// TODO: Check project id

	p, notEnded, invalidTime, err := pomodoro.Start(userId, *post, false)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notEnded {
		// 409: Conflict
		c.Logger().Debug("pomodoro not ended")
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "pomodoro not ended"}, "	")
	}
	if invalidTime {
		// 409: Conflict
		c.Logger().Debug("start time must not before time last pomodoro ended")
		return c.JSONPretty(http.StatusConflict, map[string]string{"message": "start time must not before time last pomodoro ended"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
