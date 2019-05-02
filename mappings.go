package fairway

import (
	"strings"
)

func mappings(name string, routes Routes) ([]byte, error) {
	c := Contexts{}
	c.Contexts = map[string]Mapping{}
	c.Contexts[name] = Mapping{}

	d := make([]DispatcherServlet, 0)

	for _, route := range routes {
		ds := DispatcherServlet{Handler: route.Name, Predicate: ""}
		details := Details{}
		hm := HandlerMethod{ClassName: route.Name, Name: route.Name, Descriptor: "(http.ResponseWriter,*http.Request)"}

		consumes := make([]Consumes, 0)
		headers := make([]Headers, 0)
		methods := make([]string, 0)
		params := make([]Params, 0)
		patterns := make([]string, 0)
		produces := make([]Consumes, 0)

		methods = append(methods, route.Method)
		patterns = append(patterns, route.Pattern)

		consumes = append(consumes, Consumes{MediaType: route.Produces, Negated: false})
		if strings.Contains(route.Produces, "+json") {
			consumes = append(consumes, Consumes{MediaType: "application/json", Negated: false})
		}

		rmc := RequestMappingConditions{
			Consumes: consumes,
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
	m := Mapping{}
	dis := Dispatcher{DispatcherServlet: d}
	m.Mappings.DispatcherServlets = dis
	m.Mappings.ParentID = appName
	c.Contexts[name] = m
	return toJson(c), nil
}

type Contexts struct {
	Contexts map[string]Mapping `json:"contexts"`
}

type Mapping struct {
	Mappings  MapData `json:"mappings"`
	Bootstrap MapData `json:"bootstrap"`
}

type MapData struct {
	DispatcherServlets Dispatcher `json:"dispatcherServlets"`
	ParentID           string     `json:"parentId"`
}

type Dispatcher struct {
	DispatcherServlet []DispatcherServlet `json:"dispatcherServlet"`
	ServerFilters     []Filters           `json:"serverletFilters"`
	Serverlets        []Serverlets        `json:"serverlets"`
}

type DispatcherServlet struct {
	Handler   string  `json:"handler"`
	Predicate string  `json:"predicate"`
	Details   Details `json:"details"`
}

type Details struct {
	HandlerMethod            HandlerMethod            `json:"handlerMethod"`
	RequestMappingConditions RequestMappingConditions `json:"requestMappingConditions"`
}

type HandlerMethod struct {
	Name       string `json:"name"`
	ClassName  string `json:"className"`
	Descriptor string `json:"descriptor"`
}

type RequestMappingConditions struct {
	Consumes []Consumes `json:"consumes"`
	Headers  []Headers  `json:"headers"`
	Methods  []string   `json:"methods"`
	Params   []Params   `json:"params"`
	Patterns []string   `json:"patterns"`
	Produces []Consumes `json:"produces"`
}

type Consumes struct {
	MediaType string `json:"mediaType"`
	Negated   bool   `json:"negated"`
}

type Headers struct {
}

type Produces struct {
}
type Params struct {
}

type Filters struct {
	ServletNameMappings []string `json:"serverletMappings"`
	UrlPatternMappings  []string `json:"urlPatternMappings"`
	Name                string   `json:"name"`
	ClassName           string   `json:"className"`
}

type Serverlets struct {
	Mappings  []string `json:"mappings"`
	Name      string   `json:"name"`
	ClassName string   `json:"className"`
}
