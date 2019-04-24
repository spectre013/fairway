package fairway

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Handler     http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"Info",
		"GET",
		"/actuator/info",
		nil,
		http.HandlerFunc(Info),
	},
	Route{
		"Health",
		"POST",
		"/actuator/health",
		nil,
		http.HandlerFunc(Health),
	},
	Route{
		"Env",
		"GET",
		"/actuator/env",
		nil,
		http.HandlerFunc(Env),
	},
	Route{
		"Env",
		"GET",
		"/actuator/env/",
		nil,
		http.HandlerFunc(Env),
	},
	Route{
		"Metrics",
		"GET",
		"/actuator/metrics",
		nil,
		http.HandlerFunc(Metrics),
	},
	Route{
		"Metrics Property",
		"GET",
		"/actuator/metrics/",
		nil,
		http.HandlerFunc(Metrics),
	},
	Route{
		"Actuator",
		"GET",
		"/actuator",
		nil,
		http.HandlerFunc(Actuator),
	},
	Route{
		"Loggers",
		"GET",
		"/actuator/loggers",
		nil,
		http.HandlerFunc(Loggers),
	},
	Route{
		"Loggers",
		"GET",
		"/actuator/loggers/",
		nil,
		http.HandlerFunc(Loggers),
	},
	Route{
		"Error",
		"GET",
		"/actuator/",
		nil,
		http.HandlerFunc(Error),
	},
}

func BuildRoutes(routes Routes, e *http.ServeMux) *http.ServeMux {
	loadGitInfo()
	for _, route := range routes {
		if route.HandlerFunc != nil {
			e.Handle(route.Pattern, Logger(route.HandlerFunc, route.Name))
		} else {
			e.Handle(route.Pattern, Logger(route.Handler, route.Name))
		}
	}
	return e
}
