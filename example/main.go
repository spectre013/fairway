package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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
		Url:         "http://eurekaserver:8761/eureka",
		HostName:    "c00064.issinc.com",
		Port:        "8900",
		SecurePort:  "8943",
		RestService: true,
		PreferIP:    true,
	}

	eureka := fairway.Init(config)
	eurekaRoutes := eureka.Routes

	var routes = fairway.Routes{
		fairway.Route{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: index,
			Handler:     nil, //http.FileServer(http.Dir("/Users/brian.paulson/pa/tb/tb-ui/tb-search/dist/search")),
		},
	}

	routes = fairway.CombineRoutes(routes, eurekaRoutes)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	router := http.NewServeMux()

	router = fairway.BuildRoutes(routes, router)

	log.Println("Server is up and listening on ", listenAddr)
	http.ListenAndServe(listenAddr, router)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
