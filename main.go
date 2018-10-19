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

type EurekaConfig struct {
	Name        string
	Url         string
	VipAddress  string
	IpAddress   string
	HostName    string
	Port        string
	SecurePort  string
	RestService bool
}

func Init(config EurekaConfig) EurekaClient {
	handleSigterm(config) // Graceful shutdown on Ctrl+C or kill
	router := buildRouter()
	go Register(config) // Performs Eureka registration
	// start server and Block if not a rest service...
	if !config.RestService {
		go startWebServer(router)
		wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
		wg.Add(1)
		wg.Wait()
	}

	var e Eureka
	return EurekaClient{Client: e, Router: router}
}

func handleSigterm(config EurekaConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		Deregister(config)
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
