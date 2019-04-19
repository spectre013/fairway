# Go Eureka

a simple Eureka client for go that will register your service with eureka and send hearbeats as well as deregister on shutdown or panic. 

There are two ways to implement fairway, one as a rest service and as a stand alone service. 



## Rest service implementation
```golang
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/spectre013/fairway"
)

func main() {
	eureka_host := os.Getenv("EUREKA_HOST")
	service_name := os.Getenv("SERVICE_NAME")
	vip_address := os.Getenv("VIP_ADDRESS")
	ip_address := os.Getenv("IP_ADDRESS")
	host_name := os.Getenv("HOST_NAME")
	port := os.Getenv("PORT")
	secure_port := os.Getenv("SECURE_PORT")
	rest_service, _ := strconv.ParseBool(os.Getenv("REST_SERVICE"))

	config := fairway.EurekaConfig{
		Name:        service_name,
		Url:         eureka_host,
		VipAddress:  vip_address,
		IpAddress:   ip_address,
		HostName:    host_name,
		Port:        port,
		SecurePort:  secure_port,
		RestService: rest_service,
	}

	eureka := fairway.Init(config)
	eurekaRoutes := eureka.Routes

	var routes = fairway.Routes{
		fairway.Route{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: Index,
		},
	}
	e := echo.New()
	routes = fairway.CombineRoutes(routes, eurekaRoutes)
	e = fairway.BuildRoutes(routes, e)
	startServer(port, e)
}

func startServer(port string, e *echo.Echo) {
	log.Println("Starting HTTP service at " + port)
	fairway.PrintRoutes(e)
	e.Logger.Fatal(e.Start(":" + port))
}

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

```

## Service implementation
```golang
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/spectre013/fairway"
)

func main() {
	eureka_host := os.Getenv("EUREKA_HOST")
	service_name := os.Getenv("SERVICE_NAME")
	vip_address := os.Getenv("VIP_ADDRESS")
	ip_address := os.Getenv("IP_ADDRESS")
	host_name := os.Getenv("HOST_NAME")
	port := os.Getenv("PORT")
	secure_port := os.Getenv("SECURE_PORT")
	rest_service, _ := strconv.ParseBool(os.Getenv("REST_SERVICE"))

	config := fairway.EurekaConfig{
		Name:        service_name,
		Url:         eureka_host,
		VipAddress:  vip_address,
		IpAddress:   ip_address,
		HostName:    host_name,
		Port:        port,
		SecurePort:  secure_port,
		RestService: rest_service,
	}

	go fairway.Init(config)

	dosomething()

}

func dosomething() {
	forever := make(chan bool)
	go func() {

	}()

	log.Printf(" [*] Waiting... To exit press CTRL+C")
	<-forever
}

```
