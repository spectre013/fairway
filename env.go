package fairway

import (
	"encoding/json"
	"os"
	"strings"
)

func env() ([]byte, error) {
	env := new(envObject)
	env.SystemEnvironment = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env.SystemEnvironment[pair[0]] = pair[1]
	}

	b, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type envObject struct {
	Configuration     map[string]interface{} `json:"configuration"`
	SystemEnvironment map[string]string      `json:"systemEnvironment"`
}
