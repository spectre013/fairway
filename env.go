package fairway

import (
	"os"
	"strconv"
	"strings"
)

type EnvData struct {
	Property Property `json:"property,omitempty"`
	PropertySources []Sources `json:"propertySources"`

}

type Sources struct {
	Name string `json:"name"`
	Property map[string]Property `json:"property,omitempty"`
}

type Property struct {
	Value string `json:"value,omitempty"`
	Origin string `json:"origin,omitempty"`
	Source string `json:"source,omitempty"`
}

func env(property string) ([]byte, error) {
	env := EnvData{}
	env.PropertySources = make([]Sources,0)
	sysProps := getSystemProperties(property)
	sysEnv := getSystemEnvironmentProperties(property)

	if property != "" {
		env.Property = getPropery(property,sysProps,sysEnv)
	}

	env.PropertySources = append(env.PropertySources,sysProps)
	env.PropertySources = append(env.PropertySources,sysEnv)

	return toJson(env), nil
}

func getPropery(prop string, props Sources, env Sources) Property {
	property := Property{}
	if val, ok := props.Property[prop]; ok {
		property = val
		property.Source = props.Name
	}

	if val, ok := env.Property[prop]; ok {
		property.Source = props.Name
		property = val
	}
	return property
}

func getSystemEnvironmentProperties(prop string) Sources {
	source := Sources{}
	source.Name = "systemEnvironment"
	p := map[string]Property{}
	for _, e := range os.Environ() {
		env := strings.Split(e, "=")
		if !strings.Contains(env[0],"PASSWORD") {
			p[env[0]] = Property{Value: env[1],
				Origin: "System Environment Property \"" + env[0] + "\"",
			}
		}
	}
	if prop != "" {
		if val, ok := p[prop]; ok {
			tp := map[string]Property{}
			tp[prop] = val
			source.Property = tp
		}
	} else {
		source.Property = p
	}
	return source
}

func getSystemProperties(prop string) Sources {
	source := Sources{}
	source.Name = "systemProperties"
	p := map[string]Property{}
	if prop == "PID" {
		p["PID"] = Property{Value: strconv.Itoa(os.Getpid())}
	}
	source.Property = p
	return source
}