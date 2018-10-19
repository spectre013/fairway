package goeureka

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Eureka interface {
	Register()
	StartHeartbeat()
}

var instanceId string
var eurekaURL string

func Register(name string, eurekaPath string) {
	instanceId = GetUUID()
	eurekaURL = eurekaPath

	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(dir + "/reg.json")

	tpl := string(data)
	tpl = strings.Replace(tpl, "${ipAddress}", GetLocalIP(), -1)
	tpl = strings.Replace(tpl, "${app}", name, -1)
	tpl = strings.Replace(tpl, "${port}", "8080", -1)
	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)
	fmt.Println(tpl)
	// Register.
	registerAction := HttpAction{
		Url:         eurekaURL + "/apps/" + name,
		Method:      "POST",
		ContentType: "application/json; charset=utf-8",
		Body:        tpl,
	}
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
