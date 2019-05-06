package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spectre013/fairway"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":8900", "server listen address")
	flag.Parse()

	config := fairway.EurekaConfig{
		Name:        "eureka-test",
		URL:         "http://eurekaserver:8761/eureka",
		HostName:    "c00064.issinc.com",
		Port:        "8900",
		SecurePort:  "8943",
		RestService: true,
		PreferIP:    true,
		Username:    "user",
		Password:    "password",
		Secure:      false,
	}

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
			Handler:  http.FileServer(http.Dir("/Users/brian.paulson/go/src/github.com/spectre013/fairway/example/site")),
			Static:   true,
		},
	}

	routes = fairway.CombineRoutes(routes, eurekaRoutes)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	router := mux.NewRouter()

	router = fairway.BuildRoutes(routes, router)

	log.Println("Server is up and listening on ", listenAddr)
	err := http.ListenAndServe(listenAddr, router)
	if err != nil {
		logger.Fatal("Unable to start server", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
