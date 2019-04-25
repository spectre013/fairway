package fairway

import (
	"os"
	"strconv"
	"strings"
)

type EnvData struct {
	ActiveProfiles  []string  `json:"activeProfiles"`
	Property        Property  `json:"property,omitempty"`
	PropertySources []Sources `json:"propertySources"`
}

type Sources struct {
	Name     string              `json:"name"`
	Property map[string]Property `json:"properties"`
}

type Property struct {
	Value  string `json:"value,omitempty"`
	Origin string `json:"origin,omitempty"`
	Source string `json:"source,omitempty"`
}

func env(property string) ([]byte, error) {
	env := EnvData{ActiveProfiles: []string{"go"}}
	env.PropertySources = make([]Sources, 0)
	sysProps := getSystemProperties(property)
	sysEnv := getSystemEnvironmentProperties(property)

	if property != "" {
		env.Property = getProperty(property, sysProps, sysEnv)
	}

	env.PropertySources = append(env.PropertySources, sysProps)
	env.PropertySources = append(env.PropertySources, sysEnv)
	env.PropertySources = append(env.PropertySources, getProperties("servletContextInitParams",""))
	env.PropertySources = append(env.PropertySources, getProperties("defaultProperties",""))
	return toJson(env), nil
}

func getProperty(prop string, props Sources, sysEnv Sources) Property {
	p := Property{}
	if val, ok := props.Property[prop]; ok {
		p = val
		p.Source = props.Name
	}

	if val, ok := sysEnv.Property[prop]; ok {
		p.Source = props.Name
		p = val
	}
	return p
}

func getSystemEnvironmentProperties(prop string) Sources {
	source := Sources{}
	source.Name = "systemEnvironment"
	p := map[string]Property{}
	for _, e := range os.Environ() {
		e := strings.Split(e, "=")
		if strings.HasPrefix(e[0], "_") {
			continue
		}
		if !strings.Contains(e[0], "PASSWORD") {
			p[strings.ToUpper(e[0])] = Property{Value: e[1],
				Origin: "System Environment Property \"" + e[0] + "\"",
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
		p["PID"] = Property{Value: strconv.Itoa(os.Getpid())}
	source.Property = p
	return source
}

func getProperties(title, prop string) Sources {
	source := Sources{}
	source.Name = title
	p := map[string]Property{}
	source.Property = p
	return source
}
