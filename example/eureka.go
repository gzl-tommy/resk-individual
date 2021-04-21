package main

import (
	"github.com/gzl-tommy/go-eureka-client/eureka"
	"time"
)

func main() {
	cfg := eureka.Config{
		DialTimeout: time.Second * 10,
	}
	client := eureka.NewClientByConfig([]string{
		"http://192.168.5.101:8761/eureka",
	}, cfg)
	appName := "Go-Example"
	instance := eureka.NewInstanceInfo(
		"test.com", appName,
		"192.168.5.101",
		8080, 30,
		false)
	client.RegisterInstance(appName, instance)
	client.Start()
	c := make(chan int, 1)
	<-c
}
