package goeureka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

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
	instanceId = GetUUID()

	reg := EurekaRegistration{}
	port := Port{PortNumber: ":" + config.Port, Enabled: true}
	secureport := Port{PortNumber: ":" + config.SecurePort, Enabled: false}
	dataCenter := DataCenter{Class: "com.netflix.appinfo.MyDataCenterInfo", Name: "MyOwn"}
	instance := Instance{
		InstanceId:     config.Name + ":" + instanceId,
		HostName:       config.HostName,
		App:            config.Name,
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

	// Register.
	registerAction := HttpAction{
		Url:         config.Url + "/apps/" + config.Name,
		Method:      "POST",
		ContentType: "application/json; charset=utf-8",
		Body:        toJson(reg),
	}

	var result bool
	for {
		fmt.Println("Attempting to register with Eureka at ", config.Url)
		fmt.Println(registerAction)
		result = DoHttpRequest(registerAction)
		log.Println(result)
		if result {
			go StartHeartbeat(config) // Performs Eureka heartbeating (async)
			fmt.Println("Eureka registration successfull ... ")
			break
		} else {
			fmt.Println("Eureka registration unsuccessfull or euraka is down will keep trying... ")
			time.Sleep(time.Second * 5)
		}
	}
}

func toJson(r EurekaRegistration) string {
	f, err := json.Marshal(r)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(f)
}

func StartHeartbeat(config EurekaConfig) {
	for {
		time.Sleep(time.Second * 30)
		heartbeat(config)
	}
}

func heartbeat(config EurekaConfig) {
	heartbeatAction := HttpAction{
		Url:    config.Url + "/apps/" + config.Name + "/" + GetLocalIP() + ":" + config.Name + ":" + instanceId,
		Method: "PUT",
	}
	DoHttpRequest(heartbeatAction)
}

func Deregister(config EurekaConfig) {
	fmt.Println("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		Url:    config.Url + "/apps/" + config.Name + "/" + config.Name + ":" + instanceId,
		Method: "DELETE",
	}
	DoHttpRequest(deregisterAction)
	fmt.Println("Deregistered application, exiting.")
}
