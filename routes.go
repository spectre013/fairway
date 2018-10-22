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
		"/info",
		Info,
	},
	Route{
		"Health",
		"POST",
		"/health",
		Health,
	},
}

func buildRoutes(routes Routes, e *echo.Echo) *echo.Echo {
	for _, route := range routes {
		e.Add(route.Method, route.Pattern, route.HandlerFunc)
	}
	return e
}
