package eureka

import (
	"fmt"
	"goeureka/util"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var instanceId string

func Register() {
	instanceId = util.GetUUID()

	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(dir + "/reg.json")

	tpl := string(data)
	tpl = strings.Replace(tpl, "${ipAddress}", util.GetLocalIP(), -1)
	tpl = strings.Replace(tpl, "${port}", "8080", -1)
	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)

	// Register.
	registerAction := HttpAction{
		Url:         "http://localhost:8761/eureka/apps/vendor",
		Method:      "POST",
		ContentType: "application/json; charset=utf-8",
		Body:        tpl,
	}
	var result bool
	for {
		fmt.Println("Attempting to register with Eureka ...")
		fmt.Println(registerAction)
		result = DoHttpRequest(registerAction)
		if result {
			fmt.Println(result)
			fmt.Println("Eureka registration successfull ... ")
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

func StartHeartbeat() {
	for {
		time.Sleep(time.Second * 30)
		fmt.Println("sending heartbeat ...")
		heartbeat()
	}
}

func heartbeat() {
	heartbeatAction := HttpAction{
		Url:    "http://localhost:8761/eureka/apps/vendor/" + util.GetLocalIP() + ":vendor:" + instanceId,
		Method: "PUT",
	}
	DoHttpRequest(heartbeatAction)
}

func Deregister() {
	fmt.Println("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		Url:    "http://localhost:8761/eureka/apps/vendor/" + util.GetLocalIP() + ":vendor:" + instanceId,
		Method: "DELETE",
	}
	DoHttpRequest(deregisterAction)
	fmt.Println("Deregistered application, exiting. Check Eureka...")
}
