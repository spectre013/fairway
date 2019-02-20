package main

import (
	"flag"
	"github.com/spectre013/goeureka"
	"log"
	"net/http"
	"os"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":8900", "server listen address")
	flag.Parse()

	config := goeureka.EurekaConfig{
		Name:        "tb-ui-search",
		Url:         "http://eurekaserver:8761/eureka",
		HostName:    "c00064.issinc.com",
		Port:        "8900",
		SecurePort:  "8943",
		RestService: true,
		PreferIP: true,
	}



	eureka := goeureka.Init(config)
	eurekaRoutes := eureka.Routes

	var routes = goeureka.Routes{
		goeureka.Route{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: nil,
			Handler:     http.FileServer(http.Dir("/Users/brian.paulson/pa/tb/tb-ui/tb-search/dist/search")),
		},
	}

	routes = goeureka.CombineRoutes(routes, eurekaRoutes)

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	router := http.NewServeMux()

	//fs := http.FileServer(http.Dir("dist/search"))
	//http.Handle("/", Log(fs))

	router = goeureka.BuildRoutes(routes, router)

	log.Println("Server is up and listening on ", listenAddr)
	http.ListenAndServe(listenAddr, router)
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
