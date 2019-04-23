package fairway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var startTime time.Time

type Eureka interface {
	Register()
	StartHeartbeat()
}

type EurekaRegistration struct {
	Instance Instance `json:"instance"`
}

type Instance struct {
	InstanceId     string     `json:"instanceId"`
	HostName       string     `json:"hostName"`
	App            string     `json:"app"`
	IpAddr         string     `json:"ipAddr"`
	VipAddress     string     `json:"vipAddress"`
	Status         string     `json:"status"`
	Port           Port       `json:"port"`
	SecurePort     Port       `json:"securePort"`
	HomePageUrl    string     `json:"homePageUrl"`
	StatusPageUrl  string     `json:"statusPageUrl"`
	HealthCheckUrl string     `json:"healthCheckUrl"`
	DataCenterInfo DataCenter `json:"dataCenterInfo"`
}

type Port struct {
	PortNumber string `json:"$"`
	Enabled    bool   `json:"@enabled"`
}

type DataCenter struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

var instanceId string
var eurekaURL string

func Register(config EurekaConfig) {
	startTime = time.Now()
	reg := CreateRegistration(config)
	registerAction := CreateHTTPAction(config, reg)

	logger.Debug(toJson(reg))

	var result bool
	for {
		logger.Info("Attempting to register with Eureka at ", config.Url)
		result = DoHttpRequest(registerAction)
		if result {
			logger.Info("Eureka registration successfull ... ")
			heartbeatStatus := make(chan bool)
			go func() { heartbeatStatus <- startHeartbeat(config) }() // Performs Eureka heartbeating (async)
			status := <-heartbeatStatus
			logger.Warn("Heartbeat request failed trying to reregister: heartbeat status -> ", status)
		} else {
			logger.Info("Eureka registration unsuccessfull or euraka is down will keep trying... ")
			time.Sleep(time.Second * 5)
		}
	}
}

func CreateHTTPAction(config EurekaConfig, reg EurekaRegistration) HttpAction {
	return HttpAction{
		Url:         config.Url + "/apps/" + config.Name,
		Method:      http.MethodPost,
		ContentType: "application/json",
		Body:        string(toJson(reg)),
	}
}

func CreateRegistration(config EurekaConfig) EurekaRegistration {
	instanceId = GetUUID()

	reg := EurekaRegistration{}
	port := Port{PortNumber: config.Port, Enabled: true}
	secureport := Port{PortNumber: config.SecurePort, Enabled: false}
	dataCenter := DataCenter{Class: "com.netflix.appinfo.MyDataCenterInfo", Name: "MyOwn"}
	instance := Instance{
		InstanceId:     config.Name + ":" + instanceId,
		HostName:       config.HostName,
		App:            strings.ToUpper(config.Name),
		IpAddr:         config.IpAddress,
		VipAddress:     config.VipAddress,
		Status:         "UP",
		Port:           port,
		SecurePort:     secureport,
		HomePageUrl:    fmt.Sprintf("http://%s:%s/", config.IpAddress, config.Port),
		StatusPageUrl:  fmt.Sprintf("http://%s:%s/actuator/info", config.IpAddress, config.Port),
		HealthCheckUrl: fmt.Sprintf("http://%s:%s/actuator/health", config.IpAddress, config.Port),
		DataCenterInfo: dataCenter}

	reg.Instance = instance
	return reg
}

func toJson(r interface{}) []byte {
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
	heartbeatAction := HttpAction{
		Url:    config.Url + "/apps/" + config.Name + "/" + config.Name + ":" + instanceId,
		Method: "PUT",
	}
	return DoHttpRequest(heartbeatAction)

}

func deregister(config EurekaConfig) {
	logger.Info("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		Url:    config.Url + "/apps/" + config.Name + "/" + config.Name + ":" + instanceId,
		Method: "DELETE",
	}
	DoHttpRequest(deregisterAction)
	logger.Info("Deregistered application, exiting.")
}
