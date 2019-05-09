// +build ignore
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spectre013/fairway"
	"log"
	"net/http"
)

var (
	conf string
)

func main() {

	flag.StringVar(&conf, "conf", "./conf.yml", "Yaml Configuration")
	flag.Parse()

	config := fairway.GetFile(conf)
	eureka := fairway.Init(config)
	eurekaRoutes := eureka.Routes

	var routes = fairway.Routes{
		fairway.Route{
			Name:     "Index",
			Method:   "GET",
			Pattern:  "/index",
			Produces: "text/html",
			Handler:  http.HandlerFunc(index),
			Static:   true,
		},
		fairway.Route{
			Name:     "Index",
			Method:   "GET",
			Pattern:  "/",
			Produces: "text/html",
			Handler:  http.FileServer(http.Dir(config.ServeDir)),
			Static:   true,
		},
	}

	routes = fairway.CombineRoutes(routes, eurekaRoutes)

	fmt.Println("Server is starting...")

	router := mux.NewRouter()

	router = fairway.BuildRoutes(routes, router)

	fmt.Println("Server is up and listening on ", config.Port)
	err := http.ListenAndServe(":"+config.Port, router)
	if err != nil {
		log.Fatal("Unable to start server", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
