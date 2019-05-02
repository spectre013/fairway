package fairway

type Links struct {
	Link map[string]Link `json:"_links"`
}

type Link struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}

func acuator(host string) ([]byte, error) {

	links := Links{}

	m := map[string]Link{}

	m["self"] = Link{Href: "http://" + host + "/actuator", Templated: false}
	m["health"] = Link{Href: "http://" + host + "/actuator/health", Templated: false}
	m["info"] = Link{Href: "http://" + host + "/actuator/info", Templated: false}
	m["env"] = Link{Href: "http://" + host + "/actuator/env", Templated: false}
	m["env-toMatch"] = Link{Href: "http://" + host + "/actuator/env/{toMatch}", Templated: true}
	m["metrics"] = Link{Href: "http://" + host + "/actuator/metrics", Templated: true}
	m["metrics-requiredMetricName"] = Link{Href: "http://" + host + "/actuator/env/{requiredMetricName}", Templated: true}
	m["loggers"] = Link{Href: "http://" + host + "/actuator/loggers", Templated: false}
	m["loggers-name"] = Link{Href: "http://" + host + "/actuator/loggers/{name}", Templated: true}
	m["mappings"] = Link{Href: "http://" + host + "/actuator/mappings", Templated: true}

	links.Link = m
	return toJson(links), nil
}
