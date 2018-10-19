package goeureka

import (
	"encoding/json"
	"fmt"
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

func Register(name string, eurekaPath string, vip_address string) {
	instanceId = GetUUID()
	eurekaURL = eurekaPath

	reg := EurekaRegistration{}
	port := Port{PortNumber: ":8080", Enabled: true}
	secureport := Port{PortNumber: ":8443", Enabled: false}
	dataCenter := DataCenter{Class: "com.netflix.appinfo.MyDataCenterInfo", Name: "MyOwn"}
	instance := Instance{
		InstanceId:     name + ":" + instanceId,
		HostName:       "c00064.issinc.com",
		App:            name,
		IpAddr:         "172.20.3.154",
		VipAddress:     vip_address,
		Status:         "UP",
		Port:           port,
		SecurePort:     secureport,
		HomePageUrl:    "http://172.20.3.154:8181/",
		StatusPageUrl:  "http://172.20.3.154:8181/info",
		HealthCheckUrl: "http://172.20.3.154:8181/health",
		DataCenterInfo: dataCenter}

	reg.Instance = instance

	// Register.
	registerAction := HttpAction{
		Url:         eurekaURL + "/apps/" + name,
		Method:      "POST",
		ContentType: "application/json; charset=utf-8",
		Body:        toJson(reg),
	}
	fmt.Println(registerAction)

	var result bool
	for {
		fmt.Println("Attempting to register with Eureka at ", eurekaPath)
		result = DoHttpRequest(registerAction)
		if result {
			go StartHeartbeat(name) // Performs Eureka heartbeating (async)
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

func StartHeartbeat(name string) {
	for {
		time.Sleep(time.Second * 30)
		fmt.Println("sending heartbeat ...")
		heartbeat(name)
	}
}

func heartbeat(name string) {
	heartbeatAction := HttpAction{
		Url:    eurekaURL + "/apps/" + name + "/" + GetLocalIP() + ":" + name + ":" + instanceId,
		Method: "PUT",
	}
	DoHttpRequest(heartbeatAction)
}

func Deregister(name string) {
	fmt.Println("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		Url:    eurekaURL + "/apps/" + name + "/" + name + ":" + instanceId,
		Method: "DELETE",
	}
	DoHttpRequest(deregisterAction)
	fmt.Println("Deregistered application, exiting. Check Eureka...")
}
