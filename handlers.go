package goeureka

import (
	"net/http"

	"github.com/labstack/echo"
)


func Info(c echo.Context) error {
	return c.JSON(http.StatusOK, info())
}

func Health(c echo.Context) error {
	health := map[string]string{"status": "UP"}
	return c.JSON(http.StatusOK, health)
}

func Env(c echo.Context) error {
	return c.JSON(http.StatusOK, env())
}

func Metrics(c echo.Context) error {
	return c.JSON(http.StatusOK, metrics())
}
