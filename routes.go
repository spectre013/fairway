package goeureka

import "github.com/labstack/echo"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc echo.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Info",
		"GET",
		"/actuator/info",
		Info,
	},
	Route{
		"Health",
		"POST",
		"/actuator/health",
		Health,
	},
	Route{
		"Env",
		"GET",
		"/actuator/env",
		Env,
	},
	Route{
		"Metrics",
		"GET",
		"/actuator/metrics",
		Metrics,
	},
}

func BuildRoutes(routes Routes, e *echo.Echo) *echo.Echo {
	loadGitInfo()
	for _, route := range routes {
		e.Add(route.Method, route.Pattern, route.HandlerFunc)
	}
	return e
}
