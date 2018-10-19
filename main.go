package goeureka

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type EurekaClient struct {
	Client Eureka
	Router *mux.Router
}

func Init(name string, eurekaPath string, vip_address string, restService bool) EurekaClient {
	log.Println(name, eurekaPath)
	handleSigterm(name) // Graceful shutdown on Ctrl+C or kill
	router := buildRouter()
	go Register(name, eurekaPath, vip_address) // Performs Eureka registration
	// start server and Block if not a rest service...
	if !restService {
		go startWebServer(router)
		wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
		wg.Add(1)
		wg.Wait()
	}

	var e Eureka
	return EurekaClient{Client: e, Router: router}
}

func handleSigterm(name string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		Deregister(name)
		os.Exit(1)
	}()
}

func buildRouter() *mux.Router {
	return NewRouter()
}

func startWebServer(router *mux.Router) {
	log.Println("Starting HTTP service at 23456")
	srv := &http.Server{
		Handler: router,
		Addr:    ":23456",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
