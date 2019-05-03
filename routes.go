package fairway

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//Route data type
type Route struct {
	Name     string
	Method   string
	Pattern  string
	Produces string
	Handler  http.Handler
}

//Routes to be combined with main app routes to set up API
type Routes []Route

var routes = Routes{
	Route{
		"Info",
		"GET",
		"/actuator/info",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Info),
	},
	Route{
		"Health",
		"GET",
		"/actuator/health",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Health),
	},
	Route{
		"Env",
		"GET",
		"/actuator/env",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Env),
	},
	Route{
		"Env",
		"GET",
		"/actuator/env/{toMatch}",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Env),
	},
	Route{
		"Metrics",
		"GET",
		"/actuator/metrics",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Metrics),
	},
	Route{
		"Metrics Property",
		"GET",
		"/actuator/metrics/{requiredMetricName}",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Metrics),
	},
	Route{
		"Actuator",
		"GET",
		"/actuator",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Actuator),
	},
	Route{
		"Loggers",
		"GET",
		"/actuator/loggers",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Loggers),
	},
	Route{
		"Loggers",
		"GET",
		"/actuator/loggers/{name}",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Loggers),
	},
	Route{
		"Loggers",
		"POST",
		"/actuator/loggers/{name}",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(UpdateLogger),
	},
	Route{
		"Mappings",
		"GET",
		"/actuator/mappings",
		"application/vnd.spring-boot.actuator.v2+json;charset=UTF-8",
		http.HandlerFunc(Mappings),
	},
}

// BuildRoutes - Builds route handlers
func BuildRoutes(routes Routes, e *mux.Router) *mux.Router {
	loadGitInfo()
	for _, route := range routes {
		if secure.Enable && strings.HasPrefix(route.Pattern, "/actuator") {
			route.Handler = basicAuth(route.Handler, secure.User, secure.Password, "Password required to access actuator endpoints")
		}
		e.Handle(route.Pattern, route.Handler).Methods(route.Method)

	}
	e.Use(loggingMiddleware)
	return e
}

func basicAuth(handler http.Handler, username, password, realm string) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%v", r)
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("unauthorized.\n"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
