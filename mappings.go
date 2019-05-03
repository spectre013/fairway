package fairway

import (
	"strings"
)

func mappings(name string, routes Routes) ([]byte, error) {
	c := contexts{}
	c.Contexts = map[string]mapping{}
	c.Contexts[name] = mapping{}

	d := make([]dispatcherServlet, 0)

	for _, route := range routes {
		ds := dispatcherServlet{Handler: route.Name, Predicate: ""}
		details := details{}
		hm := handlerMethod{ClassName: route.Name, Name: route.Name, Descriptor: "(http.ResponseWriter,*http.Request)"}

		consumer := make([]consumes, 0)
		headers := make([]headers, 0)
		methods := make([]string, 0)
		params := make([]params, 0)
		patterns := make([]string, 0)
		produces := make([]consumes, 0)

		methods = append(methods, route.Method)
		patterns = append(patterns, route.Pattern)

		consumer = append(consumer, consumes{MediaType: route.Produces, Negated: false})
		if strings.Contains(route.Produces, "+json") {
			consumer = append(consumer, consumes{MediaType: "application/json", Negated: false})
		}

		rmc := requestMappingConditions{
			Consumes: consumer,
			Headers:  headers,
			Methods:  methods,
			Params:   params,
			Patterns: patterns,
			Produces: produces,
		}

		details.HandlerMethod = hm
		details.RequestMappingConditions = rmc
		ds.Details = details
		d = append(d, ds)
	}
	m := mapping{}
	dis := dispatcher{DispatcherServlet: d}
	m.Mappings.DispatcherServlets = dis
	m.Mappings.ParentID = appName
	c.Contexts[name] = m
	return toJSON(c), nil
}

type contexts struct {
	Contexts map[string]mapping `json:"contexts"`
}

type mapping struct {
	Mappings  mapData `json:"mappings"`
	Bootstrap mapData `json:"bootstrap"`
}

type mapData struct {
	DispatcherServlets dispatcher `json:"dispatcherServlets"`
	ParentID           string     `json:"parentId"`
}

type dispatcher struct {
	DispatcherServlet []dispatcherServlet `json:"dispatcherServlet"`
	ServerFilters     []filters           `json:"serverletFilters"`
	Serverlets        []serverlets        `json:"serverlets"`
}

type dispatcherServlet struct {
	Handler   string  `json:"handler"`
	Predicate string  `json:"predicate"`
	Details   details `json:"details"`
}

type details struct {
	HandlerMethod            handlerMethod            `json:"handlerMethod"`
	RequestMappingConditions requestMappingConditions `json:"requestMappingConditions"`
}

type handlerMethod struct {
	Name       string `json:"name"`
	ClassName  string `json:"className"`
	Descriptor string `json:"descriptor"`
}

type requestMappingConditions struct {
	Consumes []consumes `json:"consumes"`
	Headers  []headers  `json:"headers"`
	Methods  []string   `json:"methods"`
	Params   []params   `json:"params"`
	Patterns []string   `json:"patterns"`
	Produces []consumes `json:"produces"`
}

type consumes struct {
	MediaType string `json:"mediaType"`
	Negated   bool   `json:"negated"`
}

type headers struct {
}

type produces struct {
}
type params struct {
}

type filters struct {
	ServletNameMappings []string `json:"serverletMappings"`
	URLPatternMappings  []string `json:"urlPatternMappings"`
	Name                string   `json:"name"`
	ClassName           string   `json:"className"`
}

type serverlets struct {
	Mappings  []string `json:"mappings"`
	Name      string   `json:"name"`
	ClassName string   `json:"className"`
}
