package fairway

import (
	"os"
	"strconv"
	"strings"
)

type envData struct {
	ActiveProfiles  []string  `json:"activeProfiles"`
	Property        property  `json:"property,omitempty"`
	PropertySources []sources `json:"propertySources"`
}

type sources struct {
	Name     string              `json:"name"`
	Property map[string]property `json:"properties"`
}

type property struct {
	Value  string `json:"value,omitempty"`
	Origin string `json:"origin,omitempty"`
	Source string `json:"source,omitempty"`
}

func env(property string) ([]byte, error) {
	env := envData{ActiveProfiles: []string{"go"}}
	env.PropertySources = make([]sources, 0)
	sysProps := getSystemProperties(property)
	sysEnv := getSystemEnvironmentProperties(property)

	if property != "" {
		env.Property = getProperty(property, sysProps, sysEnv)
	}

	env.PropertySources = append(env.PropertySources, sysProps)
	env.PropertySources = append(env.PropertySources, sysEnv)
	env.PropertySources = append(env.PropertySources, getProperties("servletContextInitParams", ""))
	env.PropertySources = append(env.PropertySources, getProperties("defaultProperties", ""))
	return toJSON(env), nil
}

func getProperty(prop string, props sources, sysEnv sources) property {
	p := property{}
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

func getSystemEnvironmentProperties(prop string) sources {
	source := sources{}
	source.Name = "systemEnvironment"
	p := map[string]property{}
	for _, e := range os.Environ() {
		e := strings.Split(e, "=")
		if strings.HasPrefix(e[0], "_") {
			continue
		}
		if !strings.Contains(e[0], "PASSWORD") {
			p[strings.ToUpper(e[0])] = property{Value: e[1],
				Origin: "System Environment Property \"" + e[0] + "\"",
			}
		}
	}
	if prop != "" {
		if val, ok := p[prop]; ok {
			tp := map[string]property{}
			tp[prop] = val
			source.Property = tp
		}
	} else {
		source.Property = p
	}
	return source
}

func getSystemProperties(prop string) sources {
	source := sources{}
	source.Name = "systemProperties"
	p := map[string]property{}
	p["PID"] = property{Value: strconv.Itoa(os.Getpid())}
	source.Property = p
	return source
}

func getProperties(title, prop string) sources {
	source := sources{}
	source.Name = title
	p := map[string]property{}
	source.Property = p
	return source
}
