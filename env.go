package goeureka

import (
	"os"
	"strings"
)

func env() *envObject {
	env := new(envObject)
	env.SystemEnvironment = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env.SystemEnvironment[pair[0]] = pair[1]
	}

	return env
}

type envObject struct {
	Configuration map[string]interface{} `json:"configuration"`
	SystemEnvironment map[string]string `json:"systemEnvironment"`
}