package handler

import (
	"flow-pomodoro/flags"
	"flow-pomodoro/jwt"
	"flow-pomodoro/pomodoro"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func PostEnd(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	post := new(pomodoro.PostEnd)
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

	p, notStarted, invalidTime, err := pomodoro.End(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notStarted {
		// 400: Bad request
		c.Logger().Debug("pomodoro not started")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "pomodoro not started"}, "	")
	}
	if invalidTime {
		// 400: Bad request
		c.Logger().Debug("end time must not before time pomodoro started")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "end time must not before time pomodoro started"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
