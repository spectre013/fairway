package goeureka

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
)

type EurekaClient struct {
	Client Eureka
	Routes Routes
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

	config.IpAddress = GetOutboundIP().String()
	config.VipAddress = strings.ToUpper(config.Name)
	handleSigterm(config) // Graceful shutdown on Ctrl+C or kill
	routes := routes
	go Register(config) // Performs Eureka registration
	// start server and Block if not a rest service...
	if !config.RestService {
		go startWebServer(routes, config.Port)
		wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
		wg.Add(1)
		wg.Wait()
	}

	var e Eureka
	return EurekaClient{Client: e, Routes: routes}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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

func CombineRoutes(routes Routes, eurekaRouts Routes) Routes {

	for _, route := range eurekaRouts {
		routes = append(routes, route)
	}
	return routes
}

func PrintRoutes() {
	//for _, route := range e{
	//	fmt.Println(fmt.Sprintf("Mapped (%s) with method (%s) to %s", route.Path, route.Method, route.Name))
	//}
	httpMux := reflect.ValueOf(http.DefaultServeMux).Elem()
	finList := httpMux.FieldByIndex([]int{1})
	fmt.Println(finList)
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func startWebServer(router Routes, port string) {
	e := http.NewServeMux()
	e = BuildRoutes(routes, e)
	// //log.Println("Starting HTTP service at " + port)
	// e.Logger.Fatal(e.Start(":" + port))
}
