package fairway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var startTime time.Time

// Eureka Interface
type Eureka interface {
	Register()
	StartHeartbeat()
}

// EurekaRegistration data type
type EurekaRegistration struct {
	Instance instance `json:"instance"`
}

type instance struct {
	InstanceID     string     `json:"instanceId"`
	HostName       string     `json:"hostName"`
	App            string     `json:"app"`
	IPADDR         string     `json:"ipAddr"`
	VipAddress     string     `json:"vipAddress"`
	Status         string     `json:"status"`
	Port           port       `json:"port"`
	SecurePort     port       `json:"securePort"`
	HomePageURL    string     `json:"homePageUrl"`
	StatusPageURL  string     `json:"statusPageUrl"`
	HealthCheckURL string     `json:"healthCheckUrl"`
	DataCenterInfo dataCenter `json:"dataCenterInfo"`
}

type port struct {
	PortNumber string `json:"$"`
	Enabled    bool   `json:"@enabled"`
}

type dataCenter struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

var instanceID string
var eurekaURL string
var appName string

// Register Rgister application with EurekA
func Register(config EurekaConfig) {
	startTime = time.Now()
	reg := createRegistration(config)
	registerAction := createHTTPAction(config, reg)

	appName = config.Name

	logger.Debug(string(toJSON(reg)))

	var result bool
	for {
		logger.Info("Attempting to register with Eureka at ", config.URL)
		result = DoHTTPRequest(registerAction)
		if result {
			logger.Info("Eureka registration successful ... ")
			heartbeatStatus := make(chan bool)
			go func() { heartbeatStatus <- startHeartbeat(config) }() // Performs Eureka heartbeating (async)
			status := <-heartbeatStatus
			logger.Warn("Heartbeat request failed trying to reregister: heartbeat status -> ", status)
		} else {
			logger.Info("Eureka registration successful or euraka is down will keep trying... ")
			time.Sleep(time.Second * 5)
		}
	}
}

func createHTTPAction(config EurekaConfig, reg EurekaRegistration) HTTPAction {
	return HTTPAction{
		URL:         config.URL + "/apps/" + config.Name,
		Method:      http.MethodPost,
		ContentType: "application/json",
		Body:        string(toJSON(reg)),
	}
}

func createRegistration(config EurekaConfig) EurekaRegistration {
	instanceID = getUUID()

	reg := EurekaRegistration{}
	unsecurePort := port{PortNumber: config.Port, Enabled: true}
	secureport := port{PortNumber: config.SecurePort, Enabled: false}
	dataCenter := dataCenter{Class: "com.netflix.appinfo.MyDataCenterInfo", Name: "MyOwn"}
	instance := instance{
		InstanceID:     config.Name + ":" + instanceID,
		HostName:       config.HostName,
		App:            strings.ToUpper(config.Name),
		IPADDR:         config.IPAddress,
		VipAddress:     config.VipAddress,
		Status:         "UP",
		Port:           unsecurePort,
		SecurePort:     secureport,
		HomePageURL:    fmt.Sprintf("http://%s:%s/", config.IPAddress, config.Port),
		StatusPageURL:  fmt.Sprintf("http://%s:%s/actuator/info", config.IPAddress, config.Port),
		HealthCheckURL: fmt.Sprintf("http://%s:%s/actuator/health", config.IPAddress, config.Port),
		DataCenterInfo: dataCenter}

	reg.Instance = instance
	return reg
}

func toJSON(r interface{}) []byte {
	f, err := json.Marshal(r)
	if err != nil {
		logger.Error("error:", err)
	}
	return f
}

func startHeartbeat(config EurekaConfig) bool {
	for {
		time.Sleep(time.Second * 30)
		status := heartbeat(config)
		if status == false {
			break
		}
	}
	return false
}

func heartbeat(config EurekaConfig) bool {
	heartbeatAction := HTTPAction{
		URL:    config.URL + "/apps/" + config.Name + "/" + config.Name + ":" + instanceID,
		Method: "PUT",
	}
	return DoHTTPRequest(heartbeatAction)

}

func deregister(config EurekaConfig) {
	logger.Info("Trying to deregister application...")
	// Deregister
	deregisterAction := HTTPAction{
		URL:    config.URL + "/apps/" + config.Name + "/" + config.Name + ":" + instanceID,
		Method: "DELETE",
	}
	DoHTTPRequest(deregisterAction)
	logger.Info("Deregistered application, exiting.")
}
