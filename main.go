package main

import (
	"flow-pomodoro/flags"
	"flow-pomodoro/handler"
	"flow-pomodoro/jwt"
	"flow-pomodoro/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func DatetimeStrValidation(fl validator.FieldLevel) bool {
	_, err1 := time.Parse("2006-1-2T15:4:5", fl.Field().String())
	_, err2 := time.Parse(time.RFC3339, fl.Field().String())
	_, err3 := strconv.ParseUint(fl.Field().String(), 10, 64)
	return err1 == nil || err2 == nil || err3 == nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	cv.validator.RegisterValidation("datetime", DatetimeStrValidation)

	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func main() {
	// Get command line params / env variables
	f := flags.Get()

	//
	// Setup echo and middlewares
	//

	// Echo instance
	e := echo.New()

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: int(*f.GzipLevel),
	}))

	// Log level
	e.Logger.SetLevel(log.Lvl(*f.LogLevel))

	// Validator instance
	e.Validator = &CustomValidator{validator: validator.New()}

	// JWT
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*f.JwtSecret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/-/readiness"
		},
	}))

	//
	// Check health of external service
	//

	// flow-projects
	if *flags.Get().ServiceUrlProjects == "" {
		e.Logger.Fatal("`--service-url-projects` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlProjects+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-projects` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-projects`")
	}
	// flow-todos
	if *flags.Get().ServiceUrlTodos == "" {
		e.Logger.Fatal("`--service-url-todos` option is required")
	}
	if status, err := utils.HttpGet(*flags.Get().ServiceUrlTodos+"/-/readiness", nil); err != nil {
		e.Logger.Fatalf("failed to check health of external service `flow-todos` %s", err)
	} else if status != http.StatusOK {
		e.Logger.Fatal("failed to check health of external service `flow-todos`")
	}

	//
	// Routes
	//

	// Health check route
	e.GET("/-/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, "flow-pomodoro is Healthy.\n")
	})

	// Restricted routes
	e.GET("/", handler.GetList)
	e.GET("/aggregated", handler.GetAggregated)
	e.POST("/start", handler.PostStart)
	e.POST("/end", handler.PostEnd)
	e.DELETE(":id", handler.Delete)
	e.DELETE("/", handler.DeleteAll)

	//
	// Start echo
	//
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *f.Port)))
}
