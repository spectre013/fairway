package fairway

type links struct {
	Link map[string]link `json:"_links"`
}

type link struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}

func acuator(host string) ([]byte, error) {

	links := links{}

	m := map[string]link{}

	m["self"] = link{Href: "http://" + host + "/actuator", Templated: false}
	m["health"] = link{Href: "http://" + host + "/actuator/health", Templated: false}
	m["info"] = link{Href: "http://" + host + "/actuator/info", Templated: false}
	m["env"] = link{Href: "http://" + host + "/actuator/env", Templated: false}
	m["env-toMatch"] = link{Href: "http://" + host + "/actuator/env/{toMatch}", Templated: true}
	m["metrics"] = link{Href: "http://" + host + "/actuator/metrics", Templated: true}
	m["metrics-requiredMetricName"] = link{Href: "http://" + host + "/actuator/env/{requiredMetricName}", Templated: true}
	m["loggers"] = link{Href: "http://" + host + "/actuator/loggers", Templated: false}
	m["loggers-name"] = link{Href: "http://" + host + "/actuator/loggers/{name}", Templated: true}
	m["mappings"] = link{Href: "http://" + host + "/actuator/mappings", Templated: true}

	links.Link = m
	return toJSON(links), nil
}
