package fairway

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Error(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Error Route not found", http.StatusNotFound)
}

func Info(w http.ResponseWriter, r *http.Request) {
	json, err := info()
	writeResponse(json, err, w)
}

func Actuator(w http.ResponseWriter, r *http.Request) {
	actuator, err := acuator(r.Host)
	writeResponse(actuator, err, w)
}

func Loggers(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	var logs []byte
	var err error
	var log LogStruct
	if r.Method == "POST" {
		key := ""
		if len(path) > 3 {
			key = path[3]
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeResponse(nil, err, w)
		}
		err = json.Unmarshal(bodyBytes, &log)
		if err != nil {
			writeResponse(nil, err, w)
		}
		logs, err = loggersUpdate(key, log)
	} else {
		logs, err = loggers()
	}
	writeResponse(logs, err, w)
}

func Health(w http.ResponseWriter, r *http.Request) {
	health, err := health()
	writeResponse(health, err, w)
}

func Env(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	prop := ""
	if len(path) > 3 {
		prop = path[3]
	}
	json, err := env(prop)
	writeResponse(json, err, w)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	metric := ""
	if len(path) > 3 {
		metric = path[3]
	}
	json, err := metrics(metric, r.URL.Query())
	writeResponse(json, err, w)
}

func writeResponse(data []byte, err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v2+json;charset=UTF-8")

	w.Write(data)
}
