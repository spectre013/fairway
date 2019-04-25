package fairway

type Links struct {
	Link map[string]Link `json:"_links"`
}

type Link struct {
	Href      string `json:"href"`
	templated bool   `json:"templated"`
}

func acuator(host string) ([]byte, error) {

	links := Links{}

	m := map[string]Link{}

	m["self"] = Link{Href: "http://" + host + "/actuator", templated: false}
	m["health"] = Link{Href: "http://" + host + "/actuator/health", templated: false}
	m["info"] = Link{Href: "http://" + host + "/actuator/info", templated: false}
	m["env"] = Link{Href: "http://" + host + "/actuator/env", templated: false}
	m["env-toMatch"] = Link{Href: "http://" + host + "/actuator/env/{toMatch}", templated: true}
	m["metrics"] = Link{Href: "http://" + host + "/actuator/metrics", templated: true}
	m["metrics-requiredMetricName"] = Link{Href: "http://" + host + "/actuator/env/{requiredMetricName}", templated: true}
	m["loggers"] = Link{Href: "http://" + host + "/actuator/loggers", templated: false}
	m["loggers-name"] = Link{Href: "http://" + host + "/actuator/loggers/{name}", templated: true}
	m["mappings"] = Link{Href: "http://" + host + "/actuator/mappings", templated: true}

	links.Link = m
	return toJson(links), nil
}
