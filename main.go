package main

import (
	"flag"
	"flow-pomodoros/jwt"
	"flow-pomodoros/mysql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		var intValue, err = strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Priority: command line params > env variables > default value
var (
	port        = flag.Int("port", getIntEnv("PORT", 1323), "Server port")
	logLevel    = flag.Int("log-level", getIntEnv("LOG_LEVEL", 2), "Log level (1: 'DEBUG', 2: 'INFO', 3: 'WARN', 4: 'ERROR', 5: 'OFF', 6: 'PANIC', 7: 'FATAL'")
	gzipLevel   = flag.Int("gzip-level", getIntEnv("GZIP_LEVEL", 6), "Gzip compression level")
	mysqlHost   = flag.String("mysql-host", getEnv("MYSQL_HOST", "db"), "MySQL host")
	mysqlPort   = flag.Int("mysql-port", getIntEnv("MYSQL_PORT", 3306), "MySQL port")
	mysqlDB     = flag.String("mysql-database", getEnv("MYSQL_DATABASE", "flow-pomodoros"), "MySQL database")
	mysqlUser   = flag.String("mysql-user", getEnv("MYSQL_USER", "flow-pomodoros"), "MySQL user")
	mysqlPasswd = flag.String("mysql-password", getEnv("MYSQL_PASSWORD", ""), "MySQL password")
	jwtIssuer   = flag.String("jwt-issuer", getEnv("JWT_ISSUER", "flow-users"), "JWT issuer")
	jwtSecret   = flag.String("jwt-secret", getEnv("JWT_SECRET", ""), "JWT secret")
)

type CustomValidator struct {
	validator *validator.Validate
}

func DatetimeStrValidation(fl validator.FieldLevel) bool {
	_, err1 := time.Parse("2006-1-2T15:4:5", fl.Field().String())
	_, err2 := strconv.ParseUint(fl.Field().String(), 10, 64)
	fmt.Printf("err1: %v\n", err1)
	fmt.Printf("err2: %v\n", err2)
	return err1 == nil || err2 == nil
}

func datetimeStrConv(str string) (t time.Time, err error) {
	// y-m-dTh:m:s or unix timestamp
	t, err1 := time.Parse("2006-1-2T15:4:5", str)
	u, err2 := strconv.ParseInt(str, 10, 64)
	if err1 == nil {
		return
	}
	if err2 == nil {
		t = time.Unix(u, 0)
		return
	}
	err = fmt.Errorf("\"%s\" is not a unix timestamp or string format \"2006-1-2T15:4:5\"", str)
	return
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
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: *gzipLevel,
	}))
	e.Logger.SetLevel(log.Lvl(*logLevel))
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.JwtCustumClaims{},
		SigningKey: []byte(*jwtSecret),
	}))

	// Setup db client instance
	e.Logger.Info(mysql.SetDSNTCP(*mysqlUser, *mysqlPasswd, *mysqlHost, *mysqlPort, *mysqlDB))

	// Restricted routes
	e.GET("/", get)
	e.GET("/aggregated", getAggregated)
	e.POST("/start", postStart)
	e.POST("/end", postEnd)
	e.DELETE(":id", delete)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
