package goeureka

import (
	"encoding/json"
	"net/http"
)


func Info(w http.ResponseWriter, r *http.Request) {
	json, err := info()
	if err != nil {
		http.Error(w,"Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(json)
}

func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		health := map[string]string{"status": "UP"}
		json, err := json.Marshal(health)
		if err != nil {
			http.Error(w,"Internal Server Error", http.StatusInternalServerError)
		}
		w.Write(json)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func Env(w http.ResponseWriter, r *http.Request) {
	json, err := env()
	if err != nil {
		http.Error(w,"Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(json)
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	json, err := metrics()
	if err != nil {
		http.Error(w,"Internal Server Error", http.StatusInternalServerError)
	}
	w.Write(json)
}
