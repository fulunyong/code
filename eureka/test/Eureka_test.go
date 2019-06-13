package test

import (
	"github.com/fulunyong/code/eureka"
	"testing"
)

func TestEureka(t *testing.T) {
	client := eureka.NewClient([]string{
		"http://106.13.10.111:1002/eureka",
		"http://106.13.10.111:1001/eureka",
	})

	instance := eureka.NewInstanceInfo("127.0.0.1", "go", "127.0.0.1"+":"+":8001", "127.0.0.1", 8001, 30, false) //Create a new instance to register
	instanceId := "127.0.0.1" + ":" + ":8001"
	_ = client.RegisterInstance(instance.App, instance) // Register new instance in your eureka(s)
	_, _ = client.GetInstance(instance.App, instanceId) // retrieve the instance from "test.com" inside "test"" app
	client.LocalInstanceInfo = *instance
	er := client.SendHeartbeat(instance.App, instanceId) // say to eureka th
	if er != nil {
		t.Fatal(er)
	}
}
