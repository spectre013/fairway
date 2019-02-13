package goeureka

import (
	"net/http"
	"os"
	"testing"
)
var BuildHttpRequest = buildHttpRequest
func getConfig() EurekaConfig {
	return EurekaConfig{
		Name:        "test",
		Url:         "http://eurekaserver:8761/eureka",
		HostName:    "test",
		Port:        os.Getenv("SEARCH_UI_SERVICE_PORT"),
		SecurePort:  "8943",
		RestService: false,
	}
}

func TestCreateHTTPAction(t *testing.T) {
	config := getConfig()
	reg := CreateRegistration(config)

	httpaction := CreateHTTPAction(config, reg)

	expected := config.Url + "/apps/" + config.Name

	if httpaction.Url != expected {
		t.Errorf("got %s wanted %s", httpaction.Url, expected)
	}
}

func TestCreateRegistration(t *testing.T) {
	config := getConfig()
	reg := CreateRegistration(config)

	if reg.Instance.HostName != config.HostName {
		t.Errorf("got %s wanted %s", reg.Instance.HostName, config.HostName)
	}
}

func TestCombineRoutes(t *testing.T) {
	total := 5
	addRoute := Routes{
		Route{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: nil,
			Handler:     http.FileServer(http.Dir("dist/search")),
		},
	}

	combined := CombineRoutes(addRoute, routes)
	if len(combined) != 5 {
		t.Errorf("got %d wanted %d", len(combined), total)
	}
}

func TestBuildHttpRequest(t *testing.T) {
	config := getConfig()
	reg := CreateRegistration(config)

	httpaction := CreateHTTPAction(config, reg)
	req := BuildHttpRequest(httpaction)
	if req == nil {
		t.Error("Request creation failed")
	}
}
