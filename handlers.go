package fairway

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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
func UpdateLogger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	bodyBytes, err := ioutil.ReadAll(r.Body)
	var log LogStruct
	var logs []byte
	if err != nil {
		writeResponse(nil, err, w)
	}
	err = json.Unmarshal(bodyBytes, &log)
	if err != nil {
		writeResponse(nil, err, w)
	}
	logs, err = loggersUpdate(key, log)
	writeResponse(logs, err, w)
}
func Loggers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]

	logs, err := loggers(key)
	writeResponse(logs, err, w)
}

func Health(w http.ResponseWriter, r *http.Request) {
	health, err := health()
	writeResponse(health, err, w)
}

func Mappings(w http.ResponseWriter, r *http.Request) {
	m, err := mappings(appName, appRoutes)
	writeResponse(m, err, w)
}

func Env(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prop := vars["toMatch"]
	json, err := env(prop)
	writeResponse(json, err, w)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	//requiredMetricName
	vars := mux.Vars(r)
	metric := vars["requiredMetricName"]
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
