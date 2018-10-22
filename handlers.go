package goeureka

import (
	"net/http"

	"github.com/labstack/echo"
)

func Info(c echo.Context) error {
	status := map[string]string{"status": "OK"}
	status["status"] = "OK"
	return c.JSON(http.StatusOK, status)
}

func Health(c echo.Context) error {
	health := map[string]string{"health": "OK"}
	return c.JSON(http.StatusOK, health)
}
